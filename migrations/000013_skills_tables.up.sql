-- Table to store jobinheritlist , which contains the Job IDs used in skill list files
-- We don't do a full detailed track list for those, as they should be generally stable
CREATE TABLE `skills_jobs` (
	`constant` varchar(255) NOT NULL,
	`job_id` int NOT NULL,
	`inherited_job` int,
	`inherited_job2` int,
	`first_update` varchar(255) NOT NULL,
	`last_update` varchar(255) NOT NULL,
	`deleted` boolean NOT NULL DEFAULT FALSE,
	`updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`constant`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `skills_history` (
	`history_id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
	`previous_history_id` int,
	`skill_id` int NOT NULL,
	`file_version` int NOT NULL,
	`update` varchar(255) NOT NULL,
	`constant` varchar(255),
	`name` text,
	`description` text,
	`max_level` int,
	`sp_cost` json,
	`can_select_level` boolean,
	`attack_range` json,
	`required_skills` json,
	`job_required_skills` json,
	`skill_scale` json,
	`created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
	KEY `idx_skill_history_skill_id` (`skill_id`),
	KEY `idx_skill_history_update` (`update`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `skills` (
	`skill_id` int NOT NULL PRIMARY KEY,
	`latest_history_id` int NOT NULL,
	`deleted` boolean NOT NULL DEFAULT FALSE,
	`updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (`latest_history_id`) REFERENCES `skills_history`(`history_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
