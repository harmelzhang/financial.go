package model

// 行业分类与股票代码之间的关系
type CategoryStockCode struct {
	// 行业代码
	CategoryCode string
	// 股票代码
	StockCode string
}

// 行业分类与股票代码关系表所有列信息
type categoryStockCodeColumns struct {
	// 行业代码
	CategoryCode string
	// 股票代码
	StockCode string
}

// 行业分类与股票代码关系表信息
type categoryStockCodeTableInfo struct {
	// 表名
	table string
	// 所有列名
	columns categoryStockCodeColumns
}

var CategoryStockCodeTableInfo = categoryStockCodeTableInfo{
	table: "category_stock_code",
	columns: categoryStockCodeColumns{
		CategoryCode: "category_code",
		StockCode:    "stock_code",
	},
}

// 数据表名
func (info *categoryStockCodeTableInfo) Table() string {
	return info.table
}

// 字段名（列名）
func (info *categoryStockCodeTableInfo) Columns() categoryStockCodeColumns {
	return info.columns
}
