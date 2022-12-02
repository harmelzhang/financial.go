package models

// Category 行业分类
type Category struct {
	Id       string // 行业ID（网易财经）
	Name     string // 名称
	ParentId string // 父分类ID
}
