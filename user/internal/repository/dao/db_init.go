package dao

import (
	"context"
	"fmt"

	"go-micro-toDoList/user/global"
	"go-micro-toDoList/user/internal/repository/model"
	"time"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var _db *gorm.DB

func InitDB() {
	// 1. dsn
	dsn := dsn()

	// 2. set gorm logger
	var mylogger logger.Interface
	if global.Config.Server.ENV == "debug" {
		mylogger = logger.Default.LogMode(logger.Info)
	} else {
		mylogger = logger.Default
	}

	// 3. gorm open + config
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mylogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true},
	})
	if err != nil {
		global.Logger.Fatalf(fmt.Sprintf("[%s] Mysql connection failed.", dsn))
	}

	// 4. DB()
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 30)

	// 5. migrate
	_db = db
	Migrate()
}

func dsn() string {
	m := global.Config.Mysql
	return m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.Database + "?" + m.Charset
}

func Migrate() {
	err := _db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.User{},
	)
	if err != nil {
		global.Logger.Error("[ error ] Table schema migration failed.")
		return
	}
	global.Logger.Info("[ success ] Table schema migration successful.")
}

func DBWithContext(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
