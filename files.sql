CREATE TABLE `file_info` (
     `id` varchar(128) NOT NULL COMMENT '文件ID',
     `key` varchar(512) NOT NULL COMMENT '文件对象Key',
     `filename` varchar(512) DEFAULT NULL COMMENT '原始文件名',
     `file_size` int(11) unsigned DEFAULT NULL COMMENT '文件大小',
     `created_at` datetime NOT NULL COMMENT '创建时间',
     PRIMARY KEY (`id`)
) COMMENT='文件信息';