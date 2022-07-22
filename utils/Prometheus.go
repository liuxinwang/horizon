package utils

import (
	"context"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"horizon/config"
	model2 "horizon/model"
	"os"
	"time"
)

type Prom struct {
	Api               v1.API
	StartTime         time.Time
	EndTime           time.Time
	Instance          model2.Instance
	NodeExporterInst  string
	MySQLExporterInst string
}

type MysqlUserResult struct {
	User string
	Host string
}

func PromAPI() v1.API {
	address := fmt.Sprintf("http://%s:%d", config.Conf.Prometheus.Host, config.Conf.Prometheus.Port)
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}
	v1api := v1.NewAPI(client)
	return v1api
}

func (p *Prom) MetricQuery(query string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _, err := p.Api.Query(ctx, query, time.Now())
	resultByte, _ := json.Marshal(result)
	return resultByte, err
}

func (p *Prom) MetricQueryRange(query string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := v1.Range{
		Start: p.StartTime,
		End:   p.EndTime,
		Step:  time.Minute,
	}
	result, _, err := p.Api.QueryRange(ctx, query, r)
	resultByte, _ := json.Marshal(result)
	return resultByte, err
}

func (p *Prom) MetricCPU() ([]byte, error) {
	query := fmt.Sprintf("100*(1-avg(irate(node_cpu_seconds_total{instance='%s', mode='idle'}[1m])) without (cpu)) ", p.NodeExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricMemory() ([]byte, error) {
	query := fmt.Sprintf("(1 - (node_memory_MemAvailable_bytes{instance=~'%s'} / (node_memory_MemTotal_bytes{instance=~'%s'})))* 100"+
		" or (1 - ((node_memory_Buffers_bytes{instance=~'%s'} + node_memory_Cached_bytes{instance=~'%s'} + node_memory_MemFree_bytes{instance=~'%s'})"+
		" / (node_memory_MemTotal_bytes{instance=~'%s'})))* 100",
		p.NodeExporterInst, p.NodeExporterInst, p.NodeExporterInst, p.NodeExporterInst, p.NodeExporterInst, p.NodeExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricSwap() ([]byte, error) {
	query := fmt.Sprintf("(node_memory_SwapTotal_bytes{instance=~'%s'} - node_memory_SwapFree_bytes{instance=~'%s'}) / 1024 / 1024",
		p.NodeExporterInst, p.NodeExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricDisk() ([]byte, error) {
	query := fmt.Sprintf("100 - ((node_filesystem_avail_bytes{instance=~'%s',mountpoint='/',fstype=~'ext4|xfs'} * 100)"+
		" / node_filesystem_size_bytes{instance=~'%s',mountpoint='/',fstype=~'ext4|xfs'})",
		p.NodeExporterInst, p.NodeExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricIOPS() ([]byte, error) {
	query := fmt.Sprintf("irate(node_disk_reads_completed_total{instance=~'%s'}[1m])"+
		" + irate(node_disk_writes_completed_total{instance=~'%s'}[1m])",
		p.NodeExporterInst, p.NodeExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricDeadlock() ([]byte, error) {
	query := fmt.Sprintf("(mysql_global_status_innodb_deadlocks{instance='%s'}) or "+
		"(mysql_info_schema_innodb_metrics_lock_lock_deadlocks_total{instance='%s'})",
		p.MySQLExporterInst, p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricSlowSQLNum() ([]byte, error) {
	query := fmt.Sprintf("increase(mysql_global_status_slow_queries{instance='%s'}[30s])",
		p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricIncrementIdOverflow() ([]byte, error) {
	query := fmt.Sprintf("(1-mysql_info_schema_auto_increment_column{instance='%s',schema!~'test|mysql'} / "+
		"mysql_info_schema_auto_increment_column_max{instance='%s',schema!~'test|mysql'})*100 < 20",
		p.MySQLExporterInst, p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricLockWait() ([]byte, error) {
	query := fmt.Sprintf("increase(mysql_global_status_innodb_row_lock_current_waits{instance='%s'}[1m])",
		p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricBigTableNum() ([]byte, error) {
	query := fmt.Sprintf("increase(mysql_global_status_innodb_row_lock_current_waits{instance='%s'}[1m])",
		p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricThreadsRunningNum() ([]byte, error) {
	query := fmt.Sprintf("mysql_global_status_threads_running{instance='%s'}",
		p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricThreadsConnected() ([]byte, error) {
	query := fmt.Sprintf("(mysql_global_status_threads_connected{instance='%s'}"+
		" / mysql_global_variables_max_connections{instance='%s'}) * 100",
		p.MySQLExporterInst, p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricIBPCacheHitsRate() ([]byte, error) {
	query := fmt.Sprintf("(increase(mysql_global_status_innodb_buffer_pool_read_requests{instance='%s'}[1m]))"+
		" / (increase(mysql_global_status_innodb_buffer_pool_reads{instance='%s'}[1m])"+
		" + increase(mysql_global_status_innodb_buffer_pool_read_requests{instance='%s'}[1m])) * 100",
		p.MySQLExporterInst, p.MySQLExporterInst, p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricQPS() ([]byte, error) {
	query := fmt.Sprintf("irate(mysql_global_status_queries{instance='%s'}[1m])",
		p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricTPS() ([]byte, error) {
	query := fmt.Sprintf("sum(rate(mysql_global_status_commands_total{instance='%s',"+
		"command=~'(commit|rollback)'}[1m])) without (command)",
		p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricHighRiskAccount() ([]byte, error) {
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=3s"
	dsn = fmt.Sprintf(dsn, p.Instance.User, DecryptAES([]byte(config.Conf.General.SecretKey), p.Instance.Password), p.Instance.Ip, p.Instance.Port, "mysql")
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	weakPasswordList := []string{
		"", // ''
		"*81F5E21E35407D884A6CD4A731AEBFB6AF209E1B", // root
		"*6BB4837EB74329105EE4568DDA7DC67ED2CA2AD9", // 123456
		"*84AAC12F54AB666ECFC2A83C676908C8BBC381B1", // 12345678
		"*CC67043C7BCFF5EEA5566BD9B1F3C74FD9A5CF5D", // 123456789
		"*6B5EDDE567F4F29018862811195DBD14B8ADDD2A", // 1234567890
		"*B6E7D9CB4385CA81E24FF70D5705954B78AD583B", // 0123456789
	}
	var mysqlUsersResult []MysqlUserResult
	result := db.Table("user").Select("user", "host").Where(
		"host = '%' OR authentication_string IN ?", weakPasswordList).Find(&mysqlUsersResult)
	fmt.Printf(string(result.RowsAffected))
	resultJson, _ := json.Marshal(mysqlUsersResult)
	return resultJson, result.Error
}

func (p *Prom) MetricHAStatus() ([]byte, error) {
	// TODO
	return nil, nil
}

func (p *Prom) MetricReplicationStatus() ([]byte, error) {
	query := fmt.Sprintf("mysql_slave_status_slave_io_running{instance='%s'}"+
		" + mysql_slave_status_slave_sql_running{instance='%s'}",
		p.MySQLExporterInst, p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricReplicationDelay() ([]byte, error) {
	query := fmt.Sprintf("mysql_slave_status_seconds_behind_master{instance='%s'}",
		p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricBackupStatus() ([]byte, error) {
	// TODO
	return nil, nil
}

func (p *Prom) MetricNetworkTrafficIn() ([]byte, error) {
	query := fmt.Sprintf("irate(mysql_global_status_bytes_received{instance='%s'}[1m])",
		p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}

func (p *Prom) MetricNetworkTrafficOut() ([]byte, error) {
	query := fmt.Sprintf("irate(mysql_global_status_bytes_sent{instance='%s'}[1m])",
		p.MySQLExporterInst)
	return p.MetricQueryRange(query)
}
