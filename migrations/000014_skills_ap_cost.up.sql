-- Adds "ap_cost" column to skills_history and view
ALTER TABLE `skills_history` ADD COLUMN `ap_cost` json DEFAULT NULL AFTER `sp_cost`;
