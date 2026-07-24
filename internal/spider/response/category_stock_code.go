package response

// 股票代码
type StockCode struct {
	// 代码
	Code string `json:"securityCode"`
	// 中证1级代码
	CicsLeve1Code string `json:"cics1stCode"`
	// 中证2级代码
	CicsLeve2Code string `json:"cics2ndCode"`
	// 中证3级代码
	CicsLeve3Code string `json:"cics3rdCode"`
	// 中证4级代码
	CicsLeve4Code string `json:"cics4thCode"`
	// 证券会1级代码
	CsrcLeve1Code string `json:"csrc1stCode"`
	// 证券会2级代码
	CsrcLeve2Code string `json:"csrc2ndCode"`
}

// 股票代码数据
type StockCodeData struct {
	// 交易日期
	Date string `json:"tradeDate"`
	// 股票代码
	List []StockCode `json:"list"`
}

// 股票代码响应
type StockCodeResult struct {
	// 状态码
	Code string `json:"code"`
	// 是否成功
	Success bool `json:"success"`
	// 消息
	Message string `json:"msg"`
	// 股票代码数据
	Data StockCodeData `json:"data"`
}
