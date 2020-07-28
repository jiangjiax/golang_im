/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80019
 Source Host           : localhost:3306
 Source Schema         : golang_im

 Target Server Type    : MySQL
 Target Server Version : 80019
 File Encoding         : 65001

 Date: 28/07/2020 16:27:46
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for im_conversation
-- ----------------------------
DROP TABLE IF EXISTS `im_conversation`;
CREATE TABLE `im_conversation`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT COMMENT '会话id',
  `app_id` bigint(0) UNSIGNED NOT NULL COMMENT 'app_id',
  `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '用户id',
  `receiver_id` bigint(0) NOT NULL COMMENT '会话人id',
  `receiver_type` tinyint(0) NOT NULL COMMENT '会话人类型 1 个人 2 群组',
  `disturb` tinyint(0) NOT NULL DEFAULT 1 COMMENT '1 2免打扰',
  `top` tinyint(0) NOT NULL DEFAULT 1 COMMENT '1 2置顶',
  `sender_id` bigint(0) UNSIGNED NULL DEFAULT NULL COMMENT '发送者id',
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `sender_name` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '发送者昵称',
  `status` tinyint(0) NULL DEFAULT 0 COMMENT '0 正常 1 删除',
  `help` bigint(0) NULL DEFAULT 0 COMMENT '客服id',
  `message` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '最新内容',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_app_id_object_seq`(`app_id`, `user_id`, `receiver_id`, `receiver_type`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2540 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '会话表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_device
-- ----------------------------
DROP TABLE IF EXISTS `im_device`;
CREATE TABLE `im_device`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `app_id` bigint(0) UNSIGNED NOT NULL COMMENT 'app_id',
  `user_id` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
  `type` tinyint(0) NOT NULL COMMENT '设备类型：1 Android 2 IOS 3 Windows  4 MacOS 5 Web',
  `brand` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '手机厂商',
  `model` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '机型',
  `system_version` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '系统版本',
  `status` tinyint(0) NOT NULL DEFAULT 0 COMMENT '在线状态：0 离线 1 在线',
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `del` tinyint(0) NOT NULL DEFAULT 1 COMMENT '1 正常；0 删除',
  `identification` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '唯一标识',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_app_id_user_id`(`app_id`, `user_id`, `identification`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2125 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '设备' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_friend
-- ----------------------------
DROP TABLE IF EXISTS `im_friend`;
CREATE TABLE `im_friend`  (
  `id` int(0) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(0) NOT NULL COMMENT 'userid',
  `fid` bigint(0) NOT NULL COMMENT '好友的id',
  `remark` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '备注',
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `way` tinyint(0) NULL DEFAULT NULL COMMENT '1 扫码 2 搜索手机号 3 搜索昵称 4 ID号 5 接受好友请求',
  `examine` tinyint(0) NULL DEFAULT 0 COMMENT '0 等待中 1 同意请求 2 拒绝请求',
  `examinetext` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '好友请求文字',
  `examine_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '请求时间',
  `app_id` bigint(0) NOT NULL DEFAULT 1 COMMENT 'app_id',
  `state` int(0) NULL DEFAULT 1 COMMENT '1 正常 2 删除',
  `is_read` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0 未读 1 已读',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_app_id_group_id_user_id`(`app_id`, `fid`, `user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 129 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '好友表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group
-- ----------------------------
DROP TABLE IF EXISTS `im_group`;
CREATE TABLE `im_group`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `app_id` bigint(0) NOT NULL COMMENT 'app_id',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '群组名称',
  `introduction` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '群组简介',
  `user_num` int(0) NOT NULL DEFAULT 0 COMMENT '群组人数',
  `type` tinyint(0) NOT NULL DEFAULT 1 COMMENT '群组类型，1：小群；2：大群',
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `user_id` int(0) NULL DEFAULT NULL COMMENT '申请人id ',
  `status` int(0) NULL DEFAULT 1 COMMENT '状态 1 审核通过 2 审核中 3审核失败 4 删除或解散',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '群头像',
  `way` tinyint(0) NULL DEFAULT 1 COMMENT '1 直接创建  2 面对面建群',
  `coordinatex` double(100, 6) NULL DEFAULT 0.000000 COMMENT '面对面建群坐标x',
  `coordinatey` double(100, 6) NULL DEFAULT 0.000000 COMMENT '面对面建群坐标y',
  `commandword` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '面对面建群口令',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 86 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '群组' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_group_user
-- ----------------------------
DROP TABLE IF EXISTS `im_group_user`;
CREATE TABLE `im_group_user`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `app_id` bigint(0) NOT NULL DEFAULT 1 COMMENT 'app_id',
  `group_id` bigint(0) UNSIGNED NOT NULL COMMENT '组id',
  `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '用户id',
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `type` tinyint(0) NOT NULL DEFAULT 1 COMMENT '1 群成员 2 群管理 3 群主',
  `label` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '用户在群组的昵称',
  `status` tinyint(0) NOT NULL DEFAULT 0 COMMENT '0 正常 1 删除',
  `examine` tinyint(0) NULL DEFAULT 0 COMMENT '0 等待中 1 同意请求 2 拒绝请求',
  `examinetext` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT '' COMMENT '请求文字',
  `examine_time` datetime(0) NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '请求时间',
  `is_read` bigint(0) NULL DEFAULT 0 COMMENT '群组用户消息索引',
  `friend_read` bigint(0) NOT NULL DEFAULT 0 COMMENT '好友请求 0 未读 1 已读',
  `way` tinyint(0) NULL DEFAULT NULL COMMENT '1 扫码 2 搜索手机号 3 搜索昵称 4 ID号 5 接受加群请求',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_app_id_group_id_user_id`(`app_id`, `group_id`, `user_id`) USING BTREE,
  INDEX `idx_app_id_user_id`(`app_id`, `user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 279 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '群组成员关系' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_message
-- ----------------------------
DROP TABLE IF EXISTS `im_message`;
CREATE TABLE `im_message`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `app_id` int(0) NOT NULL COMMENT 'app_id',
  `object_id` bigint(0) UNSIGNED NOT NULL COMMENT '自己的id',
  `sender_type` tinyint(0) NOT NULL DEFAULT 2 COMMENT '发送者类型 1 系统 2 用户 3 第三方业务系统',
  `sender_id` bigint(0) UNSIGNED NOT NULL COMMENT '发送者id',
  `receiver_type` tinyint(0) NOT NULL COMMENT '接收者类型 1 个人 2 群组',
  `receiver_id` bigint(0) UNSIGNED NOT NULL COMMENT '接收者id，如果是单聊信息，则为user_id，如果是群组消息，则为group_id',
  `to_user_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '需要@的用户id列表，多个用户用，隔开',
  `type` tinyint(0) NOT NULL COMMENT '消息类型',
  `content` varchar(4094) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '消息内容',
  `send_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '消息发送时间',
  `status` tinyint(0) NOT NULL DEFAULT 0 COMMENT '消息状态 0 未处理 1 消息撤回 2 删除',
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `is_read` tinyint(0) NULL DEFAULT 0 COMMENT '0 未读 1 已读',
  `conversation_id` bigint(0) NOT NULL COMMENT '会话id',
  `help` tinyint(0) NULL DEFAULT 1 COMMENT '1 普通消息 2 客服消息',
  `conversation_message` varchar(4094) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT NULL COMMENT '消息内容字符串',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2156 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '消息' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_trends
-- ----------------------------
DROP TABLE IF EXISTS `im_trends`;
CREATE TABLE `im_trends`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(0) NOT NULL COMMENT '发送人id',
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `writing` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '动态文字',
  `imgs` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '动态图片',
  `videos` text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '动态视频',
  `app_id` varchar(50) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT 'app_id 1奶酪',
  `to_user_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '需要@的用户id列表，多个用户用，隔开',
  `status` tinyint(0) NULL DEFAULT 1 COMMENT '1 正常 0 删除',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 106 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '动态表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_trends_comment
-- ----------------------------
DROP TABLE IF EXISTS `im_trends_comment`;
CREATE TABLE `im_trends_comment`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `trends_id` bigint(0) NOT NULL COMMENT '回复的动态id',
  `reply_id` bigint(0) NULL DEFAULT 0 COMMENT '回复人id（看istype，如果是评论动态就存这条动态的发送者id，如果是回复就存被回复的人的id）',
  `comment_id` bigint(0) NOT NULL COMMENT '评论id',
  `app_id` bigint(0) UNSIGNED NOT NULL COMMENT 'app_id 1 奶糖',
  `user_id` bigint(0) NULL DEFAULT NULL COMMENT '发送人id',
  `to_user_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '需要@的用户id列表，多个用户用，隔开',
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `writing` varchar(500) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '评论文字',
  `status` tinyint(0) NULL DEFAULT 0 COMMENT '0 正常 1 删除',
  `istype` tinyint(0) NOT NULL COMMENT '1 评论动态 2 回复',
  `isread` tinyint(0) NULL DEFAULT 0 COMMENT '0 未读 1 已读',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `NK_UTA`(`trends_id`, `user_id`, `app_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 253 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '动态评论表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_trends_handle
-- ----------------------------
DROP TABLE IF EXISTS `im_trends_handle`;
CREATE TABLE `im_trends_handle`  (
  `id` bigint(0) NOT NULL AUTO_INCREMENT,
  `trends_id` bigint(0) NOT NULL COMMENT '动态id',
  `user_id` bigint(0) NULL DEFAULT NULL COMMENT '操作者id',
  `app_id` bigint(0) UNSIGNED NOT NULL COMMENT 'app_id',
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `status` tinyint(0) NULL DEFAULT 1 COMMENT '0 正常 1 删除',
  `istype` tinyint(0) NOT NULL COMMENT '1 点赞 2 转发',
  `reply_id` tinyint(0) NULL DEFAULT 0 COMMENT '发布动态的人id',
  `platform` tinyint(0) NULL DEFAULT 0 COMMENT '1 微博 2 朋友圈 3 qq 4 微信',
  `isread` tinyint(0) NULL DEFAULT 0 COMMENT '0 未读 1 已读',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_tid_uid_aid_istype`(`trends_id`, `user_id`, `app_id`, `istype`, `platform`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 683 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '记录点赞转发收藏的人' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for im_user
-- ----------------------------
DROP TABLE IF EXISTS `im_user`;
CREATE TABLE `im_user`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `app_id` bigint(0) UNSIGNED NOT NULL COMMENT 'app_id',
  `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '用户id',
  `nickname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '昵称',
  `sex` tinyint(0) NOT NULL COMMENT '性别 0 未知 1 男 2 女',
  `avatar_url` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '用户头像链接',
  `extra` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '附加属性',
  `create_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `update_time` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  `sign` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '用户个性签名',
  `status` int(0) NOT NULL DEFAULT 1 COMMENT '1 正常 0 删除',
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '账号',
  `pwd` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `u_app_id_user_id`(`app_id`, `user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 110 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_bin COMMENT = '用户' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
