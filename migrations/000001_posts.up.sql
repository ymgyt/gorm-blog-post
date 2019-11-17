BEGIN;


CREATE TABLE IF NOT EXISTS `authors`
(
    `id`      bigint unsigned NOT NULL AUTO_INCREMENT,
    `name`    varchar(100)    NOT NULL,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `reviews`
(
    `id`        bigint unsigned NOT NULL AUTO_INCREMENT,
    `body`      text            NOT NULL,
    `author_id` bigint unsigned NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`)
)
    ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `posts`
(
    `id`           bigint unsigned NOT NULL AUTO_INCREMENT,
    `kind`         int             NOT NULL,
    `title`        varchar(100)    NOT NULL,
    `author_id`    bigint unsigned NOT NULL,
    `post_font`    varchar(100)    NOT NULL,
    `post_theme`   varchar(100)    NOT NULL,
    `published_at` DATETIME,
    `created_at`   DATETIME        NOT NULL,
    `updated_at`   DATETIME        NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`)
)
    ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `contents`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT,
    `post_ref`   bigint unsigned NOT NULL,
    `body`       text            NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`post_ref`) REFERENCES `posts` (`id`)
)
    ENGINE = InnoDB;


CREATE TABLE IF NOT EXISTS `meta`
(
    `id`            bigint unsigned NOT NULL AUTO_INCREMENT,
    `version`       int             NOT NULL,
    `resource_id`   bigint unsigned NOT NULL,
    `resource_type` varchar(100)    NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`resource_id`) REFERENCES `posts` (`id`)
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