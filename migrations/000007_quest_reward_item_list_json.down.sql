-- Convert reward_item_list column back to text type
ALTER TABLE `quest_history` MODIFY COLUMN `reward_item_list` text;
