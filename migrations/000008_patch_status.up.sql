-- Add status column to patches table
ALTER TABLE patches
ADD COLUMN `status` ENUM('pending', 'extracted', 'gone', 'skipped') NOT NULL;
