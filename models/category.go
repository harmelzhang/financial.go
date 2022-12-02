package models

import (
	"financial/config"
	"financial/utils/http"
	"fmt"
)

// Category 行业分类
type Category struct {
	Id       string // 行业ID（网易财经）
	Name     string // 名称
	ParentId string // 父分类ID
}

func (category *Category) BuildStockCodesUrl(page int, count int) string {
	return fmt.Sprintf(config.FetchStockCodesUrl, category.Id, page, count)
}

// GetStocks 获取分类下的股票信息
func (category *Category) GetStocks() []Stock {
	stocks := make([]Stock, 0)

	if category.ParentId == "" {
		return stocks
	}

	page := 0
	for {
		data := http.Get(category.BuildStockCodesUrl(page, config.FetchStockCodesCount))
		fmt.Println(string(data))
		break
	}

	return stocks
}
