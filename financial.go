package main

import (
	categoryConfig "financial/config/category"
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
