-- name: GetLatestUpdate :one
SELECT `last_update_date` FROM `loader_controller` WHERE `name` = ?;

-- name: UpsertLatestUpdate :exec
INSERT INTO `loader_controller` (`name`, `last_update_date`)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE
	last_update_date = VALUES(last_update_date);
