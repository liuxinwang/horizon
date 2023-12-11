package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"horizon/config"
	"horizon/model"
	"horizon/utils"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// DataMigrateJobSelectByList 查询列表
func DataMigrateJobSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	var Db = model.Db
	var dataMigrateJobs []model.DataMigrateJob
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.DataMigrateJob{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 查询条件处理
	if name, isExist := c.GetQuery("name"); isExist == true && strings.Trim(name, " ") != "" {
		Db = Db.Where("name like ?", "%"+name+"%")
	}
	if status, isExist := c.GetQuery("Status"); isExist == true && status != "0" {
		Db = Db.Where("status = ?", status)
	}

	// 执行查询
	Db.Order("created_at desc").Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&dataMigrateJobs)
	Db.Model(&model.DataMigrateJob{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &dataMigrateJobs, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// DataMigrateJobSelectById 查看信息
func DataMigrateJobSelectById(c *gin.Context) {
	type tmpData struct {
		SourceInstance model.Instance `json:"sourceInstance"`
		TargetInstance model.Instance `json:"targetInstance"`
	}
	var td tmpData
	var dataMigrateJob model.DataMigrateJob
	model.Db.Where("id = ?", c.Param("id")).First(&dataMigrateJob)
	model.Db.Where("inst_id = ?", dataMigrateJob.SourceInstId).First(&td.SourceInstance)
	model.Db.Where("inst_id = ?", dataMigrateJob.TargetInstId).First(&td.TargetInstance)
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": &dataMigrateJob, "extData": &td, "err": ""})
}

// DataMigrateJobInsert 新增实例
func DataMigrateJobInsert(c *gin.Context) {
	type tmpDataMigrateJob struct {
		model.DataMigrateJob
		SourceDbTables []string `json:"sourceDbTables"`
	}
	// 参数映射到对象
	var tDataMigrateJob tmpDataMigrateJob
	if err := c.ShouldBind(&tDataMigrateJob); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}
	// 获取实例环境类型
	var sourceInstance model.Instance
	sourceInstance.InstId = tDataMigrateJob.SourceInstId
	model.Db.First(&sourceInstance, "inst_id = ?", tDataMigrateJob.SourceInstId)
	var targetInstance model.Instance
	targetInstance.InstId = tDataMigrateJob.TargetInstId
	model.Db.First(&targetInstance, "inst_id = ?", tDataMigrateJob.TargetInstId)
	// 插入主表
	result := model.Db.Create(&tDataMigrateJob.DataMigrateJob)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error})
		return
	}
	// 插入明细表
	var dataMigrateJobDetails []*model.DataMigrateJobDetail
	for _, tableName := range tDataMigrateJob.SourceDbTables {
		var dataMigrateJobDetail model.DataMigrateJobDetail
		dataMigrateJobDetail.DataMigrateJobId = tDataMigrateJob.DataMigrateJob.ID
		dataMigrateJobDetail.TableName = tableName
		dataMigrateJobDetails = append(dataMigrateJobDetails, &dataMigrateJobDetail)
	}
	err := getTablesEstimateRows(sourceInstance, tDataMigrateJob.SourceDb, dataMigrateJobDetails)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error})
		return
	}
	result = model.Db.Create(&dataMigrateJobDetails)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
	}
}

