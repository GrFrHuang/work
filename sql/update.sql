-- 2017-11-08 添加游戏停运表
CREATE TABLE `game_outage` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `game_id` int(11) DEFAULT NULL,
  `incr_time` bigint(20) DEFAULT NULL COMMENT '关新增时间',
  `recharge_time` bigint(20) DEFAULT NULL COMMENT '关充值时间',
  `server_time` bigint(20) DEFAULT NULL COMMENT '关服务器时间',
  `create_time` bigint(20) DEFAULT NULL COMMENT '提交时间',
  `create_person` int(11) DEFAULT NULL COMMENT '创建人',
  `desc` text COLLATE utf8_unicode_ci COMMENT '备注',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci


#   2017-11-13 修改合同上传文件为多个
ALTER TABLE `work_together`.`contract`     CHANGE `file_id` `file_id` VARCHAR(255) DEFAULT '0' NULL ;

UPDATE contract ,
  (SELECT id,CONCAT("[",file_id,"]") AS file_id FROM contract WHERE file_id != "0") a
SET contract.file_id=a.file_id
WHERE  contract.id=a.id

UPDATE contract SET file_id="" WHERE file_id ="0"

-- 2017-11-22 添加游戏停运表
ALTER TABLE distribution_company ADD  yunduan_responsible_person int(11) DEFAULT NULL COMMENT '云端负责人'
ALTER TABLE distribution_company ADD  youliang_responsible_person int(11) DEFAULT NULL COMMENT '有量负责人'

-- 2018-11-28 主合同相关
ALTER TABLE `work_together`.`contract`     ADD COLUMN `is_main` INT DEFAULT '2' NULL COMMENT '是否在主合同中，1:是，2:否' AFTER `effective_state`;




---2017.12.25 新增用户人脸信息表
CREATE TABLE `user_face` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(10) DEFAULT NULL COMMENT '用户id',
  `status` int(2) DEFAULT '3' COMMENT '人脸状态（1绑定成功，2正在验证,3未绑定,4未通过）',
  `path` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '注册时的人脸截图',
  `remarks` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '备注',
  `create_time` bigint(20) DEFAULT NULL COMMENT '创建时间',
  `update_time` bigint(20) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


---2017.12.25 新增用户签到信息表
CREATE TABLE `user_sign` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL COMMENT '用户ID',
  `status` int(11) DEFAULT NULL COMMENT '1已经审核2未审核3拒绝通过',
  `tag` int(11) DEFAULT NULL COMMENT '1（早会签到） 2其他签到',
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `path` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `date` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '签到年月日',
  `remarks` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '备注',
  `sign_time` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '签到时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


# 2018-01-02
SHOW VARIABLES LIKE '%table_open_cache%';
SHOW VARIABLES LIKE '%table_definition_cache%';
SET GLOBAL table_open_cache=16384;
SET GLOBAL table_definition_cache=16384;



CREATE TABLE `user_face` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(10) DEFAULT NULL COMMENT '用户id',
  `status` int(2) DEFAULT '3' COMMENT '人脸状态（1绑定成功，2正在验证,3未绑定,4未通过）',
  `path` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '注册时的人脸截图',
  `name` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '姓名',
  `remarks` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '备注',
  `create_time` bigint(20) DEFAULT NULL COMMENT '创建时间',
  `update_time` bigint(20) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=38 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



CREATE TABLE `user_face_verify` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL COMMENT '用户ID',
  `tag` int(11) DEFAULT NULL COMMENT '1打款 2二维码登录 3登录',
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '用户名',
  `image` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '人脸记录',
  `remarks` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '备注',
  `create_time` bigint(20) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=82 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



CREATE TABLE `user_sign` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL COMMENT '用户ID',
  `status` int(11) DEFAULT NULL COMMENT '1已经审核2未审核3拒绝通过',
  `tag` int(11) DEFAULT NULL COMMENT '1（早会签到） 2其他签到',
  `name` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL,
  `path` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL,
  `date` varchar(50) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '签到年月日',
  `remarks` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '备注',
  `sign_time` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '签到时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=44 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

