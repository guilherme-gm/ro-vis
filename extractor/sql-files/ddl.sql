CREATE TABLE `items` (
	`Id` int NOT NULL,
	`Patch` tinytext NOT NULL,
	`_FileVersion` int DEFAULT NULL,
	`IdentifiedName` tinytext DEFAULT NULL,
	`IdentifiedDescription` text DEFAULT NULL,
	`IdentifiedSprite` tinytext DEFAULT NULL,
	`UnidentifiedName` tinytext DEFAULT NULL,
	`UnidentifiedDescription` text DEFAULT NULL,
	`UnidentifiedSprite` tinytext DEFAULT NULL,
	`SlotCount` tinyint DEFAULT NULL,
	`IsBook` tinyint DEFAULT NULL,
	`CanUseBuyingStore` tinyint DEFAULT NULL,
	`CardPrefix` tinytext DEFAULT NULL,
	`CardPostfix` tinytext DEFAULT NULL,
	`CardIllustration` tinytext DEFAULT NULL,
	`ClassNum` int DEFAULT NULL,
	`MoveInfo_canDrop` tinyint DEFAULT NULL,
	`MoveInfo_canTrade` tinyint DEFAULT NULL,
	`MoveInfo_canMoveToStorage` tinyint DEFAULT NULL,
	`MoveInfo_canMoveToCart` tinyint DEFAULT NULL,
	`MoveInfo_canSellToNpc` tinyint DEFAULT NULL,
	`MoveInfo_canMail` tinyint DEFAULT NULL,
	`MoveInfo_canAuction` tinyint DEFAULT NULL,
	`MoveInfo_canMoveToGuildStorage` tinyint DEFAULT NULL,
	`MoveInfo_commentName` tinytext DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `items_history` (
	`HistoryId` varchar(100) NOT NULL,
	`Patch` tinytext NOT NULL,
	`PreviousId` varchar(100) DEFAULT NULL,
	`Id` int NOT NULL,
	`_FileVersion` int DEFAULT NULL,
	`IdentifiedName` tinytext DEFAULT NULL,
	`IdentifiedDescription` text DEFAULT NULL,
	`IdentifiedSprite` tinytext DEFAULT NULL,
	`UnidentifiedName` tinytext DEFAULT NULL,
	`UnidentifiedDescription` text DEFAULT NULL,
	`UnidentifiedSprite` tinytext DEFAULT NULL,
	`SlotCount` tinyint DEFAULT NULL,
	`IsBook` tinyint DEFAULT NULL,
	`CanUseBuyingStore` tinyint DEFAULT NULL,
	`CardPrefix` tinytext DEFAULT NULL,
	`CardPostfix` tinytext DEFAULT NULL,
	`CardIllustration` tinytext DEFAULT NULL,
	`ClassNum` int DEFAULT NULL,
	`MoveInfo_canDrop` tinyint DEFAULT NULL,
	`MoveInfo_canTrade` tinyint DEFAULT NULL,
	`MoveInfo_canMoveToStorage` tinyint DEFAULT NULL,
	`MoveInfo_canMoveToCart` tinyint DEFAULT NULL,
	`MoveInfo_canSellToNpc` tinyint DEFAULT NULL,
	`MoveInfo_canMail` tinyint DEFAULT NULL,
	`MoveInfo_canAuction` tinyint DEFAULT NULL,
	`MoveInfo_canMoveToGuildStorage` tinyint DEFAULT NULL,
	`MoveInfo_commentName` tinytext DEFAULT NULL,
	UNIQUE KEY `HistoryId` (`HistoryId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `quests` (
	`Id` int(11) NOT NULL,
	`Patch` tinytext NOT NULL,
	`_FileVersion` int(11) NOT NULL,
	`Title` tinytext DEFAULT NULL,
	`Description` mediumtext DEFAULT NULL,
	`Summary` tinytext DEFAULT NULL,
	`OldImage` tinytext DEFAULT NULL,
	`IconName` tinytext DEFAULT NULL,
	`NpcSpr` tinytext DEFAULT NULL,
	`NpcNavi` tinytext DEFAULT NULL,
	`NpcPosX` int(11) DEFAULT NULL,
	`NpcPosY` int(11) DEFAULT NULL,
	`RewardEXP` tinytext DEFAULT NULL,
	`RewardJEXP` tinytext DEFAULT NULL,
	`RewardItemList` text DEFAULT NULL,
	`CoolTimeQuest` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `quests_history` (
	`HistoryId` varchar(100) NOT NULL,
	`Patch` tinytext NOT NULL,
	`PreviousId` varchar(100) DEFAULT NULL,
	`Id` int(11) NOT NULL,
	`_FileVersion` int(11) NOT NULL,
	`Title` tinytext DEFAULT NULL,
	`Description` mediumtext DEFAULT NULL,
	`Summary` tinytext DEFAULT NULL,
	`OldImage` tinytext DEFAULT NULL,
	`IconName` tinytext DEFAULT NULL,
	`NpcSpr` tinytext DEFAULT NULL,
	`NpcNavi` tinytext DEFAULT NULL,
	`NpcPosX` int(11) DEFAULT NULL,
	`NpcPosY` int(11) DEFAULT NULL,
	`RewardEXP` tinytext DEFAULT NULL,
	`RewardJEXP` tinytext DEFAULT NULL,
	`RewardItemList` text DEFAULT NULL,
	`CoolTimeQuest` int(11) DEFAULT NULL,
	UNIQUE KEY `HistoryId` (`HistoryId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;

CREATE TABLE `updates` (
	`id` varchar(30) NOT NULL,
	`order` int(11) NOT NULL,
	`updates` mediumtext NOT NULL,
	`patches` mediumtext NOT NULL,
	PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;