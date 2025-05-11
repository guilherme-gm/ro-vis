-- name: GetCurrentItems :many
SELECT `item_history`.*, `items`.`deleted`
FROM `items`
INNER JOIN `item_history` ON `items`.`latest_history_id` = `item_history`.`history_id`;

-- name: GetItemIdsInUpdate :many
SELECT `item_history`.`history_id`, `item_history`.`item_id`
FROM `item_history`
WHERE `item_history`.`update` = ?;

-- name: UpsertItem :execresult
INSERT INTO `items` (`item_id`, `latest_history_id`, `deleted`)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
	latest_history_id = VALUES(latest_history_id),
	deleted = VALUES(deleted);
