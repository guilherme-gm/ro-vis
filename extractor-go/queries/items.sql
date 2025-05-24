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

-- name: GetChangedItems :many
SELECT sqlc.embed(current), sqlc.embed(previous), latest.update lastUpdate
FROM `item_history` current
LEFT JOIN `previous_item_history_vw` previous ON `previous`.`history_id` = `current`.`previous_history_id`
LEFT JOIN `items` latest_id ON `latest_id`.`item_id` = `current`.`item_id`
LEFT JOIN `item_history` latest ON `latest_id`.`latest_history_id` = `latest`.`history_id`
WHERE `current`.`update` = ?
ORDER BY `current`.`history_id`
LIMIT ?, ?;

-- name: CountChangedItemsInPatch :one
SELECT COUNT(*)
FROM `item_history`
WHERE `update` = ?;

-- name: GetItemHistory :many
SELECT sqlc.embed(current), sqlc.embed(previous)
FROM `item_history` current
LEFT JOIN `previous_item_history_vw` previous ON `previous`.`history_id` = `current`.`previous_history_id`
WHERE `current`.`item_id` = ?
ORDER BY `current`.`history_id` ASC
LIMIT ?, ?;
