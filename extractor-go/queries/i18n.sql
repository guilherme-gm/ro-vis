-- name: GetCurrentI18ns :many
SELECT `i18n_history`.*, `i18ns`.`deleted`
FROM `i18ns`
INNER JOIN `i18n_history` ON `i18ns`.`latest_history_id` = `i18n_history`.`history_id`;

-- name: GetI18nsIdsInUpdate :many
SELECT `i18n_history`.`history_id`, `i18n_history`.`i18n_id`
FROM `i18n_history`
WHERE `i18n_history`.`update` = ?;

-- name: UpsertI18n :execresult
INSERT INTO `i18ns` (`i18n_id`, `latest_history_id`, `deleted`)
VALUES (?, ?, ?)
ON DUPLICATE KEY UPDATE
	latest_history_id = VALUES(latest_history_id),
	deleted = VALUES(deleted);

-- name: GetChangedI18ns :many
SELECT sqlc.embed(current), sqlc.embed(previous), latest.update lastUpdate
FROM `i18n_history` current
LEFT JOIN `previous_i18n_history_vw` previous ON `previous`.`history_id` = `current`.`previous_history_id`
LEFT JOIN `i18ns` latest_id ON `latest_id`.`i18n_id` = `current`.`i18n_id`
LEFT JOIN `i18n_history` latest ON `latest_id`.`latest_history_id` = `latest`.`history_id`
WHERE `current`.`update` = ?
ORDER BY `current`.`history_id`
LIMIT ?, ?;

-- name: CountChangedI18nsInUpdate :one
SELECT COUNT(*)
FROM `i18n_history`
WHERE `update` = ?;

-- name: GetI18nHistory :many
SELECT sqlc.embed(current), sqlc.embed(previous)
FROM `i18n_history` current
LEFT JOIN `previous_i18n_history_vw` previous ON `previous`.`history_id` = `current`.`previous_history_id`
WHERE `current`.`i18n_id` = ?
ORDER BY `current`.`history_id` ASC
LIMIT ?, ?;

-- name: GetI18nList :many
SELECT `i18n_history`.`i18n_id`, `i18n_history`.`pt_br_text`, `i18n_history`.`update` lastUpdate
FROM `i18ns`
INNER JOIN `i18n_history` ON `i18n_history`.`history_id` = `i18ns`.`latest_history_id`
WHERE `i18ns`.`deleted` = FALSE
ORDER BY `i18n_history`.`i18n_id` ASC
LIMIT ?, ?;

-- name: CountI18ns :one
SELECT COUNT(*)
FROM `i18ns`
INNER JOIN `i18n_history` ON `i18n_history`.`history_id` = `i18ns`.`latest_history_id`
WHERE `i18ns`.`deleted` = FALSE;

-- name: GetStrings :many
SELECT `i18n_history`.`i18n_id`, `i18n_history`.`pt_br_text`
FROM `i18ns`
INNER JOIN `i18n_history` ON `i18n_history`.`history_id` = `i18ns`.`latest_history_id`
WHERE `i18ns`.`i18n_id` IN (sqlc.slice('ids'));
