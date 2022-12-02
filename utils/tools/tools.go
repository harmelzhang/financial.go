package tools

import (
	"golang.org/x/net/html"
	"log"
)

// IndexOf 查询指定元素在数组中的位置
func IndexOf[T string](source []T, target T) int {
	for index, item := range source {
		if item == target {
			return index
		}
	}
	return -1
}

// FetchColData 获取指定二维数组索引列的数据
func FetchColData(table [][]string, colIndex int) []string {
	data := make([]string, 0)
	for _, row := range table {
		for i := 0; i < len(row); i++ {
			if i == colIndex {
				data = append(data, row[i])
			}
		}
	}
	return data
}

// FindAttrVal 查询指定属性的值
func FindAttrVal(node *html.Node, attrName string) string {
	value, hasAttr := "", false
	for _, attr := range node.Attr {
		if attr.Key == attrName {
			hasAttr = true
			value = attr.Val
		}
	}
	if !hasAttr {
		log.Fatalln("找不到指定属性")
	}
	return value
}
