package index

import (
	"financial/config"
	"financial/utils/http"
	"financial/utils/tools"
	"financial/utils/xls"
	"fmt"
	"log"
)

// IndexType 指数类型
type IndexType string

const (
	HS300 IndexType = "000300" // 沪深300
	ZZ500 IndexType = "000905" // 中证500
	SZ50  IndexType = "000016" // 上证50
	KC50  IndexType = "000688" // 科创50
	HLZS  IndexType = "000015" // 红利指数
)

var IndexTypeMap = map[IndexType]string{
	HS300: "沪深300",
	ZZ500: "中证500",
	SZ50:  "上证50",
	KC50:  "科创50",
	HLZS:  "红利指数",
}

// 指数样本信息
var indexUrlMap = map[IndexType]string{}

// 指数样本信息
var indexStockMap = make(map[IndexType][]string)

func init() {
	// 初始化指数样本信息地址
	for indexType, _ := range IndexTypeMap {
		indexUrlMap[indexType] = fmt.Sprintf(config.FetchIndexUrlTemplate, indexType)
	}

	log.Println("初始化主要指数样本信息")
	for indexType, url := range indexUrlMap {
		data := xls.ReadXls(http.Get(url), 0, 0)
		stockCodes := tools.FetchColData(data, 4)
		indexStockMap[indexType] = append(indexStockMap[indexType], stockCodes...)
	}
}

// GetStockIndexTypes 获取指定股票指数类型
func GetStockIndexTypes(stockCode string) ([]IndexType, []string) {
	indexTypes := make([]IndexType, 0, len(IndexTypeMap))
	indexTypeNames := make([]string, 0, len(IndexTypeMap))
	for indexType, stockCodes := range indexStockMap {
		if tools.IndexOf(stockCodes, stockCode) != -1 {
			indexTypes = append(indexTypes, indexType)
			indexTypeNames = append(indexTypeNames, IndexTypeMap[indexType])
		}
	}
	return indexTypes, indexTypeNames
}

// GetIndexStocks 根据指数类型查询样本股票代码
func GetIndexStocks(indexType IndexType) []string {
	return indexStockMap[indexType]
}
