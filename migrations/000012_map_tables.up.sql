CREATE TABLE `maps_history` (
	`history_id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
	`previous_history_id` int,
	`map_id` varchar(255) NOT NULL,
	`file_version` int NOT NULL,
	`update` varchar(255) NOT NULL,
	`name` text,
	`special_code` int,
	`mp3_name` varchar(255),
	`npcs` json,
	`warps` json,
	`spawns` json,
	`created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
	KEY `idx_map_history_map_id` (`map_id`),
	KEY `idx_map_history_update` (`update`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `maps` (
	`map_id` varchar(255) NOT NULL PRIMARY KEY,
	`latest_history_id` int NOT NULL,
	`deleted` boolean NOT NULL DEFAULT FALSE,
	`updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	FOREIGN KEY (`latest_history_id`) REFERENCES `maps_history`(`history_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- sqlc doesn't handle left joins very well, so we make it into a view so
-- it is forced to create a new type. Not ideal, but should not cause huge performance issues
CREATE VIEW `previous_map_history_vw` AS (
	SELECT prev.* FROM `maps_history` AS curr LEFT JOIN `maps_history` AS prev ON prev.history_id = curr.previous_history_id
);
