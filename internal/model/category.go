package model

// 行业分类
type Category struct {
	// 分类类型（证监会、中证）
	Type string `json:"type"`
	// 行业代码
	Code string `json:"code"`
	// 名称
	Name string `json:"name"`
	// 层级
	Level string `json:"level"`
	// 显示顺序
	DisplayOrder int `json:"display_order"`
	// 父行业代码
	ParentCode any `json:"parent_code"`
}

// 行业分类表所有列信息
type categoryColumns struct {
	// 分类类型（证监会、中证）
	Type string
	// 行业代码
	Code string
	// 名称
	Name string
	// 层级
	Level string
	// 显示顺序
	DisplayOrder string
	// 父行业代码
	ParentCode string
}

// 行业分类表信息
type categoryTableInfo struct {
	// 表名
	table string
	// 所有列名
	columns categoryColumns
}

var CategoryTableInfo = categoryTableInfo{
	table: "category",
	columns: categoryColumns{
		Type:         "type",
		Code:         "code",
		Name:         "name",
		Level:        "level",
		DisplayOrder: "display_order",
		ParentCode:   "parent_code",
	},
}

// 数据表名
func (info *categoryTableInfo) Table() string {
	return info.table
}

// 字段名（列名）
func (info *categoryTableInfo) Columns() categoryColumns {
	return info.columns
}
