package config

import (
	"financial/utils/http"
	"financial/utils/tools"
	"financial/utils/xls"
	"log"
)

// IndexType 指数类型
type IndexType uint8

const (
	HS300 IndexType = iota
	SZ50
	ZZ500
)

// 指数样本信息获取地址
var indexUrlMap = map[IndexType]string{
	HS300: "https://csi-web-dev.oss-cn-shanghai-finance-1-pub.aliyuncs.com/static/html/csindex/public/uploads/file/autofile/cons/000300cons.xls",
	SZ50:  "https://csi-web-dev.oss-cn-shanghai-finance-1-pub.aliyuncs.com/static/html/csindex/public/uploads/file/autofile/cons/000016cons.xls",
	ZZ500: "https://csi-web-dev.oss-cn-shanghai-finance-1-pub.aliyuncs.com/static/html/csindex/public/uploads/file/autofile/cons/000905cons.xls",
}

// 指数样本信息
var indexTypeMap = make(map[IndexType][]string)

func init() {
	log.Println("初始化，读取主要指数样本信息")
	for indexType, url := range indexUrlMap {
		data := xls.ReadXls(http.Get(url), 0, 0)
		stockCodes := tools.FetchColData(data, 4)
		indexTypeMap[indexType] = append(indexTypeMap[indexType], stockCodes...)
	}
}

// GetStockIndexTypes 获取指定股票指数类型
func GetStockIndexTypes(stockCode string) []IndexType {
	indexTypes := make([]IndexType, 0, 3)
	for indexType, stockCodes := range indexTypeMap {
		if tools.IndexOf(stockCodes, stockCode) != -1 {
			indexTypes = append(indexTypes, indexType)
		}
	}
	return indexTypes
}
