-- sqlc doesn't handle left joins very well, so we make it into a view so
-- it is forced to create a new type. Not ideal, but should not cause huge performance issues
DROP VIEW IF EXISTS `previous_skill_history_vw`;

CREATE VIEW `previous_skill_history_vw` AS (
	SELECT prev.* FROM `skills_history` AS curr LEFT JOIN `skills_history` AS prev ON prev.history_id = curr.previous_history_id
);
