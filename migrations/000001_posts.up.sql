BEGIN;

CREATE TABLE IF NOT EXISTS `metas`
(
    `id`      bigint unsigned NOT NULL AUTO_INCREMENT,
    `version` int             NOT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `posts`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `meta_id`    bigint unsigned NOT NULL,
    `content`    text            NOT NULL,
    `created_at` DATETIME        NOT NULL,
    `updated_at` DATETIME        NOT NULL,
    FOREIGN KEY (`meta_id`) REFERENCES `metas` (`id`),
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `authors`
(
    `id`      bigint unsigned NOT NULL AUTO_INCREMENT,
    `name`    varchar(100)    NOT NULL,
    `post_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`)
)
    ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `comments`
(
    `id`      bigint unsigned NOT NULL AUTO_INCREMENT,
    `body`    text            NOT NULL,
    `post_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`)
)
    ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `tags`
(
    `id`   bigint unsigned NOT NULL AUTO_INCREMENT,
    `body` varchar(100),
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `post_tags`
(
    `id`     bigint unsigned NOT NULL AUTO_INCREMENT,
    `tag_id` bigint unsigned NOT NULL,
    `post_id`bigint unsigned NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`),
    FOREIGN KEY (`post_id`) REFERENCES `posts` (`id`)
)
    ENGINE = InnoDB;

COMMIT;