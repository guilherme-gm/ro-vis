-- name: GetLatestPatch :one
SELECT `latest_patch_id` FROM `loader_controller` WHERE `name` = ?;

-- name: UpsertLatestPatch :exec
INSERT INTO `loader_controller` (`name`, `latest_patch_id`, `latest_patch_name`)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
	latest_patch_id = VALUES(latest_patch_id),
	latest_patch_name = VALUES(latest_patch_name);
