CREATE DATABASE IF NOT EXISTS `douban_zufang` character set utf8mb4;

USE douban_zufang;

CREATE TABLE IF NOT EXISTS `TopicInfo`
(
    `id`          varchar(20)  NOT NULL,
    `link`        varchar(256) NOT NULL,
    `title`       varchar(100) NOT NULL,
    `create_time` datetime     NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

create user IF NOT EXISTS 'douban'@'%' identified by 'douban@123';

grant all on douban_zufang.* to 'douban'@'%';
