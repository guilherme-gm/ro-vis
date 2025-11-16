-- name: GetSkillJobs :many
SELECT `skills_jobs`.* FROM `skills_jobs`;

-- name: GetCurrentSkills :many
SELECT `skills_history`.*, `skills`.`deleted`
FROM `skills`
INNER JOIN `skills_history` ON `skills`.`latest_history_id` = `skills_history`.`history_id`;

-- name: GetSkillsIdsInUpdate :many
SELECT `skills_history`.`history_id`, `skills_history`.`skill_id`
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
