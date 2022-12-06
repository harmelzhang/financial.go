package db

import (
	"database/sql"
	"financial/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// 获取数据库连接
func getConnection() *sql.DB {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.DB_USERNAME, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatalf("数据库连接出错 : %s", err)
	}
	return db
}

// ExecSQL 执行SQL
func ExecSQL(sql string, args ...any) *sql.Rows {
	db := getConnection()
	defer func() {
		_ = db.Close()
	}()

	rows, err := db.Query(sql, args...)
	if err != nil {
		log.Fatalf("SQL执行出错 : %s", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	return rows
}
