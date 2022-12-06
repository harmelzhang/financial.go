package main

import (
	categoryConfig "financial/config/category"
	_ "financial/config/indexSample"
)

func main() {
	for _, category := range categoryConfig.GetCategorys() {
		category.IntoDb()
		stocks := category.GetStocks()
		for _, stock := range stocks {
			stock.BuildStockInfo()
			stock.IntoDb()
		}
	}
}
