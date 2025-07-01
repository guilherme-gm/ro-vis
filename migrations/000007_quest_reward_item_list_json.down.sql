-- Convert reward_item_list column back to text type
ALTER TABLE `quest_history` MODIFY COLUMN `reward_item_list` text;

-- Update the view to use the new column type
DROP VIEW `previous_quest_history_vw`;
CREATE VIEW previous_quest_history_vw AS (
	SELECT prev.* FROM quest_history curr LEFT JOIN quest_history prev ON prev.history_id = curr.previous_history_id
);
