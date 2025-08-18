package dao

import (
	"context"

	"harmel.cn/financial/internal/model"
)

type categoryStockCodeDao struct{}

// 行业分类与股票代码数据访问层
var CategoryStockCodeDao = new(categoryStockCodeDao)

// 插入记录
func (dao *categoryStockCodeDao) Insert(ctx context.Context, entity *model.CategoryStockCode) (err error) {
	_, err = DB(ctx, model.CategoryStockCodeTableInfo.Table()).Insert(entity)
	return
}

// 删除所有数据
func (dao *categoryStockCodeDao) DeleteAll(ctx context.Context) (err error) {
	_, err = DB(ctx, model.CategoryStockCodeTableInfo.Table()).Delete()
	return
}
