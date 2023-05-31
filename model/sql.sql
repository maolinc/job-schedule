/*Table structure for table `job_record` */

DROP TABLE IF EXISTS `job_record`;

CREATE TABLE `job_record` (
                              `id` bigint NOT NULL AUTO_INCREMENT,
                              `job_id` bigint NOT NULL COMMENT '任务id',
                              `start_time` datetime DEFAULT NULL COMMENT '开始时间',
                              `end_time` datetime DEFAULT NULL COMMENT '结束时间',
                              `result` varchar(5000) DEFAULT NULL COMMENT '执行结果完整信息json格式',
                              `status` varchar(20) DEFAULT NULL COMMENT '结果状态，ok | error',
                              `use_milli` bigint DEFAULT '0' COMMENT '耗时,毫秒',
                              `exec_type` varchar(20) DEFAULT NULL COMMENT '执行类型，action | compensate',
                              PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='执行记录表';

/*Data for the table `job_record` */

insert  into `job_record`(`id`,`job_id`,`start_time`,`end_time`,`result`,`status`,`use_milli`,`exec_type`) values
                                                                                                               (1,1,'2023-05-31 17:31:51','2023-05-31 17:31:51','{\"data\":\"处理2023-05-30 17:31:51至2023-05-31 17:31:51登录日志：循环100次，共处理1000000条\",\"key\":\"LoginLogJob\",\"params\":\"\",\"error\":\"\"}','ok',34,'action'),
                                                                                                               (2,1,'2023-05-31 17:31:56','2023-05-31 17:31:56','{\"data\":\"处理2023-05-30 17:31:56至2023-05-31 17:31:56登录日志：循环100次，共处理1000000条\",\"key\":\"LoginLogJob\",\"params\":\"\",\"error\":\"\"}','ok',36,'action'),
                                                                                                               (3,1,'2023-05-31 17:32:01','2023-05-31 17:32:01','{\"data\":\"处理2023-05-30 17:32:01至2023-05-31 17:32:01登录日志：循环100次，共处理1000000条\",\"key\":\"LoginLogJob\",\"params\":\"\",\"error\":\"\"}','ok',35,'action'),
                                                                                                               (4,1,'2023-05-31 17:32:06','2023-05-31 17:32:06','{\"data\":\"处理2023-05-30 17:32:06至2023-05-31 17:32:06登录日志：循环100次，共处理1000000条\",\"key\":\"LoginLogJob\",\"params\":\"\",\"error\":\"\"}','ok',31,'action'),
                                                                                                               (5,1,'2023-05-31 17:32:11','2023-05-31 17:32:11','{\"data\":\"处理2023-05-30 17:32:11至2023-05-31 17:32:11登录日志：循环100次，共处理1000000条\",\"key\":\"LoginLogJob\",\"params\":\"\",\"error\":\"\"}','ok',34,'action'),
                                                                                                               (6,2,'2023-05-31 17:32:31','2023-05-31 17:32:31','{\"data\":\"同步2023-05-30 17:32:31至2023-05-31 17:32:31的订单：共处理20000条订单\",\"key\":\"SyncOrderJob\",\"params\":\"\",\"error\":\"\"}','ok',33,'action'),
                                                                                                               (7,2,'2023-05-31 17:32:41','2023-05-31 17:32:41','{\"data\":\"同步2023-05-30 17:32:41至2023-05-31 17:32:41的订单：共处理20000条订单\",\"key\":\"SyncOrderJob\",\"params\":\"\",\"error\":\"\"}','ok',49,'action'),
                                                                                                               (8,2,'2023-05-31 17:32:51','2023-05-31 17:32:51','{\"data\":\"同步2023-05-30 17:32:51至2023-05-31 17:32:51的订单：共处理20000条订单\",\"key\":\"SyncOrderJob\",\"params\":\"\",\"error\":\"\"}','ok',35,'action'),
                                                                                                               (9,2,'2023-05-31 17:34:40','2023-05-31 17:34:40','{\"data\":\"同步2023-04-08 00:00:00至2023-04-08 23:59:59的订单：共处理20000条订单\",\"key\":\"SyncOrderJob\",\"params\":{\"endTime\":\"2023-04-08 23:59:59\",\"startTime\":\"2023-04-08 00:00:00\"},\"error\":\"\"}','ok',31,'compensate');

/*Table structure for table `job_schedule` */

DROP TABLE IF EXISTS `job_schedule`;

CREATE TABLE `job_schedule` (
                                `id` bigint NOT NULL AUTO_INCREMENT,
                                `job_name` varchar(128) DEFAULT '' COMMENT '任务名称',
                                `des` varchar(500) DEFAULT '' COMMENT '描述',
                                `cron` varchar(50) DEFAULT NULL COMMENT '调度时间，cron表达式',
                                `execute_count` bigint DEFAULT '0' COMMENT '调度次数',
                                `fail_count` bigint DEFAULT '0' COMMENT '失败次数',
                                `now_status` varchar(20) DEFAULT 'wait' COMMENT '当前状态，wait | readly | running',
                                `job_status` varchar(20) DEFAULT 'stop' COMMENT '任务转态, stop不加入调度 | enable加入调度中',
                                `key` varchar(128) DEFAULT '' COMMENT '与程序中的任务key一致，唯一',
                                `parallel` varchar(20) DEFAULT 'reject' COMMENT '是否允许并发执行，allow | reject',
                                `lock` int DEFAULT '0' COMMENT '锁，0解锁, 1上锁',
                                `delete_flag` int DEFAULT '0',
                                `create_time` datetime DEFAULT NULL,
                                `update_time` datetime DEFAULT NULL,
                                PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='调度任务表';

/*Data for the table `job_schedule` */

insert  into `job_schedule`(`id`,`job_name`,`des`,`cron`,`execute_count`,`fail_count`,`now_status`,`job_status`,`key`,`parallel`,`lock`,`delete_flag`,`create_time`,`update_time`) values
                                                                                                                                                                                       (1,'计算用户留存','补偿参数格式：{startTime: string, endTime: string}\n参数描述: startTime: 开始时间,格式YYYY-MM-DD hh:mm:ss; endTime: 结束时间,格式YYYY-MM-DD hh:mm:ss','1/5 * * * * *',5,0,'readly','stop','LoginLogJob','reject',0,0,'2023-05-31 17:30:39','2023-05-31 17:32:13'),
                                                                                                                                                                                       (2,'同步每天订单','补偿参数格式：{startTime: string, endTime: string}\n参数描述: startTime: 开始时间,格式YYYY-MM-DD hh:mm:ss; endTime: 结束时间,格式YYYY-MM-DD hh:mm:ss','1/10 * * * * *',3,0,'readly','stop','SyncOrderJob','reject',0,0,'2023-05-31 17:31:22','2023-05-31 17:32:58');
