package models

import (
	"encoding/json"
	"financial/config"
	"financial/utils/db"
	"financial/utils/http"
	"fmt"
	"log"
)

// Category 行业分类
type Category struct {
	Id           string // 行业ID（网易财经）
	Name         string // 名称
	DisplayOrder uint8  // 显示顺序
	ParentId     string // 父分类ID
}

// BuildStockCodesUrl 构建根据行业分类查询股票代码的地址
func (category *Category) BuildStockCodesUrl(page int, count int) string {
	return fmt.Sprintf(config.FetchStockCodesUrl, category.Id, page, count)
}

// GetStocks 获取分类下的股票信息
func (category *Category) GetStocks() []*Stock {
	log.Printf("获取分类 [%s] 下的所有股票代码", category.Name)

	stocks := make([]*Stock, 0)
	if category.ParentId == "" {
		return stocks
	}

	page := 0
	for {
		data := make(map[string]interface{})
		err := json.Unmarshal(http.Get(category.BuildStockCodesUrl(page, config.FetchStockCodesCount)), &data)
		if err != nil {
			log.Fatalf("解析JSON数据失败 : %s", err)
		}
		stockCodes := data["list"].([]interface{})
		if len(stockCodes) == 0 {
			break
		}
		for _, stockCode := range stockCodes {
			code := stockCode.(map[string]interface{})
			symbol := code["SYMBOL"].(string)
			stock := &Stock{
				Code:     symbol,
				Category: Category{Id: category.Id, Name: category.Name, ParentId: category.ParentId},
			}
			stocks = append(stocks, stock)
		}
		page++
	}

	return stocks
}

// IntoDb 更新数据库
func (category *Category) IntoDb() {
	sql := "REPLACE INTO category(id, name, display_order, parent_id) VALUES(?, ?, ?, ?)"
	args := []interface{}{category.Id, category.Name, category.DisplayOrder, category.ParentId}
	if category.ParentId == "" {
		args[len(args)-1] = nil
	}
	db.ExecSQL(sql, args...)
}
