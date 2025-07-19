-- Drop created_at and updated_at columns
ALTER TABLE `patches`
DROP COLUMN `updated_at`,
DROP COLUMN `created_at`;

-- Revert date column type back to DATETIME
ALTER TABLE `patches`
MODIFY COLUMN `date` DATETIME NOT NULL;