// DataMigrateJobDetailSelectByList 查询详情列表
func DataMigrateJobDetailSelectByList(c *gin.Context) {
	// 变量初始化
	pageNo, _ := strconv.Atoi(c.Query("pageNo"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	jobId := c.Query("jobId")
	var Db = model.Db
	var dataMigrateJobDetails []model.DataMigrateJobDetail
	var totalCount int64
	data := gin.H{"totalCount": 0, "data": &[]model.DataMigrateJobDetail{}, "pageNo": pageNo, "pageSize": pageSize, "totalPage": 0}

	// 执行查询
	Db.Where("data_migrate_job_id = ?", jobId).Order("created_at").Limit(pageSize).Offset((pageNo-1)*pageSize - 1).Find(&dataMigrateJobDetails)
	Db.Where("data_migrate_job_id = ?", jobId).Model(&model.DataMigrateJobDetail{}).Count(&totalCount)

	// 处理结果集并返回
	totalPage := math.Ceil(float64(totalCount) / float64(pageSize))
	if totalCount > 0 {
		data = gin.H{"totalCount": totalCount, "data": &dataMigrateJobDetails, "pageNo": pageNo, "pageSize": pageSize, "totalPage": totalPage}
	}
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": data, "err": ""})
}

// DataMigrateJobExecuteUpdate 执行
func DataMigrateJobExecuteUpdate(c *gin.Context) {
	var dataMigrateJob model.DataMigrateJob
	if err := c.ShouldBind(&dataMigrateJob); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": err.Error()})
		return
	}

	var jobStatus string
	model.Db.Select("status").Model(&dataMigrateJob).Where("id = ?", dataMigrateJob.ID).First(&jobStatus)
	if jobStatus != "NotStart" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": "任务状态不支持执行！"})
		return
	}

	// 获取源实例信息
	var sourceInstance model.Instance
	result := model.Db.First(&sourceInstance, "inst_id = ?", dataMigrateJob.SourceInstId)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": fmt.Sprintf("%v 源实例不存在", dataMigrateJob.SourceInstId)})
		return
	}
	// 获取目标实例信息
	var targetInstance model.Instance
	result = model.Db.First(&targetInstance, "inst_id = ?", dataMigrateJob.TargetInstId)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": fmt.Sprintf("%v 目标实例不存在", dataMigrateJob.TargetInstId)})
		return
	}
	// 更新状态
	result = model.Db.Model(&dataMigrateJob).Where("id = ?", dataMigrateJob.ID).Updates(&model.DataMigrateJob{Status: "Running"})
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "fail", "data": "", "err": result.Error.Error()})
		return
	}

	go func() {
		err := executeJob(dataMigrateJob, sourceInstance, targetInstance)
		if err != nil {
			// 更新状态
			result = model.Db.Model(&dataMigrateJob).Where("id = ?", dataMigrateJob.ID).Updates(&model.DataMigrateJob{Status: "Error"})
			if result.Error != nil {
				log.Errorf("update data migrate job status failed, err: %v", result.Error.Error())
				return
			}
			log.Errorf("data migrate job execute failed, err: %v", err.Error())
			return
		}

		// 更新状态
		result = model.Db.Model(&dataMigrateJob).Where("id = ?", dataMigrateJob.ID).Updates(&model.DataMigrateJob{Status: "Finished"})
		if result.Error != nil {
			log.Errorf("update data migrate job status failed, err: %v", result.Error.Error())
			return
		}
	}()
	c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "success", "data": "", "err": ""})
}

func getTablesEstimateRows(instance model.Instance, db string, dataMigrateJobDetails []*model.DataMigrateJobDetail) error {
	sourceDsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=1s"
	sourceDsn = fmt.Sprintf(sourceDsn, instance.User,
		utils.DecryptAES([]byte(config.Conf.General.SecretKey), instance.Password),
		instance.Ip, instance.Port, db)
	sourceDbConn, err := gorm.Open(mysql.Open(sourceDsn), &gorm.Config{})
	if err != nil {
		return err
	}
	for _, detail := range dataMigrateJobDetails {
		tableRows := 0
		result := sourceDbConn.Raw("select table_rows from information_schema.tables "+
			"where table_schema = ? and table_name = ?", db, detail.TableName).First(&tableRows)
		if result.Error != nil {
			return result.Error
		}
		detail.EstimateRows = uint(tableRows)
	}
	return nil
}

