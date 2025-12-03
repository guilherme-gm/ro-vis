-- name: GetSkillJobs :many
SELECT `skills_jobs`.* FROM `skills_jobs`;

-- name: GetCurrentSkills :many
SELECT `skills_history`.*, `skills`.`deleted`
FROM `skills`
INNER JOIN `skills_history` ON `skills`.`latest_history_id` = `skills_history`.`history_id`;

-- name: GetSkillsIdsInUpdate :many
SELECT `skills_history`.`history_id`, `skills_history`.`skill_id` as ID
FROM `skills_history`
WHERE `skills_history`.`update` = ?;

-- name: GetSkillList :many
SELECT `skills_history`.`skill_id`, `skills_history`.`name`, `skills_history`.`constant`, `skills_history`.`update` lastUpdate
FROM `skills`
INNER JOIN `skills_history` ON `skills_history`.`history_id` = `skills`.`latest_history_id`
WHERE `skills`.`deleted` = FALSE
ORDER BY `skills_history`.`skill_id` ASC
LIMIT ?, ?;

-- name: CountSkills :one
SELECT COUNT(*)
FROM `skills`
INNER JOIN `skills_history` ON `skills`.`latest_history_id` = `skills_history`.`history_id`
WHERE `skills`.`deleted` = FALSE;

-- name: GetSkillHistory :many
SELECT sqlc.embed(current), sqlc.embed(previous)
FROM `skills_history` current
LEFT JOIN `previous_skill_history_vw` previous ON `previous`.`history_id` = `current`.`previous_history_id`
WHERE `current`.`skill_id` = ?
ORDER BY `current`.`history_id` ASC
LIMIT ?, ?;

-- name: CountChangedSkillsInUpdate :one
SELECT COUNT(*)
FROM `skills_history`
WHERE `update` = ?;

-- name: GetChangedSkills :many
SELECT sqlc.embed(current), sqlc.embed(previous), latest.update lastUpdate
FROM `skills_history` current
LEFT JOIN `previous_skill_history_vw` previous ON `previous`.`history_id` = `current`.`previous_history_id`
LEFT JOIN `skills` latest_id ON `latest_id`.`skill_id` = `current`.`skill_id`
LEFT JOIN `skills_history` latest ON `latest_id`.`latest_history_id` = `latest`.`history_id`
WHERE `current`.`update` = ?
ORDER BY `current`.`history_id`
LIMIT ?, ?;
