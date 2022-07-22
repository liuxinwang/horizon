package model

import "time"

const CpuUtilization string = "CPU_UTILIZATION"
const MemoryUtilization string = "MEMORY_UTILIZATION"
const SwapUse string = "SWAP_USE"
const DiskUtilization string = "DISK_UTILIZATION"
const IOPSUtilization string = "IOPS_UTILIZATION"
const Deadlock string = "DEADLOCK"
const SlowSQLNum string = "SLOW_SQL_NUM"
const IncrementIdOverflow string = "INCREMENT_ID_OVERFLOW"
const LockWait string = "LOCK_WAIT"
const BigTableNum string = "BIGTABLE_NUM"
const ThreadsRunningNum string = "THREADS_RUNNING_NUM"
const ThreadsConnected string = "THREADS_CONNECTED"
const IBPCacheHitsRate string = "IBP_CACHE_HITS_RATE"
const QPS string = "QPS"
const TPS string = "TPS"
const HighRiskAccount string = "HIGH_RISK_ACCOUNT"
const HAStatus string = "HA_STATUS"
const ReplicationStatus string = "REPLICATION_STATUS"
const ReplicationDelay string = "REPLICATION_DELAY"
const BackupStatus string = "BACKUP_STATUS"
const NetworkTrafficIn string = "NETWORK_TRAFFIC_IN"
const NetworkTrafficOut string = "NETWORK_TRAFFIC_OUT"

type Metric struct {
	Key        string       `gorm:"type:varchar(50);primaryKey;comment:指标Key"`
	Name       string       `gorm:"type:varchar(20);not null;comment:名称"`
	Unit       string       `gorm:"type:varchar(20);not null;comment:单位"`
	CreatedAt  time.Time    `gorm:"type:datetime;not null;default:current_timestamp;comment:创建时间"`
	UpdatedAt  time.Time    `gorm:"type:datetime;not null;default:current_timestamp on update current_timestamp;comment:修改时间"`
	InstMetric []InstMetric `gorm:"foreignKey:Metric;references:Key"`
}