# 2018-01-08
ALTER TABLE `work_together`.`game`     ADD COLUMN `source` INT(11) NULL COMMENT '游戏来源' AFTER `game_name`;

# 2018-01-16
CREATE TABLE `remit_down_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `remit_down_id` int(11) NOT NULL COMMENT '回款表(remit_down_account)的id',
  `verify_channel_id` int(11) NOT NULL COMMENT '对账单(verify_channel)的id',
  `remit_month` varchar(7) COLLATE utf8_unicode_ci NOT NULL COMMENT '回款月份',
  `remit_money` decimal(16,2) NOT NULL DEFAULT '0.00' COMMENT '回款金额',
  `remit_type` int(11) NOT NULL DEFAULT '1' COMMENT '回款类型,1:全部,2:部分',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci



# 工作流表

CREATE TABLE `workflow_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `task_id` int(11) DEFAULT NULL,
  `node_id` int(20) DEFAULT NULL COMMENT '结点ID',
  `user_id` int(20) DEFAULT NULL COMMENT '用户ID',
  `user_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '姓名',
  `workflow_id` int(20) DEFAULT NULL COMMENT '任务流ID',
  `remarks` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '备注（操作信息）',
  `create_time` bigint(30) DEFAULT NULL COMMENT '创建时间',
  `update_time` bigint(30) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=208 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


CREATE TABLE `workflow_name` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `tb_nickname` varchar(100) CHARACTER SET utf8 DEFAULT NULL COMMENT '任务别名',
  `tb_name` varchar(50) CHARACTER SET utf8 DEFAULT NULL COMMENT '工作流表名',
  `node_ids` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '流程的结点ID',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of workflow_name
-- ----------------------------
INSERT INTO `workflow_name` VALUES ('1', '发包流程', 'workflow_send_package', '2,3,4,5,6,7');



CREATE TABLE `workflow_node` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `node_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '结点名称',
  `department_id` int(20) DEFAULT NULL COMMENT '角色ID',
  `wf_name_id` int(10) DEFAULT NULL COMMENT '那种工作流',
  `department_name` varchar(100) CHARACTER SET utf8 DEFAULT NULL COMMENT '角色信息',
  `permissions` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '权限（eg:1,2,3,4  |增,删,改,查）',
  `rollback` int(2) DEFAULT NULL COMMENT '1 能回退,2不行 ',
  `current_node_id` varchar(100) CHARACTER SET utf8 DEFAULT NULL COMMENT '上一步结点，号分割',
  `next_node_id` varchar(100) CHARACTER SET utf8 DEFAULT NULL COMMENT '下一步结点，号分割',
  `rollback_id` int(10) DEFAULT NULL COMMENT '能回退到的结点',
  `workflow_show` int(2) DEFAULT NULL COMMENT ' 1显示渠道信息 2显示合同信息 ',
  `channel_proportion_hide` int(2) DEFAULT NULL COMMENT ' 1显示分成信息',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- ----------------------------
-- Records of workflow_node
-- ----------------------------
INSERT INTO `workflow_node` VALUES ('2', '商务渠道信息', '237', '1', '商务', '1,2,3,4', '2', '', '3', null, '1', '1');
INSERT INTO `workflow_node` VALUES ('3', '运营部审核', '10', '1', '运营', '4', '1', '2', '4', '2', '1', '0');
INSERT INTO `workflow_node` VALUES ('4', '财务部审核', '7', '1', '财务', '1,2,3,4', '1', '3', '5', '2', '1', '1');
INSERT INTO `workflow_node` VALUES ('5', '结算部审核', '6', '1', '结算', '1,2,3,4', '1', '4', '6', '2', '2', '1');
INSERT INTO `workflow_node` VALUES ('6', '法务部审核', '239', '1', '法务', '1,2,3,4', '1', '5', '7', '2', '2', '1');
INSERT INTO `workflow_node` VALUES ('7', '客服部出包', '244', '1', '客服', '1,2,3,4', '1', '6', '2', '2', '1', '0');


