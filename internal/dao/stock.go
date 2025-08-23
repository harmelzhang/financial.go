package dao

import (
	"context"

	"harmel.cn/financial/internal/model"
)

type stockDao struct{}

// 股票数据访问层
var StockDao = new(stockDao)

// 是否存在
func (dao *stockDao) IsExist(ctx context.Context, code string) (exist bool, err error) {
	cnt, err := DB(ctx, model.StockTableInfo.Table()).
		Where(model.StockTableInfo.Columns().Code, code).
		Count()
	if err != nil {
		return
	}
	if cnt > 0 {
		exist = true
	}
	return
}

// 根据代码查询
func (dao *stockDao) QueryByCode(ctx context.Context, code string) (entity *model.Stock, err error) {
	err = DB(ctx, model.StockTableInfo.Table()).
		Where(model.StockTableInfo.Columns().Code, code).
		Scan(&entity)
	return
}

// 插入数据
func (dao *stockDao) Insert(ctx context.Context, entity *model.Stock) (err error) {
	_, err = DB(ctx, model.StockTableInfo.Table()).Data(entity).Insert()
	return
}

// 更新数据
func (dao *stockDao) Update(ctx context.Context, entity *model.Stock) (err error) {
	_, err = DB(ctx, model.StockTableInfo.Table()).
		Data(entity).
		Where(model.StockTableInfo.Columns().Code, entity.Code).
		Update()
	return
}

// 根据代码查询股票信息
func (dao *stockDao) FindStockByCode(ctx context.Context, code string) (entity *model.Stock, err error) {
	err = DB(ctx, model.StockTableInfo.Table()).
		Where(model.StockTableInfo.Columns().Code, code).
		Scan(&entity)
	return
}
