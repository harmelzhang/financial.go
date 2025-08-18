package service

import (
	"context"

	"harmel.cn/financial/internal/dao"
	"harmel.cn/financial/internal/model"
)

type financialService struct{}

// 股票逻辑处理对象
var FinancialService = new(financialService)

// 更新或插入
func (s *financialService) Replace(ctx context.Context, entity *model.Financial) (err error) {
	exist, err := dao.FinancialDao.IsExist(ctx, entity.StockCode, entity.ReportDate)
	if err != nil {
		return
	}
	if exist {
		err = dao.FinancialDao.Update(ctx, entity)
	} else {
		err = dao.FinancialDao.Insert(ctx, entity)
	}
	return
}
