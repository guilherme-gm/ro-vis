CREATE TABLE `quest_history` (
	`history_id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
	`previous_history_id` int,
	`quest_id` int NOT NULL,
	`file_version` int NOT NULL,
	`patch` varchar(255) NOT NULL,
	`title` varchar(255),
	`description` text,
	`summary` varchar(255),
	`old_image` varchar(255),
	`icon_name` varchar(255),
	`npc_spr` varchar(255),
	`npc_navi` varchar(255),
	`npc_pos_x` int,
	`npc_pos_y` int,
	`reward_exp` varchar(255),
	`reward_jexp` varchar(255),
	`reward_item_list` text,
	`cool_time_quest` int
);

CREATE TABLE `quests` (
	`quest_id` int NOT NULL PRIMARY KEY,
	`latest_history_id` int NOT NULL,
	`deleted` boolean NOT NULL
);
