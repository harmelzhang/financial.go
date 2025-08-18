package service

import (
	"context"

	"harmel.cn/financial/internal/dao"
	"harmel.cn/financial/internal/model"
)

type categoryStockCodeService struct{}

// 行业分类与股票代码业务逻辑处理对象
var CategoryStockCodeService = new(categoryStockCodeService)

// 插入记录
func (cscService *categoryStockCodeService) Insert(ctx context.Context, entity *model.CategoryStockCode) (err error) {
	err = dao.CategoryStockCodeDao.Insert(ctx, entity)
	return
}

// 删除所有数据
func (cscService *categoryStockCodeService) DeleteAll(ctx context.Context) (err error) {
	err = dao.CategoryDao.DeleteAll(ctx)
	return
}
