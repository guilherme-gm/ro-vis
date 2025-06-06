-- Renames quest_history's `patch` column to `update` to make it consistent
-- with other tables.
ALTER TABLE `quest_history` RENAME COLUMN `patch` TO `update`;

-- sqlc doesn't handle left joins very well, so we make it into a view so
-- it is forced to create a new type. Not ideal, but should not cause huge performance issues
CREATE VIEW previous_quest_history_vw AS (
	SELECT prev.* FROM quest_history curr LEFT JOIN quest_history prev ON prev.history_id = curr.previous_history_id
);
