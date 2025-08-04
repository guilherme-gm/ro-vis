-- name: GetCurrentMaps :many
SELECT `maps_history`.*, `maps`.`deleted`
FROM `maps`
INNER JOIN `maps_history` ON `maps_history`.`history_id` = `maps`.`latest_history_id`;

-- name: GetMapsIdsInUpdate :many
SELECT `maps_history`.`history_id`, `maps_history`.`map_id`
FROM `maps_history`
WHERE `maps_history`.`update` = ?;

-- name: UpsertMap :execresult
INSERT INTO `maps` (`map_id`, `latest_history_id`, `deleted`)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
	latest_history_id = VALUES(latest_history_id),
	deleted = VALUES(deleted);

-- name: GetChangedMaps :many
SELECT sqlc.embed(current), sqlc.embed(previous), latest.update lastUpdate
FROM `maps_history` current
LEFT JOIN `previous_map_history_vw` previous ON `previous`.`history_id` = `current`.`previous_history_id`
LEFT JOIN `maps` latest_id ON `latest_id`.`map_id` = `current`.`map_id`
LEFT JOIN `maps_history` latest ON `latest_id`.`latest_history_id` = `latest`.`history_id`
WHERE `current`.`update` = ?
ORDER BY `current`.`history_id`
LIMIT ?, ?;

-- name: CountChangedMapsInUpdate :one
SELECT COUNT(*)
FROM `maps_history`
WHERE `update` = ?;

-- name: GetMapHistory :many
SELECT sqlc.embed(current), sqlc.embed(previous)
FROM `maps_history` current
LEFT JOIN `previous_map_history_vw` previous ON `previous`.`history_id` = `current`.`previous_history_id`
WHERE `current`.`map_id` = ?
ORDER BY `current`.`history_id` ASC
LIMIT ?, ?;

-- name: GetMapList :many
SELECT `maps_history`.`map_id`, `maps_history`.`name`, `maps_history`.`update` lastUpdate
FROM `maps`
INNER JOIN `maps_history` ON `maps_history`.`history_id` = `maps`.`latest_history_id`
WHERE `maps`.`deleted` = FALSE
ORDER BY `maps_history`.`map_id` ASC
LIMIT ?, ?;

-- name: CountMaps :one
SELECT COUNT(*)
FROM `maps`
INNER JOIN `maps_history` ON `maps`.`latest_history_id` = `maps_history`.`history_id`
WHERE `maps`.`deleted` = FALSE;
