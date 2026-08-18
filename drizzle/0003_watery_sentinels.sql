CREATE TABLE `administrator_license_policies` (
	`id` int AUTO_INCREMENT NOT NULL,
	`warningDays` int NOT NULL DEFAULT 30,
	`expiringSoonDays` int NOT NULL DEFAULT 7,
	`gracePeriodDays` int NOT NULL DEFAULT 0,
	`updatedBy` int,
	`updatedAt` timestamp NOT NULL DEFAULT (now()) ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT `administrator_license_policies_id` PRIMARY KEY(`id`)
);
