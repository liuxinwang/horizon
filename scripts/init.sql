-- 2022.7.16
-- create database
CREATE DATABASE `horizon` DEFAULT CHARACTER SET utf8mb4;
USE horizon;
-- insert admin user
INSERT INTO users(id, user_name, `nick_name`, `password`, `status`, created_at, updated_at)
VALUES('1', 'admin', '管理员', '$2a$10$kfcWYvoUMzY7BtQw/IrMd.aVilOs6.b/xbt560dqR467qz28TpL8K', 'Enabled', NOW(), NOW());
-- insert role
INSERT INTO `roles`(`id`, `name`, `describe`)VALUES('admin', '系统管理员', '');
-- insert admin user role
INSERT INTO `user_roles`(`user_id`, `role_id`)VALUES('1', 'admin');
-- insert inspection metrics
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('BACKUP_STATUS','备份状态','正常/异常',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('BIGTABLE_NUM','大表数量','个',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('CPU_UTILIZATION','CPU使用率','%',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('DEADLOCK','死锁','次',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('DISK_UTILIZATION','磁盘使用率',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('HA_STATUS','HA状态','开/关',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('HIGH_RISK_ACCOUNT','高危账号',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('IBP_CACHE_HITS_RATE','缓存命中率','%',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('INCREMENT_ID_OVERFLOW','自增主键溢出','个',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('IOPS_UTILIZATION','IOPS使用率','%',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('LOCK_WAIT','锁等待','次',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('MEMORY_UTILIZATION','内存使用率','%',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('NETWORK_TRAFFIC_IN','网络流量in','bytes',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('NETWORK_TRAFFIC_OUT','网络流量out','bytes',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('QPS','请求数','个/秒',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('REPLICATION_DELAY','复制延迟','秒',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('REPLICATION_STATUS','复制状态','正常/异常',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('SLOW_SQL_NUM','慢SQL数量','个',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('SWAP_USE','SWAP使用量','MB',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('THREADS_CONNECTED','连接数使用率','%',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('THREADS_RUNNING_NUM','并发线程数','个',now(),now());
INSERT INTO `metrics` (`key`,`name`,`unit`,`created_at`,`updated_at`) VALUES ('TPS','事务数','个/秒',now(),now());
-- insert menus
INSERT INTO `menus` (`id`,`name`,`parent_id`,`meta`,`component`,`redirect`,`path`,`created_at`,`updated_at`,`action_data`,`action_list`) VALUES (1,'Instance',0,'{\"icon\": \"table\", \"show\": true, \"title\": \"menu.instance\"}','RouteView','/instance/list',NULL,'2023-04-17 15:32:49','2023-05-04 16:29:15','[]','[]');
INSERT INTO `menus` (`id`,`name`,`parent_id`,`meta`,`component`,`redirect`,`path`,`created_at`,`updated_at`,`action_data`,`action_list`) VALUES (2,'System',0,'{\"icon\": \"table\", \"show\": true, \"title\": \"menu.system\"}','RouteView','/system/user',NULL,'2023-04-17 15:32:49','2023-05-04 16:29:15','[]','[]');
INSERT INTO `menus` (`id`,`name`,`parent_id`,`meta`,`component`,`redirect`,`path`,`created_at`,`updated_at`,`action_data`,`action_list`) VALUES (1001,'InstanceList',1,'{\"icon\": \"table\", \"show\": true, \"title\": \"menu.instance.list\"}','InstanceList',NULL,'/instance/list','2023-04-17 15:32:49','2023-05-04 16:29:15','[{\"action\": \"add\", \"describe\": \"新增\", \"defaultCheck\": false}, {\"action\": \"query\", \"describe\": \"查询\", \"defaultCheck\": false}, {\"action\": \"edit\", \"describe\": \"修改\", \"defaultCheck\": false}, {\"action\": \"delete\", \"describe\": \"删除\", \"defaultCheck\": false}]','[\"add\", \"query\", \"edit\", \"delete\"]');
INSERT INTO `menus` (`id`,`name`,`parent_id`,`meta`,`component`,`redirect`,`path`,`created_at`,`updated_at`,`action_data`,`action_list`) VALUES (1002,'InstanceInspection',1,'{\"icon\": \"table\", \"show\": true, \"title\": \"menu.instance.inspection\"}','InstanceInspection',NULL,'/instance/inspection','2023-04-17 15:32:49','2023-05-04 16:29:15','[{\"action\": \"get\", \"describe\": \"详情\", \"defaultCheck\": false}, {\"action\": \"query\", \"describe\": \"查询\", \"defaultCheck\": false}]','[\"get\", \"query\"]');
INSERT INTO `menus` (`id`,`name`,`parent_id`,`meta`,`component`,`redirect`,`path`,`created_at`,`updated_at`,`action_data`,`action_list`) VALUES (2001,'SystemUser',2,'{\"icon\": \"table\", \"show\": true, \"title\": \"menu.system.user\"}','SystemUser',NULL,'/system/user','2023-04-17 15:32:49','2023-05-04 16:29:53','[{\"action\": \"add\", \"describe\": \"新增\", \"defaultCheck\": false}, {\"action\": \"query\", \"describe\": \"查询\", \"defaultCheck\": false}, {\"action\": \"edit\", \"describe\": \"修改\", \"defaultCheck\": false}, {\"action\": \"delete\", \"describe\": \"删除\", \"defaultCheck\": false}, {\"action\": \"resetPassword\", \"describe\": \"重置密码\", \"defaultCheck\": false}]','[\"add\", \"query\", \"edit\", \"delete\", \"resetPassword\"]');
INSERT INTO `menus` (`id`,`name`,`parent_id`,`meta`,`component`,`redirect`,`path`,`created_at`,`updated_at`,`action_data`,`action_list`) VALUES (2002,'SystemRole',2,'{\"icon\": \"table\", \"show\": true, \"title\": \"menu.system.role\"}','SystemRole',NULL,'/system/role','2023-04-17 15:32:49','2023-05-04 16:29:53','[{\"action\": \"add\", \"describe\": \"新增\", \"defaultCheck\": false}, {\"action\": \"query\", \"describe\": \"查询\", \"defaultCheck\": false}, {\"action\": \"edit\", \"describe\": \"修改\", \"defaultCheck\": false}, {\"action\": \"delete\", \"describe\": \"删除\", \"defaultCheck\": false}, {\"action\": \"funcPerms\", \"describe\": \"功能权限\", \"defaultCheck\": false}]','[\"add\", \"query\", \"edit\", \"delete\", \"funcPerms\"]');
-- insert admin role permissions (ALL menu)
INSERT IGNORE INTO role_permissions(role_id, menu_id, action_data, action_list)
SELECT 'admin' AS role_id, id AS menu_id, action_data, action_list FROM menus;