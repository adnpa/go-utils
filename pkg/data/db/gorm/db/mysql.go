package db

import (
	"adnpa/go-utils/pkg/data/db/gorm/model"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//https://gorm.io/docs/

var db *gorm.DB

func init() {
	var err error

	//todo 替换config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "123456", "127.0.0.1", 3306, "mydb")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AllowGlobalUpdate = true
	//todo
	//sqlDb, _ := db.DB()
	//sqlDb.SetMaxOpenConns(conf.Cfg.MaxOpenConns)
	//sqlDb.SetMaxIdleConns(conf.Cfg.MaxIdleConns)
	migration()
	return
}

func migration() {
	db.Set(`gorm:table_options`, "charset=utf8mb4").AutoMigrate(&model.User{})
}

func NewDBClient(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}
