
---2017.10.17 新增游戏利润表
CREATE TABLE `profit` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `game_id` int(11) DEFAULT NULL,
  `channel_code` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `date` date DEFAULT NULL,
  `amount` double DEFAULT '0' COMMENT '流水',
  `profit` double DEFAULT '0' COMMENT '利润',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci


