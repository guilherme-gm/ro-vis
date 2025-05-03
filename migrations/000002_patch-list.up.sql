CREATE TABLE `patches` (
	`id` int(11) NOT NULL AUTO_INCREMENT,
	`name` varchar(255) NOT NULL,
	`date` datetime NOT NULL,
	`files` json NOT NULL,
	PRIMARY KEY (`id`)
);
