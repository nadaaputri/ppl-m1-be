-- Peran dan divisi menjadi daftar global (bersama semua partner):
-- kolom partner_id dihapus dari roles dan divisions, code menjadi unik global.

-- 0. Pengaman: tidak boleh ada baris yang masih dimiliki partner (datanya akan bercampur / hilang).
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM "roles" WHERE "partner_id" IS NOT NULL)
     OR EXISTS (SELECT 1 FROM "divisions" WHERE "partner_id" IS NOT NULL) THEN
    RAISE EXCEPTION 'roles atau divisions masih punya baris ber-partner_id. Bersihkan data dummy dulu, lalu jalankan ulang migrasi ini.';
  END IF;
END $$;

-- 1. Hapus index lama yang bergantung pada partner_id
DROP INDEX IF EXISTS "roles_code_platform_key";
DROP INDEX IF EXISTS "roles_partner_id_code_key";
DROP INDEX IF EXISTS "divisions_code_platform_key";
DROP INDEX IF EXISTS "divisions_partner_id_code_key";

-- 2. Hapus kolom
ALTER TABLE "roles" DROP COLUMN "partner_id";
ALTER TABLE "divisions" DROP COLUMN "partner_id";

-- 3. code unik global
CREATE UNIQUE INDEX "roles_code_key" ON "roles"("code");
CREATE UNIQUE INDEX "divisions_code_key" ON "divisions"("code");