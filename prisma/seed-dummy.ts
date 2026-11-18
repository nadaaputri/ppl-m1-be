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
const DUMMY_DIVISION_PREFIX = "dummy_";

function dbHost(): string {
  try {
    return new URL(process.env.DATABASE_URL!).host;
  } catch {
    return "(DATABASE_URL tidak terbaca)";
  }
}

// divisions.code unik global, jadi upsert aman dijalankan berulang
async function ensureDivision(code: string, name: string) {
  return prisma.division.upsert({
    where: { code },
    update: {},
    create: { code, name },
  });
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

  // Peran bawaan platform (dibuat oleh seed.ts)
  const partnerRole = await prisma.role.findUniqueOrThrow({ where: { code: "partner" } });
  const managerRole = await prisma.role.findUniqueOrThrow({ where: { code: "manajer" } });

  // Divisi uji (global)
  const penjualan = await ensureDivision(`${DUMMY_DIVISION_PREFIX}penjualan`, "Penjualan (dummy)");
  const keuangan = await ensureDivision(`${DUMMY_DIVISION_PREFIX}keuangan`, "Keuangan (dummy)");

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

  // Akun bawahan partner A (peran manajer + divisi), dibuat oleh admin partner A
  await ensureUser({
    email: `manajer1.partner1${DUMMY_DOMAIN}`,
    name: "Manajer 1 Partner 1 (dummy)",
    passwordHash,
    partnerId: PARTNER_A,
    roleId: managerRole.id,
    divisionId: penjualan.id,
    createdById: adminA.id,
  });
  await ensureUser({
    email: `manajer2.partner1${DUMMY_DOMAIN}`,
    name: "Manajer 2 Partner 1 (dummy)",
    passwordHash,
    partnerId: PARTNER_A,
    roleId: managerRole.id,
    divisionId: keuangan.id,
    createdById: adminA.id,
  });

  // Akun nonaktif (AC-2: akun nonaktif tidak bisa login, datanya tetap tersimpan)
  await ensureUser({
    email: `nonaktif.partner1${DUMMY_DOMAIN}`,
    name: "Akun Nonaktif Partner 1 (dummy)",
    passwordHash,
    partnerId: PARTNER_A,
    roleId: managerRole.id,
    divisionId: penjualan.id,
    createdById: adminA.id,
    isActive: false,
  });

  console.log("\nData dummy siap:");
  console.table([
    { email: `admin.partner1${DUMMY_DOMAIN}`, partner_id: PARTNER_A, role: "partner", divisi: "-", aktif: true },
    { email: `admin.partner2${DUMMY_DOMAIN}`, partner_id: PARTNER_B, role: "partner", divisi: "-", aktif: true },
    { email: `manajer1.partner1${DUMMY_DOMAIN}`, partner_id: PARTNER_A, role: "manajer", divisi: penjualan.code, aktif: true },
    { email: `manajer2.partner1${DUMMY_DOMAIN}`, partner_id: PARTNER_A, role: "manajer", divisi: keuangan.code, aktif: true },
    { email: `nonaktif.partner1${DUMMY_DOMAIN}`, partner_id: PARTNER_A, role: "manajer", divisi: penjualan.code, aktif: false },
  ]);
}

async function clean() {
  // Urutan: log -> akun -> divisi dummy (FK akun ke divisi memakai RESTRICT).
  // Peran tidak dihapus karena peran bawaan platform dipakai semua partner.
  const [logs, users, divisions] = await prisma.$transaction([
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
    prisma.division.deleteMany({
      where: { code: { startsWith: DUMMY_DIVISION_PREFIX } },
    }),
  ]);

  console.log(
    `Dibersihkan: ${logs.count} log, ${users.count} akun, ${divisions.count} divisi.`,
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