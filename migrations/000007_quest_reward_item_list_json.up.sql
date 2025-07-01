-- Convert reward_item_list column to JSON type
ALTER TABLE `quest_history` MODIFY COLUMN `reward_item_list` json;
