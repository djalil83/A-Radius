CREATE TABLE `administrator_license_audit` (
	`id` int AUTO_INCREMENT NOT NULL,
	`licenseId` int NOT NULL,
	`actorId` int NOT NULL,
	`action` enum('CREATE','EDIT','EXTEND','SUSPEND','ACTIVATE') NOT NULL,
	`snapshot` text NOT NULL,
	`createdAt` timestamp NOT NULL DEFAULT (now()),
	CONSTRAINT `administrator_license_audit_id` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `administrator_licenses` (
	`id` int AUTO_INCREMENT NOT NULL,
	`administratorId` int NOT NULL,
	`branchId` varchar(64) NOT NULL,
	`companyName` varchar(180) NOT NULL,
	`packageName` varchar(100) NOT NULL,
	`status` enum('ACTIVE','SUSPENDED','EXPIRED','TRIAL') NOT NULL DEFAULT 'ACTIVE',
	`startDate` timestamp NOT NULL,
	`endDate` timestamp NOT NULL,
	`activePeriodDays` int NOT NULL DEFAULT 365,
	`price` decimal(14,2) NOT NULL DEFAULT '0',
	`currency` varchar(3) NOT NULL DEFAULT 'IDR',
	`autoRenewal` boolean NOT NULL DEFAULT false,
	`dueDate` timestamp NOT NULL,
	`maxCustomers` int NOT NULL DEFAULT 0,
	`maxUsers` int NOT NULL DEFAULT 0,
	`maxRouters` int NOT NULL DEFAULT 0,
	`maxOlt` int NOT NULL DEFAULT 0,
	`maxOdp` int NOT NULL DEFAULT 0,
	`maxVouchers` int NOT NULL DEFAULT 0,
	`maxTechnicians` int NOT NULL DEFAULT 0,
	`maxPartners` int NOT NULL DEFAULT 0,
	`storageLimitGb` int NOT NULL DEFAULT 0,
	`usedCustomers` int NOT NULL DEFAULT 0,
	`usedUsers` int NOT NULL DEFAULT 0,
	`usedRouters` int NOT NULL DEFAULT 0,
	`usedOlt` int NOT NULL DEFAULT 0,
	`usedStorageGb` int NOT NULL DEFAULT 0,
	`features` text NOT NULL,
	`createdAt` timestamp NOT NULL DEFAULT (now()),
	`updatedAt` timestamp NOT NULL DEFAULT (now()) ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT `administrator_licenses_id` PRIMARY KEY(`id`),
	CONSTRAINT `administrator_licenses_administratorId_unique` UNIQUE(`administratorId`)
);
--> statement-breakpoint
CREATE INDEX `administrator_license_audit_license_idx` ON `administrator_license_audit` (`licenseId`);