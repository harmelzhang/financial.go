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

// 根据股票代码查询数据库中所有财报报告日期
func (s *financialService) GetReportDates(ctx context.Context, stockCode string) (reportDates []string, err error) {
	return dao.FinancialDao.GetReportDates(ctx, stockCode)
}

// 根据类型查询所有报告
func (s *financialService) GetByType(ctx context.Context, stockCode, reportType string) (financials []*model.Financial, err error) {
	return dao.FinancialDao.GetByType(ctx, stockCode, reportType)
}
