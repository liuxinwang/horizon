package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"horizon/config"
	"log"
)

var Db *gorm.DB

func InitDb() {
	dbConf := config.Conf.Mysql
	dsn := "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=3s"
	dsn = fmt.Sprintf(dsn, dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Db)
	Db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	err := Db.AutoMigrate(
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
		&Workflow{})
	if err != nil {
		log.Fatal(err.Error())
	}
}
