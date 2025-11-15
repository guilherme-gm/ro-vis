-- Adds "is_passive", "SkillCastFixedDelay",  SkillCastStatDelay, SkillSinglePostDelay, SkillGlobalPostDelay column to skills_history and view
ALTER TABLE `skills_history` ADD COLUMN `is_passive` boolean DEFAULT NULL AFTER `max_level`;
ALTER TABLE `skills_history` ADD COLUMN `cast_flags` json DEFAULT NULL AFTER `skill_scale`;
ALTER TABLE `skills_history` ADD COLUMN `cast_fixed_delay` json DEFAULT NULL AFTER `cast_flags`;
ALTER TABLE `skills_history` ADD COLUMN `cast_stat_delay` json DEFAULT NULL AFTER `cast_fixed_delay`;
ALTER TABLE `skills_history` ADD COLUMN `single_post_delay` json DEFAULT NULL AFTER `cast_stat_delay`;
ALTER TABLE `skills_history` ADD COLUMN `global_post_delay` json DEFAULT NULL AFTER `single_post_delay`;
