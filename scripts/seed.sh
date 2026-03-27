#!/bin/bash
set -e

DB_URL="${MIGRATE_URL:-postgres://goodie:goodie_secret@localhost:5432/goodie?sslmode=disable}"

echo "Seeding database..."

psql "$DB_URL" <<'SQL'
-- Insert admin user (password: admin123)
INSERT INTO users (email, phone, password_hash, full_name, role, is_verified, status)
VALUES (
  'admin@goodie.vn',
  '0900000001',
  '$2a$12$LJ3m4ys5gmPOHMwmzrPXxe8b5H2Q6ROhVcLiVnDcWx3QWabCGlAXm',
  'System Admin',
  'admin',
  true,
  'active'
) ON CONFLICT (email) DO NOTHING;

-- Insert sample merchant user
INSERT INTO users (email, phone, password_hash, full_name, role, is_verified, status)
VALUES (
  'merchant@goodie.vn',
  '0900000002',
  '$2a$12$LJ3m4ys5gmPOHMwmzrPXxe8b5H2Q6ROhVcLiVnDcWx3QWabCGlAXm',
  'Demo Merchant',
  'merchant',
  true,
  'active'
) ON CONFLICT (email) DO NOTHING;

-- Insert sample client user
INSERT INTO users (email, phone, password_hash, full_name, role, is_verified, status)
VALUES (
  'client@goodie.vn',
  '0900000003',
  '$2a$12$LJ3m4ys5gmPOHMwmzrPXxe8b5H2Q6ROhVcLiVnDcWx3QWabCGlAXm',
  'Demo Client',
  'client',
  true,
  'active'
) ON CONFLICT (email) DO NOTHING;

SELECT 'Seeded ' || count(*) || ' users' FROM users;
SQL

echo "Seed complete!"
