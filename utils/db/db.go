package db

import (
	"database/sql"
	"financial/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var db *sql.DB
var err error

// 初始化数据库连接
func init() {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.DbUsername, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	db, err = sql.Open("mysql", dataSource)
	db.SetMaxIdleConns(config.DbMaxIdleConns)
	db.SetConnMaxIdleTime(config.DbMaxIdleTime * time.Minute)
	db.SetConnMaxLifetime(config.DbMaxLifeTime * time.Minute)
	if err != nil {
		log.Fatalf("数据库连接出错 : %s", err)
	}
}

// GetDb 获取数据库连接
func GetDb() *sql.DB {
	return db
}

// ExecSQL 执行SQL
func ExecSQL(sql string, args ...any) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		log.Fatalf("SQL执行出错 : %s", err)
	}
	defer func() {
		_ = rows.Close()
	}()
}
