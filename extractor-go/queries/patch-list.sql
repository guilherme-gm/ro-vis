-- name: InsertPatch :exec
INSERT INTO `patches` (`id`, `name`, `date`, `files`, `status`) VALUES (?, ?, ?, ?, ?);

-- name: ListPatches :many
-- We include pending and extracted because we may have patches marked as 'extracted' but
-- that were not processed for a new record type, while gone/skipped should never be considered.
SELECT * FROM `patches` WHERE `date` <= ? AND `status` IN ('pending', 'extracted') ORDER BY `id` ASC;

-- name: ListUpdatesPatches :many
SELECT patches.*
FROM (
	SELECT `date`
	FROM `patches`
	GROUP BY `date`
	ORDER BY `date` ASC
	LIMIT ?, ?
) dates
INNER JOIN `patches` ON `patches`.`date` = dates.`date`
ORDER BY patches.`id` ASC;

-- name: GetUpdatesCount :one
SELECT COUNT(*) FROM (SELECT `date` FROM `patches` GROUP BY `date`) updates;

-- name: GetLatestPatch :one
SELECT * FROM `patches` ORDER BY `id` DESC LIMIT 1
