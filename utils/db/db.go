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
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)
	db, err = sql.Open("mysql", dataSource)
	db.SetMaxIdleConns(config.DB_MAX_IDLE_CONNS)
	db.SetConnMaxIdleTime(config.DB_MAX_IDLE_TIME * time.Minute)
	db.SetConnMaxLifetime(config.DB_MAX_LIFE_TIME * time.Minute)
	if err != nil {
		log.Fatalf("数据库连接出错 : %s", err)
	}
}

// ExecSQL 执行SQL
func ExecSQL(sql string, args ...any) *sql.Rows {
	rows, err := db.Query(sql, args...)
	if err != nil {
		log.Fatalf("SQL执行出错 : %s", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	return rows
}
