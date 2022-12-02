package main

import (
	"financial/config/category"
)

func main() {
	for _, category := range category.GetCategorys() {
		category.GetStocks()
	}
}
