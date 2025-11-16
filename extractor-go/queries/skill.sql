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
