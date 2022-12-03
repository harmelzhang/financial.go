package index

import (
	"financial/config"
	"financial/utils/http"
	"financial/utils/tools"
	"financial/utils/xls"
	"fmt"
	"log"
)

// Type 指数类型
type Type string

const (
	HS300 Type = "000300" // 沪深300
	ZZ500 Type = "000905" // 中证500
	SZ50  Type = "000016" // 上证50
	KC50  Type = "000688" // 科创50
	HLZS  Type = "000015" // 红利指数
)

var TypeMap = map[Type]string{
	HS300: "沪深300",
	ZZ500: "中证500",
	SZ50:  "上证50",
	KC50:  "科创50",
	HLZS:  "红利指数",
}

// 指数样本信息
var indexUrlMap = map[Type]string{}

// 指数样本信息
var indexStockMap = make(map[Type][]string)

func init() {
	// 初始化指数样本信息地址
	for Type, _ := range TypeMap {
		indexUrlMap[Type] = fmt.Sprintf(config.FetchIndexUrlTemplate, Type)
	}

	log.Println("初始化主要指数样本信息")
	for Type, url := range indexUrlMap {
		data := xls.ReadXls(http.Get(url), 0, 0)
		stockCodes := tools.FetchColData(data, 4)
		indexStockMap[Type] = append(indexStockMap[Type], stockCodes...)
	}
}

// GetStockTypes 获取指定股票指数类型
func GetStockTypes(stockCode string) ([]Type, []string) {
	Types := make([]Type, 0, len(TypeMap))
	TypeNames := make([]string, 0, len(TypeMap))
	for Type, stockCodes := range indexStockMap {
		if tools.IndexOf(stockCodes, stockCode) != -1 {
			Types = append(Types, Type)
			TypeNames = append(TypeNames, TypeMap[Type])
		}
	}
	return Types, TypeNames
}

// GetIndexStocks 根据指数类型查询样本股票代码
func GetIndexStocks(Type Type) []string {
	return indexStockMap[Type]
}
