-- name: GetCurrentQuests :many
SELECT `quest_history`.*, `quests`.`deleted`
FROM `quests`
INNER JOIN `quest_history` ON `quests`.`latest_history_id` = `quest_history`.`history_id`;

-- name: GetQuestsIdsInPatch :many
SELECT `quest_history`.`history_id`, `quest_history`.`quest_id`
FROM `quest_history`
WHERE `quest_history`.`patch` = ?;

-- name: UpsertQuest :execresult
INSERT INTO `quests` (`quest_id`, `latest_history_id`, `deleted`)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
	latest_history_id = VALUES(latest_history_id),
	deleted = VALUES(deleted);
