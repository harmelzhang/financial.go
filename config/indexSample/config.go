package indexSample

import (
	"financial/config"
	"financial/models"
	"financial/utils/db"
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
	for indexType, _ := range TypeMap {
		indexUrlMap[indexType] = fmt.Sprintf(config.FetchIndexUrlTemplate, indexType)
	}

	log.Println("初始化主要指数样本信息")
	for indexType, url := range indexUrlMap {
		data := xls.ReadXls(http.Get(url), 0, 0)
		stockCodes := tools.FetchColData(data, 4)
		indexStockMap[indexType] = append(indexStockMap[indexType], stockCodes...)
	}

	// 清空数据
	db.ExecSQL("DELETE FROM index_sample")
	// 插入数据
	for indexType, stockCodes := range indexStockMap {
		for _, stockCode := range stockCodes {
			indexSample := models.IndexSample{
				TypeCode:  string(indexType),
				TypeName:  TypeMap[indexType],
				StockCode: stockCode,
			}
			indexSample.IntoDb()
		}
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
