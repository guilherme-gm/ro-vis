-- Removes "is_passive", "SkillCastFixedDelay",  SkillCastStatDelay, SkillSinglePostDelay, SkillGlobalPostDelay column from skills_history and view
ALTER TABLE `skills_history` DROP COLUMN `is_passive`;
ALTER TABLE `skills_history` DROP COLUMN `cast_flags`;
ALTER TABLE `skills_history` DROP COLUMN `cast_fixed_delay`;
ALTER TABLE `skills_history` DROP COLUMN `cast_stat_delay`;
ALTER TABLE `skills_history` DROP COLUMN `single_post_delay`;
ALTER TABLE `skills_history` DROP COLUMN `global_post_delay`;
