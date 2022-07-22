-- 2022.7.16
-- create database
CREATE DATABASE `horizon` DEFAULT CHARACTER SET utf8mb4;
-- insert admin user
INSERT INTO users(id, user_name, `password`, `status`, created_at, updated_at)
VALUES('1', 'admin', '$2a$10$kfcWYvoUMzY7BtQw/IrMd.aVilOs6.b/xbt560dqR467qz28TpL8K', '1', NOW(), NOW());
