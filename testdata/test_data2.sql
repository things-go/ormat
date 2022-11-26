CREATE TABLE `test_data2` (
    `id` bigint NOT NULL AUTO_INCREMENT,
    `title` varchar(255) NOT NULL COMMENT '标题',
    `content` varchar(2048) NOT NULL COMMENT '内容',
    `value1` float unsigned NOT NULL DEFAULT '1' COMMENT '[@enum:{"0":["none","空","空注释"],"1":["key1","键1","键1注释"],"2":["key2","键2"],"3":["key3","键3"]}]',
    `value2` float(10,1) unsigned NOT NULL DEFAULT '2' COMMENT '值2,[@enum:{"0":["none","空","空注释"],"1":["key1","键1","键1注释"],"2":["key2","键2"],"3":["key3","键3"]}]',
    `value2` double(16,2) NOT NULL DEFAULT '3' COMMENT '值3',
    `value3` enum('00','SH') NOT NULL DEFAULT '4' COMMENT '值4',
	`created_at` datetime NOT NULL COMMENT '发布时间',
	`updated_at` datetime NOT NULL,
    PRIMARY KEY (`id`),
    KEY `uk_title_created_at` (`title`,`created_at`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='公告-面向所有人的消息';