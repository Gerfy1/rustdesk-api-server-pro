-- Fix user roles for existing users
-- ⚠️ NOTE: This is now done AUTOMATICALLY when the backend starts!
-- You only need to run this manually if you want to fix roles without restarting.

-- Update admin user to Super Admin (role 4)
UPDATE user SET role = 4 WHERE username = 'admin' OR is_admin = 1;

-- Update regular users to User role (role 1) if they don't have a role
UPDATE user SET role = 1 WHERE role = 0 OR role IS NULL;

-- Verify the changes
SELECT id, username, name, email, is_admin, role, 
  CASE role 
    WHEN 1 THEN 'USER'
    WHEN 2 THEN 'SUPPORT'
    WHEN 3 THEN 'SUPPORT_N2'
    WHEN 4 THEN 'SUPER_ADMIN'
    ELSE 'UNKNOWN'
  END as role_name
FROM user;
