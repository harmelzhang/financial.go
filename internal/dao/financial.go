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
