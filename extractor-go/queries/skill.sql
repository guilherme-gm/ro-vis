-- name: GetSkillJobs :many
SELECT `skills_jobs`.* FROM `skills_jobs`;

-- name: GetCurrentSkills :many
SELECT `skills_history`.*, `skills`.`deleted`
FROM `skills`
INNER JOIN `skills_history` ON `skills`.`latest_history_id` = `skills_history`.`history_id`;