func executeJob(dataMigrateJob model.DataMigrateJob, sourceInstance model.Instance, targetInstance model.Instance) error {
	// create conn
	// get source tables
	// iter tables
	// query table data
	// bulk write target table
	// update detail status
	// 迭代 WorkflowSqlDetail
	tableRows, err := model.Db.Model(&model.DataMigrateJobDetail{}).
		Where("data_migrate_job_id = ?", dataMigrateJob.ID).Order("created_at asc").Rows()
	defer tableRows.Close()
	if err != nil {
		return err
	}

	sourceDsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=1s"
	sourceDsn = fmt.Sprintf(sourceDsn, sourceInstance.User,
		utils.DecryptAES([]byte(config.Conf.General.SecretKey), sourceInstance.Password),
		sourceInstance.Ip, sourceInstance.Port, dataMigrateJob.SourceDb)
	sourceDbConn, err := gorm.Open(mysql.Open(sourceDsn), &gorm.Config{})
	if err != nil {
		return err
	}

	targetDsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=1s"
	targetDsn = fmt.Sprintf(targetDsn, targetInstance.User,
		utils.DecryptAES([]byte(config.Conf.General.SecretKey), targetInstance.Password),
		targetInstance.Ip, targetInstance.Port, dataMigrateJob.TargetDb)
	targetDbConn, err := gorm.Open(mysql.Open(targetDsn), &gorm.Config{})
	if err != nil {
		return err
	}

	var rsError error
	rsError = nil

	// 遍历迁移表
	for tableRows.Next() {
		var dataMigrateJobDetail model.DataMigrateJobDetail
		model.Db.ScanRows(tableRows, &dataMigrateJobDetail)
		columnNames, err := getTablePrimaryKeys(sourceDbConn, dataMigrateJob.SourceDb, dataMigrateJobDetail.TableName)
		if err != nil {
			return err
		}
		if len(columnNames) == 0 {
			return errors.New(fmt.Sprintf("table: %v not primary key", dataMigrateJobDetail.TableName))
		}

		// 更新状态  Running
		model.Db.Model(&dataMigrateJobDetail).Updates(model.DataMigrateJobDetail{
			Status: "Running",
		})

		nowTime := time.Now().Format("20060102150405")
		// rename target table
		renameSql := fmt.Sprintf("rename table %s to %s_backup%s",
			dataMigrateJobDetail.TableName,
			dataMigrateJobDetail.TableName,
			nowTime)
		renameResult := targetDbConn.Exec(renameSql)
		if renameResult.Error != nil {
			// 更新状态  Error
			model.Db.Model(&dataMigrateJobDetail).Updates(model.DataMigrateJobDetail{
				Status:   "Error",
				ErrorMsg: renameResult.Error.Error(),
			})
			if rsError == nil {
				rsError = renameResult.Error // 返回首次异常
			}
			continue
		}
		// get source table ddl
		createDdlMap := map[string]interface{}{}
		getDdlSql := fmt.Sprintf("show create table %s", dataMigrateJobDetail.TableName)
		getDdlResult := sourceDbConn.Raw(getDdlSql).Scan(&createDdlMap)
		if getDdlResult.Error != nil {
			// 更新状态  Error
			model.Db.Model(&dataMigrateJobDetail).Updates(model.DataMigrateJobDetail{
				Status:   "Error",
				ErrorMsg: getDdlResult.Error.Error(),
			})
			if rsError == nil {
				rsError = getDdlResult.Error // 返回首次异常
			}
			continue
		}
		// create target table
		createDdlResult := targetDbConn.Exec(fmt.Sprintf("%v", createDdlMap["Create Table"]))
		if createDdlResult.Error != nil {
			// 更新状态  Error
			model.Db.Model(&dataMigrateJobDetail).Updates(model.DataMigrateJobDetail{
				Status:   "Error",
				ErrorMsg: createDdlResult.Error.Error(),
			})
			if rsError == nil {
				rsError = createDdlResult.Error // 返回首次异常
			}
			continue
		}

		var (
			tx           = sourceDbConn.Order(columnNames[0]).Session(&gorm.Session{})
			queryDB      = tx
			rowsAffected int64
			batch        int
			batchSize    = 1000
			sleep        = 100
		)

		for {
			var results []map[string]interface{}
			queryResult := queryDB.Table(dataMigrateJobDetail.TableName).Debug().Limit(batchSize).Find(&results)
			rowsAffected += queryResult.RowsAffected
			batch++

			for _, rs := range results {
				// 批量处理找到的记录
				fmt.Println(rs)
			}

			if queryResult.Error != nil {
				// 更新状态  Error
				model.Db.Model(&dataMigrateJobDetail).Updates(model.DataMigrateJobDetail{
					Status:   "Error",
					ErrorMsg: queryResult.Error.Error(),
				})
				if rsError == nil {
					rsError = queryResult.Error // 返回首次异常
				}
				break
			}

			if int(queryResult.RowsAffected) == 0 {
				// 更新状态 Finished
				model.Db.Model(&dataMigrateJobDetail).Updates(model.DataMigrateJobDetail{Status: "Finished"})
				break
			}

			writeResult := targetDbConn.Table(dataMigrateJobDetail.TableName).Create(&results)
			if writeResult.Error != nil {
				// 更新状态  Error
				model.Db.Model(&dataMigrateJobDetail).Updates(model.DataMigrateJobDetail{
					Status:   "Error",
					ErrorMsg: writeResult.Error.Error(),
				})
				if rsError == nil {
					rsError = writeResult.Error // 返回首次异常
				}
				break
			}
			affected := writeResult.RowsAffected // 本次批量操作影响的记录数
			// 更新同步rows
			completedRows := dataMigrateJobDetail.CompletedRows + uint(affected)
			model.Db.Model(&dataMigrateJobDetail).Updates(model.DataMigrateJobDetail{
				CompletedRows: completedRows,
			})

			if int(queryResult.RowsAffected) < batchSize {
				// 更新状态 Finished
				model.Db.Model(&dataMigrateJobDetail).Updates(model.DataMigrateJobDetail{Status: "Finished"})
				break
			}

			resultsValue := reflect.Indirect(reflect.ValueOf(results))
			lastResultValue := resultsValue.Index(resultsValue.Len() - 1)
			firstPrimaryKey := columnNames[0]
			firstPrimaryKeyValue := lastResultValue.MapIndex(reflect.ValueOf(firstPrimaryKey))
			pkValue := fmt.Sprintf("%v", firstPrimaryKeyValue.Interface())
			queryDB = tx.Clauses(clause.Gt{Column: clause.Column{Table: clause.CurrentTable, Name: firstPrimaryKey}, Value: pkValue})

			time.Sleep(time.Duration(sleep) * time.Millisecond)
		}
	}
	return rsError
}

