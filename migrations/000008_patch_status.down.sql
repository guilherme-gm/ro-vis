-- Remove status column from patches table
ALTER TABLE patches
DROP COLUMN `status`;
