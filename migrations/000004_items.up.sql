CREATE TABLE `item_history` (
	`history_id` int NOT NULL AUTO_INCREMENT PRIMARY KEY,
	`previous_history_id` int,
	`item_id` int NOT NULL,
	`file_version` int NOT NULL,
	`update` varchar(255) NOT NULL,
	`identified_name` varchar(50),
	`identified_description` text,
	`identified_sprite` varchar(255),
	`unidentified_name` varchar(255),
	`unidentified_description` text,
	`unidentified_sprite` varchar(255),
	`slot_count` tinyint NOT NULL default 0,
	`is_book` boolean NOT NULL default false,
	`can_use_buying_store` boolean NOT NULL default false,
	`card_prefix` varchar(255),
	`card_is_postfix` boolean NOT NULL default false,
	`card_illustration` varchar(255),
	`class_num` int,
	`move_info` json
);

CREATE TABLE `items` (
	`item_id` int NOT NULL PRIMARY KEY,
	`latest_history_id` int NOT NULL,
	`deleted` boolean NOT NULL
);
