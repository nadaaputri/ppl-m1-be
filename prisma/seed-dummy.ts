import "dotenv/config";
import { PrismaPg } from "@prisma/adapter-pg";
import bcrypt from "bcryptjs";
import { PrismaClient } from "../generated/prisma/client";

// Data uji untuk pengujian endpoint M1.
//   Isi data dummy : npx tsx prisma/seed-dummy.ts
//   Hapus data dummy: npx tsx prisma/seed-dummy.ts --clean
// Perlu DUMMY_PASSWORD di .env (hanya untuk mode isi data).

const adapter = new PrismaPg({ connectionString: process.env.DATABASE_URL! });
const prisma = new PrismaClient({ adapter });

// ID partner dummy. Rentang >= 9000 dicadangkan untuk data uji supaya
// tidak bentrok dengan ID partner asli dari M3 (auto increment mulai dari 1).
const PARTNER_A = 9001;
const PARTNER_B = 9002;
const DUMMY_PARTNERS = [PARTNER_A, PARTNER_B];
const DUMMY_DOMAIN = "@dummy.test";

function dbHost(): string {
  try {
    return new URL(process.env.DATABASE_URL!).host;
  } catch {
    return "(DATABASE_URL tidak terbaca)";
  }
}

async function ensureDivision(partnerId: number, code: string, name: string) {
  const found = await prisma.division.findFirst({ where: { partnerId, code } });
  return found ?? prisma.division.create({ data: { partnerId, code, name } });
}

async function ensureRole(partnerId: number, code: string, name: string) {
  const found = await prisma.role.findFirst({ where: { partnerId, code } });
  return found ?? prisma.role.create({ data: { partnerId, code, name } });
}

type DummyUser = {
  email: string;
  name: string;
  passwordHash: string;
  partnerId: number;
  roleId: number;
  divisionId?: number;
  createdById?: number;
  isActive?: boolean;
};

async function ensureUser(u: DummyUser) {
  const isActive = u.isActive ?? true;
  return prisma.user.upsert({
    where: { email: u.email },
    update: {}, // akun yang sudah ada tidak ditimpa
    create: {
      email: u.email,
      name: u.name,
      passwordHash: u.passwordHash,
      partnerId: u.partnerId,
      roleId: u.roleId,
      divisionId: u.divisionId ?? null,
      createdById: u.createdById ?? null,
      isActive,
      deactivatedAt: isActive ? null : new Date(),
    },
  });
}