CREATE TABLE `workflow_send_package` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `channel_info_id` int(10) DEFAULT NULL COMMENT '渠道信息ID',
  `contract_info_id` int(10) DEFAULT NULL COMMENT '合同信息ID',
  `user_id` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '用户操作（1,2,3）顺序执行',
  `update_time` bigint(20) DEFAULT NULL,
  `task_id` int(10) DEFAULT NULL COMMENT '任务ID',
  `channel_info_hide` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '渠道信息权限（结点ID）',
  `contract_info_hide` varchar(20) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '合同信息权限（结点ID）',
  `create_time` bigint(20) DEFAULT NULL COMMENT '创建时间',
  `current_user_ids` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '当前操作的用户id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=79 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



CREATE TABLE `workflow_task` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `task_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '任务名称',
  `wf_name_id` int(20) DEFAULT NULL COMMENT '工作流ID',
  `current_progress` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '当前的进度',
  `current_progress_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '当前进度结点名',
  `create_time` bigint(30) DEFAULT NULL COMMENT '创建时间',
  `create_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '创建人的姓名',
  `update_time` bigint(30) DEFAULT NULL COMMENT '修改时间',
  `update_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '更新人的名字',
  `status` int(6) DEFAULT NULL COMMENT '状态（1成功 2失败 3挂起 4其他）',
  `remarks` varchar(100) CHARACTER SET utf8 DEFAULT NULL COMMENT '备注信息',
  `current_success` int(10) DEFAULT NULL,
  `channel_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '渠道名称',
  `game_name` varchar(100) COLLATE utf8_unicode_ci DEFAULT NULL COMMENT '游戏名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=103 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


ALTER TABLE `workflow_node`
ADD COLUMN `step`  int(5) NULL COMMENT '步骤' AFTER `channel_proportion_hide`;

# 工作流表结束


//新增小绵羊任务查看权限
INSERT INTO `work_together`.`permission` (`name`, `type`, `model`, `methods`, `field`, `condition`, `readonly`) VALUES ( '查看小绵羊任务', '2', 'small_sheep', '[\"查\"]', NULL, NULL, '1');
INSERT INTO `work_together`.`permission` (`name`, `type`, `model`, `methods`, `field`, `condition`, `readonly`) VALUES ( '编辑小绵羊任务', '2', 'small_sheep', '[\"改\",\"删\",\"增\"]', NULL, NULL, '1');


-----------------2018-2-9-------------

ALTER TABLE `channel_access`
ADD COLUMN `accessory`  varchar(255) NULL

ALTER TABLE `contract`
ADD COLUMN `accessory`  varchar(255) NULL
ADD COLUMN `accessory`  varchar(255) NULL


--------------------2018-02-23-------------------------------
ALTER TABLE `channel_access`
ADD COLUMN `pactId`  int NULL AFTER `accessory`;

ALTER TABLE `channel_access`
MODIFY COLUMN `pactId`  int(11) NULL DEFAULT NULL COMMENT '合同ID' AFTER `accessory`;


-----------2018-02-24--------------------
CREATE TABLE `user_online_logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL,
  `online_status` tinyint(4) DEFAULT NULL,
  `create_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE user_desktop_client
(
  id            INT AUTO_INCREMENT
    PRIMARY KEY,
  uid           INT      NULL,
  ctoken        CHAR(32) NULL,
  online_status TINYINT  NULL,
  create_time   INT      NULL,
  CONSTRAINT user_desktop_client_uid_uindex
  UNIQUE (uid),
  CONSTRAINT user_desktop_client_ctoken_uindex
  UNIQUE (ctoken)
) ENGINE = InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE user_desktop_screen_log
(
  id          INT AUTO_INCREMENT
    PRIMARY KEY,
  uid         INT          NULL,
  imageurl    VARCHAR(255) NULL,
  create_time INT          NULL
) ENGINE = InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


ALTER TABLE `user_desktop_screen_log`
ADD COLUMN `cycles` int(10) NULL AFTER `imageurl`;

ALTER TABLE `user_desktop_screen_log`
MODIFY COLUMN `cycles` int(10) NULL DEFAULT 1 AFTER `imageurl`;


ALTER TABLE `user`
ADD COLUMN `position`  varchar(20) NULL DEFAULT '' COMMENT '职位' AFTER `pic_file_id`;

ALTER TABLE `user`
ADD COLUMN `entry_time`  int(11) NULL DEFAULT 0 COMMENT '入职时间' AFTER `position`;

ALTER TABLE `channel_access`
MODIFY COLUMN `game_id` bigint(11) NOT NULL AFTER `id`;

ALTER TABLE `channel_access`
MODIFY COLUMN `pactId` varchar(11) NULL DEFAULT NULL COMMENT '合同ID' AFTER `accessory`;

CREATE TABLE `kpi_task` (
`id`  int(11) NOT NULL AUTO_INCREMENT ,
`task_complete_date`  int(11) NULL DEFAULT 0 COMMENT '任务完成时间' ,
`task_publish_date`  int(11) NULL DEFAULT 0 COMMENT '发布时间' ,
`assesseser_id`  int(8) NULL DEFAULT 0 COMMENT '考核人' ,
`publisher_id`  int(8) NULL DEFAULT 0 COMMENT '发布人' ,
`publish_state`  tinyint(2) NULL DEFAULT 1 COMMENT '发布状态 0.未发布 1.已发布' ,
`audit_state`  tinyint(2) NULL DEFAULT 0 COMMENT '审核状态 0.未审核 1.已审核' ,
PRIMARY KEY (`id`)
);

ALTER TABLE `kpi_task`
ADD COLUMN `create_time`  int(11) NULL AFTER `audit_state`,
ADD COLUMN `update_time`  int(11) NULL AFTER `create_time`;

CREATE TABLE `kpi_child_task` (
`id`  int(11) NOT NULL AUTO_INCREMENT ,
`task_name`  varchar(50) NULL DEFAULT '' COMMENT '任务名' ,
`period`  varchar(20) NULL DEFAULT '' COMMENT '任务时长' ,
`progress_rate`  decimal(4,2) NULL DEFAULT 0.00 COMMENT '进度' ,
`remark`  varchar(255) NULL DEFAULT '' COMMENT '备注' ,
`score`  decimal(4,2) NULL DEFAULT 0.00 COMMENT '分数' ,
`annotations`  varchar(255) NULL DEFAULT '' COMMENT '领导批注' ,
`create_time`  int(11) NULL COMMENT '创建时间' ,
`update_time`  int(11) NULL COMMENT '修改时间' ,
PRIMARY KEY (`id`)
);

ALTER TABLE `kpi_child_task`
ADD COLUMN `task_id`  int(11) NULL DEFAULT 0 COMMENT '父级任务' AFTER `update_time`;

ALTER TABLE `kpi_child_task`
ADD INDEX `task_id_index` (`task_id`) USING BTREE;

ALTER TABLE `kpi_task`
ADD COLUMN `publish_type`  tinyint(2) NULL DEFAULT 0 COMMENT '发布类型 0.立刻 1.定时' AFTER `update_time`;

ALTER TABLE `kpi_task`
ADD COLUMN `department_id`  int(4) NULL COMMENT '部门' AFTER `publish_type`;

ALTER TABLE `kpi_child_task`
ADD COLUMN `flag`  tinyint(2) NULL DEFAULT 0 COMMENT '0.被动任务  1.主动任务' AFTER `task_id`;