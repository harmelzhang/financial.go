package main

import (
	"financial/config"
	"fmt"
)

func main() {
	indexs := config.GetStockIndexTypes("600519")
	for _, index := range indexs {
		switch index {
		case config.HS300:
			fmt.Println("沪深300")
		case config.SZ50:
			fmt.Println("上证50")
		case config.ZZ500:
			fmt.Println("中证500")
		default:
			fmt.Println("查询不到所属指数信息")
		}
	}
}
