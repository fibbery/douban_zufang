CREATE DATABASE `douban_zufang` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE douban_zufang;

CREATE TABLE `TopicInfo`
(
    `id`          varchar(20)  NOT NULL,
    `link`        varchar(256) NOT NULL,
    `title`       varchar(100) NOT NULL,
    `create_time` datetime     NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

create user 'douban'@'%' identified by 'douban@123';

grant all on douban_zufang.* to 'douban'@'%';
