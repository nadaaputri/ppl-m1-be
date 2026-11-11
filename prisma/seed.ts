import "dotenv/config";
import { PrismaPg } from "@prisma/adapter-pg";
import bcrypt from "bcryptjs";
import { PrismaClient } from "../generated/prisma/client";

const adapter = new PrismaPg({ connectionString: process.env.DATABASE_URL! });
const prisma = new PrismaClient({ adapter });

const ROLES = [
  { code: "administrator", name: "Administrator" },
  { code: "manajer", name: "Manajer" },
  { code: "partner", name: "Partner" },
  { code: "customer", name: "Customer" },
];

const PERMISSIONS = [
  { code: "m1.users.read", description: "Melihat daftar dan detail akun pengguna" },
  { code: "m1.users.manage", description: "Menambah, mengubah, dan menonaktifkan akun pengguna" },
  { code: "m1.roles.manage", description: "Mengelola peran dan izin peran" },
  { code: "m1.divisions.manage", description: "Mengelola divisi dan izin divisi" },
  { code: "m1.audit.read", description: "Melihat audit log aktivitas pengguna" },
];

async function main() {
  const email = process.env.ADMIN_EMAIL;
  const name = process.env.ADMIN_NAME ?? "Administrator";
  const password = process.env.ADMIN_PASSWORD;

  if (!email || !password) {
    throw new Error("ADMIN_EMAIL dan ADMIN_PASSWORD wajib diisi di .env");
  }
  if (
    password.length < 8 ||
    !/[A-Za-z]/.test(password) ||
    !/\d/.test(password) ||
    password.toLowerCase() === email.toLowerCase()
  ) {
    throw new Error("Kata sandi admin harus minimal 8 karakter, berisi huruf dan angka, dan tidak sama dengan email");
  }

  // upsert tidak bisa dipakai: unik (partner_id, code) tidak menerima partner_id NULL di klausa where
  for (const role of ROLES) {
    const exists = await prisma.role.findFirst({
      where: { code: role.code, partnerId: null },
    });
    if (!exists) {
      await prisma.role.create({ data: { ...role, partnerId: null } });
    }
  }

  const adminRole = await prisma.role.findFirstOrThrow({
    where: { code: "administrator", partnerId: null },
  });

  // cost 10 = bcrypt.DefaultCost di backend Go
  const passwordHash = await bcrypt.hash(password, 10);

  await prisma.user.upsert({
    where: { email },
    update: {}, // akun yang sudah ada tidak ditimpa
    create: {
      email,
      passwordHash,
      name,
      roleId: adminRole.id,
      isSuperAdmin: true, // staf platform Araya (partner_id NULL)
      isActive: true,
    },
  });

  // Katalog izin: code adalah PK, jadi upsert aman dijalankan berulang
  for (const p of PERMISSIONS) {
    const mod = p.code.split(".")[0];
    await prisma.permission.upsert({
      where: { code: p.code },
      update: { module: mod, description: p.description },
      create: { code: p.code, module: mod, description: p.description },
    });
  }

  // Semua izin m1.* diberikan ke peran bawaan "partner" (admin partner)
  const partnerRole = await prisma.role.findFirstOrThrow({
    where: { code: "partner", partnerId: null },
  });
  await prisma.rolePermission.createMany({
    data: PERMISSIONS.filter((p) => p.code.startsWith("m1.")).map((p) => ({
      roleId: partnerRole.id,
      permissionCode: p.code,
    })),
    skipDuplicates: true,
  });

  console.log("Seeder selesai: role, akun admin awal, dan katalog izin siap.");
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });