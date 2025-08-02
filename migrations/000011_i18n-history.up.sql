CREATE TABLE `i18n_history` (
	`history_id` bigint NOT NULL AUTO_INCREMENT PRIMARY KEY,
	`previous_history_id` bigint,
	`i18n_id` bigint unsigned NOT NULL,
	`file_version` int NOT NULL,
	`update` varchar(255) NOT NULL,
	`container_file` varchar(255) NOT NULL,
	`en_text` text NOT NULL,
	`pt_br_text` text NOT NULL,
	`active` boolean NOT NULL
);

CREATE TABLE `i18ns` (
	`i18n_id` bigint unsigned NOT NULL PRIMARY KEY,
	`latest_history_id` bigint NOT NULL,
	`deleted` boolean NOT NULL
);

-- sqlc doesn't handle left joins very well, so we make it into a view so
-- it is forced to create a new type. Not ideal, but should not cause huge performance issues
CREATE VIEW `previous_i18n_history_vw` AS (
	SELECT prev.* FROM `i18n_history` AS curr LEFT JOIN `i18n_history` AS prev ON prev.history_id = curr.previous_history_id
);
