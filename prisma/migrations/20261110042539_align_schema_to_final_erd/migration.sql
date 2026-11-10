-- Menyelaraskan skema M1 dengan ERD terbaru.
-- Kondisi awal (db pull): users 2 baris, roles 4 baris, divisions/permissions/audit_logs kosong.

-- 0. Pengaman: tabel yang dibuat ulang atau diubah bentuknya harus kosong.
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM "role_permissions")
     OR EXISTS (SELECT 1 FROM "permissions")
     OR EXISTS (SELECT 1 FROM "audit_logs") THEN
    RAISE EXCEPTION 'role_permissions, permissions, atau audit_logs berisi data. Migrasi ini membuat ulang / mengubah bentuk tabel itu, jadi hentikan dan tinjau dulu.';
  END IF;
END $$;

-- 1. roles: tambah kolom, isi code dari name lama, ganti unik name menjadi unik per partner
ALTER TABLE "roles"
  ADD COLUMN "code" VARCHAR(50),
  ADD COLUMN "description" TEXT,
  ADD COLUMN "is_active" BOOLEAN NOT NULL DEFAULT true,
  ADD COLUMN "partner_id" INTEGER,
  ADD COLUMN "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ADD COLUMN "updated_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

UPDATE "roles" SET "code" = lower(replace("name", ' ', '_'));

ALTER TABLE "roles"
  ALTER COLUMN "code" SET NOT NULL,
  ALTER COLUMN "name" TYPE VARCHAR(100);

DROP INDEX IF EXISTS "roles_name_key";
CREATE UNIQUE INDEX "roles_partner_id_code_key" ON "roles"("partner_id", "code");
-- NULL dianggap berbeda di unique biasa, jadi peran bawaan platform dijaga terpisah
CREATE UNIQUE INDEX "roles_code_platform_key" ON "roles"("code") WHERE "partner_id" IS NULL;

-- 2. divisions: tambah code, partner_id, updated_at
ALTER TABLE "divisions"
  ADD COLUMN "code" VARCHAR(50),
  ADD COLUMN "partner_id" INTEGER,
  ADD COLUMN "updated_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

UPDATE "divisions" SET "code" = 'DIV_' || "id" WHERE "code" IS NULL;

ALTER TABLE "divisions"
  ALTER COLUMN "code" SET NOT NULL,
  ALTER COLUMN "name" TYPE VARCHAR(100),
  ALTER COLUMN "created_at" TYPE TIMESTAMPTZ(3) USING "created_at" AT TIME ZONE 'UTC';

DROP INDEX IF EXISTS "divisions_name_key";
CREATE UNIQUE INDEX "divisions_partner_id_code_key" ON "divisions"("partner_id", "code");
CREATE UNIQUE INDEX "divisions_code_platform_key" ON "divisions"("code") WHERE "partner_id" IS NULL;

-- 3. users: full_name -> name, kolom tenant dan penanda super admin, role_id jadi opsional
ALTER TABLE "users" RENAME COLUMN "full_name" TO "name";

ALTER TABLE "users"
  ADD COLUMN "partner_id" INTEGER,
  ADD COLUMN "is_super_admin" BOOLEAN NOT NULL DEFAULT false,
  ADD COLUMN "created_by" INTEGER,
  ADD COLUMN "deactivated_at" TIMESTAMPTZ(3);

ALTER TABLE "users"
  ALTER COLUMN "name" TYPE VARCHAR(150),
  ALTER COLUMN "email" TYPE VARCHAR(255),
  ALTER COLUMN "role_id" DROP NOT NULL,
  ALTER COLUMN "created_at" TYPE TIMESTAMPTZ(3) USING "created_at" AT TIME ZONE 'UTC',
  ALTER COLUMN "updated_at" TYPE TIMESTAMPTZ(3) USING "updated_at" AT TIME ZONE 'UTC';

-- FK role dan divisi dibuat ulang agar sama-sama RESTRICT (data dinonaktifkan, bukan dihapus)
ALTER TABLE "users" DROP CONSTRAINT IF EXISTS "users_role_id_fkey";
ALTER TABLE "users" DROP CONSTRAINT IF EXISTS "users_division_id_fkey";
ALTER TABLE "users"
  ADD CONSTRAINT "users_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT "users_division_id_fkey" FOREIGN KEY ("division_id") REFERENCES "divisions"("id") ON DELETE RESTRICT ON UPDATE CASCADE,
  ADD CONSTRAINT "users_created_by_fkey" FOREIGN KEY ("created_by") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE;

CREATE INDEX "users_partner_id_idx" ON "users"("partner_id");

-- 4. permissions dan role_permissions: PK berpindah dari id ke code (tabel kosong, dibuat ulang)
DROP TABLE "role_permissions";
DROP TABLE "permissions";

CREATE TABLE "permissions" (
  "code" VARCHAR(100) NOT NULL,
  "module" VARCHAR(20) NOT NULL,
  "description" TEXT NOT NULL,
  CONSTRAINT "permissions_pkey" PRIMARY KEY ("code")
);

CREATE TABLE "role_permissions" (
  "role_id" INTEGER NOT NULL,
  "permission_code" VARCHAR(100) NOT NULL,
  CONSTRAINT "role_permissions_pkey" PRIMARY KEY ("role_id", "permission_code")
);

CREATE TABLE "division_permissions" (
  "division_id" INTEGER NOT NULL,
  "permission_code" VARCHAR(100) NOT NULL,
  "created_at" TIMESTAMPTZ(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT "division_permissions_pkey" PRIMARY KEY ("division_id", "permission_code")
);

ALTER TABLE "role_permissions"
  ADD CONSTRAINT "role_permissions_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT "role_permissions_permission_code_fkey" FOREIGN KEY ("permission_code") REFERENCES "permissions"("code") ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE "division_permissions"
  ADD CONSTRAINT "division_permissions_division_id_fkey" FOREIGN KEY ("division_id") REFERENCES "divisions"("id") ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT "division_permissions_permission_code_fkey" FOREIGN KEY ("permission_code") REFERENCES "permissions"("code") ON DELETE CASCADE ON UPDATE CASCADE;

-- 5. audit_logs: entity -> resource, entity_id -> resource_id (INT), source_module dihapus, tambah kolom tenant dan pelaku
ALTER TABLE "audit_logs" RENAME COLUMN "entity" TO "resource";
ALTER TABLE "audit_logs" RENAME COLUMN "entity_id" TO "resource_id";
ALTER TABLE "audit_logs" DROP COLUMN "source_module";

ALTER TABLE "audit_logs"
  ADD COLUMN "partner_id" INTEGER,
  ADD COLUMN "actor_name" VARCHAR(150),
  ADD COLUMN "actor_email" VARCHAR(255);

ALTER TABLE "audit_logs"
  ALTER COLUMN "action" TYPE VARCHAR(80),
  ALTER COLUMN "resource" TYPE VARCHAR(50),
  ALTER COLUMN "resource_id" TYPE INTEGER USING "resource_id"::integer,
  ALTER COLUMN "ip_address" TYPE VARCHAR(45),
  ALTER COLUMN "created_at" TYPE TIMESTAMPTZ(3) USING "created_at" AT TIME ZONE 'UTC';

CREATE INDEX "audit_logs_partner_id_created_at_idx" ON "audit_logs"("partner_id", "created_at");
CREATE INDEX "audit_logs_resource_resource_id_idx" ON "audit_logs"("resource", "resource_id");