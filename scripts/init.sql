-- 2022.7.16
-- create database
CREATE DATABASE `horizon` DEFAULT CHARACTER SET utf8mb4;
-- insert admin user
INSERT INTO users(id, user_name, `password`, `status`, created_at, updated_at)
VALUES('1', 'admin', '$2a$10$kfcWYvoUMzY7BtQw/IrMd.aVilOs6.b/xbt560dqR467qz28TpL8K', '1', NOW(), NOW());
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
