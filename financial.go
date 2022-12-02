package main

import (
	"financial/config"
	"fmt"
)

func main() {
	for _, category := range config.GetCategorys() {
		fmt.Println(category.Id, category.Name)
	}
}
