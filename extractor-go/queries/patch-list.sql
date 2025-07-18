-- name: InsertPatch :exec
INSERT INTO `patches` (`id`, `name`, `date`, `files`, `status`) VALUES (?, ?, ?, ?, ?);

-- name: ListPatches :many
SELECT * FROM `patches` ORDER BY `id` ASC;

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
