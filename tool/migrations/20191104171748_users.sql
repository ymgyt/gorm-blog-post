-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DROP TABLE IF EXISTS `descriptions`;
CREATE TABLE `descriptions`
(
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `description` text,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`
(
    `id`              bigint unsigned NOT NULL AUTO_INCREMENT,
    `name`            varchar(100)    NOT NULL,
    `has_default`     varchar(100) DEFAULT 'DB-GENERATE-DEFAULT-VALUE',
    `my_scan`         varchar(100),
    `description_id` bigint unsigned NOT NULL,
    `user_meta`            text,
    `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`description_id`) REFERENCES `descriptions` (`id`)
)
    ENGINE = InnoDB;

DROP TABLE IF EXISTS `settings`;
CREATE TABLE `settings`
(
    `id`      bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id` bigint unsigned NOT NULL,
    `lang`    varchar(100) NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES users (`id`)
)
    ENGINE = InnoDB;

DROP TABLE IF EXISTS `profiles`;
CREATE TABLE `profiles`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `user_id`    bigint unsigned NOT NULL,
    `first_name` varchar(100) NOT NULL,
    `last_name`  varchar(100) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES users (`id`)
)
    ENGINE = InnoDB;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `profiles`;