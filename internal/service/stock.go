package service

import (
	"context"

	"harmel.cn/financial/internal/dao"
	"harmel.cn/financial/internal/model"
)

type stockService struct{}

// 股票逻辑处理对象
var StockService = new(stockService)

// 更新或插入
func (s *stockService) Replace(ctx context.Context, entity *model.Stock) (err error) {
	exist, err := dao.StockDao.IsExist(ctx, entity.Code)
	if err != nil {
		return
	}
	if exist {
		err = dao.StockDao.Update(ctx, entity)
	} else {
		err = dao.StockDao.Insert(ctx, entity)
	}
	return
}

// 根据股票代码查询股票信息
func (s *stockService) FindStockByCode(ctx context.Context, stockCode string) (stock *model.Stock, err error) {
	return dao.StockDao.FindStockByCode(ctx, stockCode)
}
