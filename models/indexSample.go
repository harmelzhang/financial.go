package models

import "financial/utils/db"

// IndexSample 指数样本信息
type IndexSample struct {
	TypeCode  string // 类型代码（中证指数，www.csindex.com.cn）
	TypeName  string // 类型名称（沪深300、中证500、上证50....）
	StockCode string // 股票代码
}

// IntoDb 更新数据库
func (indexSample *IndexSample) IntoDb() {
	sql := "INSERT INTO index_sample(type_code, type_name, stock_code) VALUES(?, ?, ?)"
	args := []interface{}{indexSample.TypeCode, indexSample.TypeName, indexSample.StockCode}
	db.ExecSQL(sql, args...)
}