async function seed() {
  const password = process.env.DUMMY_PASSWORD;
  if (!password) {
    throw new Error("DUMMY_PASSWORD wajib diisi di .env");
  }
  if (password.length < 8 || !/[A-Za-z]/.test(password) || !/\d/.test(password)) {
    throw new Error("DUMMY_PASSWORD minimal 8 karakter, berisi huruf dan angka");
  }

  // cost 10 = bcrypt.DefaultCost di backend Go
  const passwordHash = await bcrypt.hash(password, 10);

  const partnerRole = await prisma.role.findFirstOrThrow({
    where: { code: "partner", partnerId: null },
  });

  // Divisi: kode yang sama (penjualan) di dua partner sengaja dibuat,
  // untuk membuktikan unik per (partner_id, code).
  const penjualanA = await ensureDivision(PARTNER_A, "penjualan", "Penjualan");
  const keuanganA = await ensureDivision(PARTNER_A, "keuangan", "Keuangan");
  const penjualanB = await ensureDivision(PARTNER_B, "penjualan", "Penjualan");

  // Peran buatan partner A
  const kasirA = await ensureRole(PARTNER_A, "kasir", "Kasir");

  // Admin partner (peran bawaan "partner", tanpa divisi)
  const adminA = await ensureUser({
    email: `admin.partner1${DUMMY_DOMAIN}`,
    name: "Admin Partner 1 (dummy)",
    passwordHash,
    partnerId: PARTNER_A,
    roleId: partnerRole.id,
  });
  await ensureUser({
    email: `admin.partner2${DUMMY_DOMAIN}`,
    name: "Admin Partner 2 (dummy)",
    passwordHash,
    partnerId: PARTNER_B,
    roleId: partnerRole.id,
  });

  // Akun bawahan partner A (peran kasir + divisi), dibuat oleh admin partner A
  await ensureUser({
    email: `kasir1.partner1${DUMMY_DOMAIN}`,
    name: "Kasir 1 Partner 1 (dummy)",
    passwordHash,
    partnerId: PARTNER_A,
    roleId: kasirA.id,
    divisionId: penjualanA.id,
    createdById: adminA.id,
  });
  await ensureUser({
    email: `kasir2.partner1${DUMMY_DOMAIN}`,
    name: "Kasir 2 Partner 1 (dummy)",
    passwordHash,
    partnerId: PARTNER_A,
    roleId: kasirA.id,
    divisionId: keuanganA.id,
    createdById: adminA.id,
  });

  // Akun nonaktif (AC-11: akun nonaktif tidak bisa login)
  await ensureUser({
    email: `nonaktif.partner1${DUMMY_DOMAIN}`,
    name: "Akun Nonaktif Partner 1 (dummy)",
    passwordHash,
    partnerId: PARTNER_A,
    roleId: kasirA.id,
    divisionId: penjualanA.id,
    createdById: adminA.id,
    isActive: false,
  });

  console.log("\nData dummy siap:");
  console.table([
    { email: `admin.partner1${DUMMY_DOMAIN}`, partner_id: PARTNER_A, role: "partner", divisi: "-", aktif: true },
    { email: `admin.partner2${DUMMY_DOMAIN}`, partner_id: PARTNER_B, role: "partner", divisi: "-", aktif: true },
    { email: `kasir1.partner1${DUMMY_DOMAIN}`, partner_id: PARTNER_A, role: "kasir", divisi: "penjualan", aktif: true },
    { email: `kasir2.partner1${DUMMY_DOMAIN}`, partner_id: PARTNER_A, role: "kasir", divisi: "keuangan", aktif: true },
    { email: `nonaktif.partner1${DUMMY_DOMAIN}`, partner_id: PARTNER_A, role: "kasir", divisi: "penjualan", aktif: false },
  ]);
  console.log(`Divisi partner ${PARTNER_B} (untuk uji "divisi milik partner lain ditolak"): id ${penjualanB.id}`);
}

async function clean() {
  // Urutan: log -> akun -> divisi -> peran (FK akun ke peran/divisi memakai RESTRICT).
  // role_permissions dan division_permissions ikut terhapus lewat CASCADE.
  const [logs, users, divisions, roles] = await prisma.$transaction([
    prisma.auditLog.deleteMany({
      where: {
        OR: [
          { partnerId: { in: DUMMY_PARTNERS } },
          { actorEmail: { endsWith: DUMMY_DOMAIN } },
        ],
      },
    }),
    prisma.user.deleteMany({
      where: {
        OR: [
          { partnerId: { in: DUMMY_PARTNERS } },
          { email: { endsWith: DUMMY_DOMAIN } },
        ],
      },
    }),
    prisma.division.deleteMany({ where: { partnerId: { in: DUMMY_PARTNERS } } }),
    prisma.role.deleteMany({ where: { partnerId: { in: DUMMY_PARTNERS } } }),
  ]);

  console.log(
    `Dibersihkan: ${logs.count} log, ${users.count} akun, ${divisions.count} divisi, ${roles.count} peran.`,
  );
}

async function main() {
  const mode = process.argv.includes("--clean") ? "clean" : "seed";
  console.log(`Database: ${dbHost()}  |  Mode: ${mode}`);
  if (mode === "clean") {
    await clean();
  } else {
    await seed();
  }
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });