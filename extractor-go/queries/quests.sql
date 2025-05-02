-- name: GetCurrentQuests :many
SELECT `quest_history`.*, `quests`.`deleted`
FROM `quests`
INNER JOIN `quest_history` ON `quests`.`latest_history_id` = `quest_history`.`history_id`;

-- name: InsertQuestHistory :execresult
INSERT INTO `quest_history` (
	`previous_history_id`, -- 1
	`quest_id`, -- 2
	`file_version`, -- 3
	`patch`, -- 4
	`title`, -- 5
	`description`, -- 6
	`summary`, -- 7
	`old_image`, -- 8
	`icon_name`, -- 9
	`npc_spr`, -- 10
	`npc_navi`, -- 11
	`npc_pos_x`, -- 12
	`npc_pos_y`, -- 13
	`reward_exp`, -- 14
	`reward_jexp`, -- 15
	`reward_item_list`, -- 16
	`cool_time_quest` -- 17
)
VALUES (
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpsertQuest :execresult
INSERT INTO `quests` (`quest_id`, `latest_history_id`, `deleted`)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
	latest_history_id = VALUES(latest_history_id),
	deleted = VALUES(deleted);
