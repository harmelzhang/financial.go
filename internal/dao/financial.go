package dao

import (
	"context"

	"harmel.cn/financial/internal/model"
)

type financialDao struct{}

// 财务报表数据访问层
var FinancialDao = new(financialDao)

// 是否存在
func (dao *financialDao) IsExist(ctx context.Context, code, reportDate string) (exist bool, err error) {
	cnt, err := DB(ctx, model.FinancialTableInfo.Table()).
		Where(model.FinancialTableInfo.Columns().StockCode, code).
		Where(model.FinancialTableInfo.Columns().ReportDate, reportDate).
		Count()
	if err != nil {
		return
	}
	if cnt > 0 {
		exist = true
	}
	return
}

// 插入数据
func (dao *financialDao) Insert(ctx context.Context, entity *model.Financial) (err error) {
	_, err = DB(ctx, model.FinancialTableInfo.Table()).Data(entity).Insert()
	return
}

// 更新数据
func (dao *financialDao) Update(ctx context.Context, entity *model.Financial) (err error) {
	_, err = DB(ctx, model.FinancialTableInfo.Table()).
		Data(entity).
		Where(model.FinancialTableInfo.Columns().StockCode, entity.StockCode).
		Where(model.FinancialTableInfo.Columns().ReportDate, entity.ReportDate).
		Update()
	return
}

// 根据股票代码查询数据库中所有财报报告日期
func (dao *financialDao) GetReportDates(ctx context.Context, stockCode string) (reportDates []string, err error) {
	items, err := DB(ctx, model.FinancialTableInfo.Table()).
		Fields(model.FinancialTableInfo.Columns().ReportDate).
		Where(model.FinancialTableInfo.Columns().StockCode, stockCode).
		Array()
	if err != nil {
		return
	}
	for _, item := range items {
		reportDates = append(reportDates, item.String())
	}
	return
}

// 根据类型查询所有报告
func (dao *financialDao) GetByType(ctx context.Context, stockCode, reportType string) (financials []*model.Financial, err error) {
	err = DB(ctx, model.FinancialTableInfo.Table()).
		Where(model.FinancialTableInfo.Columns().StockCode, stockCode).
		Where(model.FinancialTableInfo.Columns().ReportType, reportType).
		Scan(&financials)
	return
}
