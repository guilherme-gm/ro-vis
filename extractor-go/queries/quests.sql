-- name: GetCurrentQuests :many
SELECT `quest_history`.*, `quests`.`deleted`
FROM `quests`
INNER JOIN `quest_history` ON `quests`.`latest_history_id` = `quest_history`.`history_id`;

-- name: GetQuestsIdsInUpdate :many
SELECT `quest_history`.`history_id`, `quest_history`.`quest_id`
FROM `quest_history`
WHERE `quest_history`.`update` = ?;

-- name: UpsertQuest :execresult
INSERT INTO `quests` (`quest_id`, `latest_history_id`, `deleted`)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
	latest_history_id = VALUES(latest_history_id),
	deleted = VALUES(deleted);

-- name: GetChangedQuests :many
SELECT sqlc.embed(current), sqlc.embed(previous), latest.update lastUpdate
FROM `quest_history` current
LEFT JOIN `previous_quest_history_vw` previous ON `previous`.`history_id` = `current`.`previous_history_id`
LEFT JOIN `quests` latest_id ON `latest_id`.`quest_id` = `current`.`quest_id`
LEFT JOIN `quest_history` latest ON `latest_id`.`latest_history_id` = `latest`.`history_id`
WHERE `current`.`update` = ?
ORDER BY `current`.`history_id`
LIMIT ?, ?;

-- name: CountChangedQuestsInUpdate :one
SELECT COUNT(*)
FROM `quest_history`
WHERE `update` = ?;

-- name: GetQuestHistory :many
SELECT sqlc.embed(current), sqlc.embed(previous)
FROM `quest_history` current
LEFT JOIN `previous_quest_history_vw` previous ON `previous`.`history_id` = `current`.`previous_history_id`
WHERE `current`.`quest_id` = ?
ORDER BY `current`.`history_id` ASC
LIMIT ?, ?;

-- name: GetQuestList :many
SELECT `quest_history`.`quest_id`, `quest_history`.`title`, `quest_history`.`update` lastUpdate
FROM `quests`
INNER JOIN `quest_history` ON `quest_history`.`history_id` = `quests`.`latest_history_id`
WHERE `quests`.`deleted` = FALSE
ORDER BY `quest_history`.`quest_id` ASC
LIMIT ?, ?;

-- name: CountQuests :one
SELECT COUNT(*)
FROM `quests`
INNER JOIN `quest_history` ON `quests`.`latest_history_id` = `quest_history`.`history_id`
WHERE `quests`.`deleted` = FALSE;
