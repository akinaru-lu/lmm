CREATE TABLE IF NOT EXISTS `user` (
	`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(31) NOT NULL,
	`password` VARCHAR(64) NOT NULL,
	`token` VARCHAR(63) NOT NULL,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE `name` (`name`),
	UNIQUE `token` (`token`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `article` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`user` INT UNSIGNED NOT NULL, -- user.id
	`uid` VARCHAR(255) NOT NULL,
	`title` VARCHAR(255) NOT NULL,
	`text` TEXT NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `uid` (`uid`),
	INDEX `created_at` (`created_at`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `article_tag` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`article_uid` VARCHAR(255) NOT NULL, 
	`name` VARCHAR(255) NOT NULL,
	PRIMARY KEY(`id`),
	UNIQUE `article_tag` (`article_uid`, `name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `blog` (
	`id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
	`user` BIGINT UNSIGNED NOT NULL,
	`title` VARCHAR(63) NOT NULL,
	`text` TEXT NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (id),
	UNIQUE `title` (`title`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `category` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(31) NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `name` (`name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `blog_category` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`blog` INT UNSIGNED NOT NULL,
	`category` INT UNSIGNED NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `blog` (`blog`),
	INDEX `category` (`category`, `blog`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `tag` (
	`id` BIGINT unsigned NOT NULL AUTO_INCREMENT,
	`blog` INT unsigned NOT NULL,
	`name` VARCHAR(31) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE `blog_tag` (`blog`, `name`),
	INDEX `tag` (`name`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;

CREATE TABLE IF NOT EXISTS `image` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`uid` VARCHAR(255) NOT NULL,
	`user` BIGINT UNSIGNED NOT NULL,
	`type` TINYINT UNSIGNED NOT NULL DEFAULT 0,
	`created_at` TIMESTAMP NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE `uid` (`uid`),
	INDEX `created_at` (`type`, `created_at`)
) ENGINE = InnoDB DEFAULT CHARACTER SET utf8;
