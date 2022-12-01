package models

// Category 行业分类
type Category struct {
	Id       string // ID
	Name     string // 名称
	ParentId string // 父分类ID
}
