import "dotenv/config";
import { PrismaPg } from "@prisma/adapter-pg";
import bcrypt from "bcryptjs";
import { PrismaClient } from "../generated/prisma/client";

const adapter = new PrismaPg({ connectionString: process.env.DATABASE_URL! });
const prisma = new PrismaClient({ adapter });

// Peran bawaan platform, dipakai bersama oleh semua partner. code dipakai di token, name untuk tampilan.
const ROLES = [
  { code: "administrator", name: "Administrator" },
  { code: "manajer", name: "Manajer" },
  { code: "partner", name: "Partner" },
  { code: "customer", name: "Customer" },
];

// Katalog izin sistem. Format code: modul.sumber.aksi (huruf kecil).
// Kolom module diambil dari segmen pertama code.
// Izin M2 sampai M5 ditambahkan di sini setelah daftarnya diterima dari tiap kelompok.
const PERMISSIONS = [
  { code: "m1.users.read", description: "Melihat daftar dan detail akun pengguna" },
  { code: "m1.users.manage", description: "Menambah, mengubah, dan menonaktifkan akun pengguna" },
  { code: "m1.roles.manage", description: "Mengelola peran dan izin peran" },
  { code: "m1.divisions.manage", description: "Mengelola divisi dan izin divisi" },
  { code: "m1.audit.read", description: "Melihat audit log aktivitas pengguna" },

  // M4 - Penjualan
  { code: "m4.products.read", description: "Melihat katalog produk dan jasa" },
  { code: "m4.products.manage", description: "Menambah, mengubah, dan menonaktifkan produk dan jasa" },
  { code: "m4.customers.read", description: "Melihat daftar dan detail pelanggan" },
  { code: "m4.customers.manage", description: "Menambah dan mengubah data pelanggan" },
  { code: "m4.quotations.read", description: "Melihat penawaran harga (quotation)" },
  { code: "m4.quotations.manage", description: "Membuat dan mengubah draf penawaran harga" },
  { code: "m4.quotations.approve", description: "Menyetujui atau menolak penawaran harga" },
  { code: "m4.invoices.read", description: "Melihat daftar dan detail faktur tagihan (invoice)" },
  { code: "m4.invoices.manage", description: "Menerbitkan faktur tagihan (konversi dari penawaran)" },
  { code: "m4.orders.read", description: "Melihat daftar dan riwayat pesanan" },
  { code: "m4.orders.create", description: "Membuat pesanan dari keranjang belanja (checkout)" },
  { code: "m4.orders.manage", description: "Memproses, menyelesaikan, atau membatalkan pesanan" },
  { code: "m4.payments.manage", description: "Memulai dan memproses pembayaran tagihan (Midtrans)" },
  { code: "m4.taxes.manage", description: "Mengatur tarif pajak penjualan" },
  { code: "m4.reports.read", description: "Melihat ringkasan penjualan dan laporan dashboard" },
];

const codesOf = (prefix: string) =>
  PERMISSIONS.filter((p) => p.code.startsWith(prefix)).map((p) => p.code);

// Izin per peran bawaan platform. Kunci = roles.code.
// Seed ini hanya menambah; izin yang dihapus dari daftar tidak otomatis dicabut dari database.
const ROLE_GRANTS: Record<string, string[]> = {
  // Admin partner: semua izin M1 + izin M4 sesuai daftar M4
  partner: [
    ...codesOf("m1."),
    "m4.products.manage",
    "m4.products.read",
    "m4.customers.manage",
    "m4.customers.read",
    "m4.quotations.manage",
    "m4.quotations.read",
    "m4.quotations.approve",
    "m4.invoices.manage",
    "m4.invoices.read",
    "m4.orders.manage",
    "m4.orders.read",
    "m4.reports.read",
  ],
  // Di daftar M4 kodenya "admin"; di platform ini peran itu adalah "administrator"
  administrator: codesOf("m4."),
  customer: [
    "m4.products.read",
    "m4.orders.create",
    "m4.orders.read",
    "m4.invoices.read",
    "m4.payments.manage",
  ],
};

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

  // roles.code unik global, jadi upsert bisa dipakai. Peran yang sudah ada tidak ditimpa.
  for (const role of ROLES) {
    await prisma.role.upsert({
      where: { code: role.code },
      update: {},
      create: role,
    });
  }

  const adminRole = await prisma.role.findUniqueOrThrow({
    where: { code: "administrator" },
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

  // Pemberian izin ke peran bawaan. Cek dulu semua kode ada di katalog (pengaman salah ketik).
  const known = new Set(PERMISSIONS.map((p) => p.code));
  for (const [roleCode, codes] of Object.entries(ROLE_GRANTS)) {
    for (const c of codes) {
      if (!known.has(c)) {
        throw new Error(`Izin "${c}" untuk peran "${roleCode}" tidak ada di katalog PERMISSIONS`);
      }
    }
    const role = await prisma.role.findUniqueOrThrow({
      where: { code: roleCode },
    });
    await prisma.rolePermission.createMany({
      data: codes.map((c) => ({ roleId: role.id, permissionCode: c })),
      skipDuplicates: true,
    });
  }

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