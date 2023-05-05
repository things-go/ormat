CREATE TABLE
    `test_data1` (
        `id` bigint NOT NULL AUTO_INCREMENT,
        `title` varchar(255) NOT NULL COMMENT '标题',
        `content` varchar(2048) DEFAULT NULL COMMENT '内容',
        `value1` float NOT NULL DEFAULT '1' COMMENT '值1,[0:空,1:键1,2:键2,3:键3]',
        `value2` decimal(10, 1) NOT NULL DEFAULT '2' COMMENT '值2,0:空,1:键1,2:键2,3:键3',
        `value3` double (16, 2) NOT NULL DEFAULT '3' COMMENT '值3',
        `value4` enum ('00', 'SH') NOT NULL DEFAULT '00' COMMENT '值4',
        `created_at` datetime NOT NULL COMMENT '发布时间',
        `updated_at` datetime DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `uk_title_created_at` (`title`, `created_at`) USING BTREE
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '公告-面向所有人的消息';