CREATE TABLE `subscription_profile_revisions` (
	`id` int AUTO_INCREMENT NOT NULL,
	`profileId` int NOT NULL,
	`version` int NOT NULL,
	`action` enum('CREATE','UPDATE','DELETE') NOT NULL,
	`snapshot` text NOT NULL,
	`actorId` int,
	`createdAt` timestamp NOT NULL DEFAULT (now()),
	CONSTRAINT `subscription_profile_revisions_id` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `subscription_profiles` (
	`id` int AUTO_INCREMENT NOT NULL,
	`name` varchar(120) NOT NULL,
	`service` enum('FTTH','PPPoE','Hotspot / Voucher','Static IP') NOT NULL,
	`category` enum('Rumahan','Bisnis','Dedicated','Hotspot') NOT NULL,
	`media` enum('Fiber Optic','Wireless','LAN','5G / LTE') NOT NULL,
	`color` varchar(7) NOT NULL DEFAULT '#1677FF',
	`isActive` boolean NOT NULL DEFAULT true,
	`description` text,
	`downloadRate` varchar(32) NOT NULL DEFAULT '1M',
	`uploadRate` varchar(32) NOT NULL DEFAULT '2M',
	`price` decimal(12,2) NOT NULL DEFAULT '0',
	`version` int NOT NULL DEFAULT 1,
	`createdBy` int,
	`updatedBy` int,
	`createdAt` timestamp NOT NULL DEFAULT (now()),
	`updatedAt` timestamp NOT NULL DEFAULT (now()) ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT `subscription_profiles_id` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE INDEX `subscription_profile_revisions_profile_idx` ON `subscription_profile_revisions` (`profileId`);--> statement-breakpoint
CREATE INDEX `subscription_profiles_service_idx` ON `subscription_profiles` (`service`);--> statement-breakpoint
CREATE INDEX `subscription_profiles_status_idx` ON `subscription_profiles` (`isActive`);