func getTablePrimaryKeys(conn *gorm.DB, dbName string, tableName string) ([]string, error) {
	var columnNames []string
	rows, err := conn.Raw("select column_name from information_schema.columns "+
		"where table_schema = ? and table_name = ? and column_key = 'PRI' "+
		"order by ordinal_position", dbName, tableName).Rows()
	if err != nil {
		return columnNames, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName string
		// ScanRows 方法用于将一行记录扫描至结构体
		err := conn.ScanRows(rows, &columnName)
		if err != nil {
			return nil, err
		}
		columnNames = append(columnNames, columnName)
	}
	return columnNames, nil
}

// findInBatches finds all records in batches of batchSize
func findInBatches(db *gorm.DB, dest interface{}, primaryKeys string, batchSize int, fc func(tx *gorm.DB, batch int) error) *gorm.DB {
	var (
		tx           = db.Order(primaryKeys).Session(&gorm.Session{})
		queryDB      = tx
		rowsAffected int64
		batch        int
	)

	// user specified offset or limit
	var totalSize int
	if c, ok := tx.Statement.Clauses["LIMIT"]; ok {
		if limit, ok := c.Expression.(clause.Limit); ok {
			if limit.Limit != nil {
				totalSize = *limit.Limit
			}

			if totalSize > 0 && batchSize > totalSize {
				batchSize = totalSize
			}

			// reset to offset to 0 in next batch
			tx = tx.Offset(-1).Session(&gorm.Session{})
		}
	}

	for {
		var tmpDest []map[string]interface{}
		result := queryDB.Debug().Limit(batchSize).Find(&tmpDest)
		rowsAffected += result.RowsAffected
		batch++

		if result.Error == nil && result.RowsAffected != 0 {
			dest = &tmpDest
			fcTx := result.Session(&gorm.Session{NewDB: true})
			fcTx.RowsAffected = result.RowsAffected
			tx.AddError(fc(fcTx, batch))
		} else if result.Error != nil {
			tx.AddError(result.Error)
		}

		if tx.Error != nil || int(result.RowsAffected) < batchSize {
			break
		}

		if totalSize > 0 {
			if totalSize <= int(rowsAffected) {
				break
			}
			if totalSize/batchSize == batch {
				batchSize = totalSize % batchSize
			}
		}

		resultsValue := reflect.Indirect(reflect.ValueOf(dest))
		lastResultValue := resultsValue.Index(resultsValue.Len() - 1)
		firstPrimaryKey := strings.Split(primaryKeys, ",")[0]
		firstPrimaryKeyValue := lastResultValue.MapIndex(reflect.ValueOf(firstPrimaryKey))
		pkValue := fmt.Sprintf("%v", firstPrimaryKeyValue.Interface())
		queryDB = tx.Clauses(clause.Gt{Column: clause.Column{Table: clause.CurrentTable, Name: firstPrimaryKey}, Value: pkValue})
		// reset dest for gorm
	}

	tx.RowsAffected = rowsAffected
	return tx
}
