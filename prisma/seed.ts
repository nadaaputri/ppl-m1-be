import "dotenv/config";
import { PrismaPg } from "@prisma/adapter-pg";
import bcrypt from "bcryptjs";
import { PrismaClient } from "../generated/prisma/client";

const adapter = new PrismaPg({ connectionString: process.env.DATABASE_URL! });
const prisma = new PrismaClient({ adapter });

// Peran bawaan platform (partner_id = NULL). code dipakai di token, name untuk tampilan.
const ROLES = [
  { code: "administrator", name: "Administrator" },
  { code: "manajer", name: "Manajer" },
  { code: "partner", name: "Partner" },
  { code: "customer", name: "Customer" },
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

  console.log("Seeder selesai: role dan akun admin awal siap.");
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });