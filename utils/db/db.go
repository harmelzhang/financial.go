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
