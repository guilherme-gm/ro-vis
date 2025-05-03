-- name: InsertPatch :exec
INSERT INTO `patches` (`name`, `date`, `files`) VALUES (?, ?, ?);

-- name: ListPatches :many
SELECT * FROM `patches`;
