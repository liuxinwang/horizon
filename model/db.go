package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"horizon/config"
	"log"
	"os"
	"time"
)

var Db *gorm.DB

func InitDb() {
	dbConf := config.Conf.Mysql
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=3s"
	dsn = fmt.Sprintf(dsn, dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Db)
	var file *os.File
	var err error
	if config.Conf.General.Environment == "dev" {
		file = os.Stdout
	} else {
		file, err = os.OpenFile(config.Conf.Log.Name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	}
	if err != nil {
		logrus.Fatalf("db init failed, err: %v", err.Error())
	}
	newLogger := logger.New(
		log.New(file, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,         // Don't include params in the SQL log
			Colorful:                  false,        // Disable color
		},
	)
	Db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: false, Logger: newLogger})
	err = Db.AutoMigrate(
		&User{},
		&Instance{},
		&Inspection{},
		&Metric{},
		&InstMetric{},
		&InspDetail{},
		&Score{},
		&Task{},
		&Role{},
		&Menu{},
		&UserRole{},
		&RolePermission{},
		&Project{},
		&ProjectRole{},
		&ProjectDatasource{},
		&ProjectUser{},
		&RuleTemplate{},
		&Workflow{},
		&WorkflowTemplate{},
		&WorkflowTemplateDetail{},
		&WorkflowRecord{},
		&WorkflowSqlDetail{},
		&DataMigrateJob{},
		&DataMigrateJobDetail{})
	if err != nil {
		log.Fatal(err.Error())
	}
}
