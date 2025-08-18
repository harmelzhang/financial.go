package model

// 指数样本（来源：中证指数）
type IndexSample struct {
	// 类型代码（中证指数，www.csindex.com.cn）
	TypeCode string
	// 股票代码
	StockCode string
}

// 指数样本表所有列信息
type indexSampleColumns struct {
	// 类型代码（中证指数，www.csindex.com.cn）
	TypeCode string
	// 股票代码
	StockCode string
}

// 指数样本表信息
type indexSampleTableInfo struct {
	// 表名
	table string
	// 所有列名
	columns indexSampleColumns
}

var IndexSampleTableInfo = indexSampleTableInfo{
	table: "index_sample",
	columns: indexSampleColumns{
		TypeCode:  "type_code",
		StockCode: "stock_code",
	},
}

// 数据表名
func (info *indexSampleTableInfo) Table() string {
	return info.table
}

// 字段名（列名）
func (info *indexSampleTableInfo) Columns() indexSampleColumns {
	return info.columns
}
