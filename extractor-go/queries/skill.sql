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
