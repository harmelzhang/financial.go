package model

// 财务报表
type Financial struct {
	// 股票代码
	StockCode string `json:"stock_code"`
	// 年份
	Year string `json:"year"`
	// 财报季期
	ReportDate string `json:"report_date"`
	// 季期类型（Q1~Q4，分别代表：一季报、半年报、三季报、年报；O，代表：其他）
	ReportType string `json:"report_type"`

	// 年度分红金额
	Dividend any `json:"dividend"`

	// 营业活动现金流量
	Ocf any `json:"ocf"`
	// 投资活动现金流量
	Cfi any `json:"cfi"`
	// 筹资活动现金流量
	Cff any `json:"cff"`
	// 分配股利、利润或偿付利息支付的现金
	AssignDividendPorfit any `json:"assign_dividend_porfit"`
	// 购建固定资产、无形资产和其他长期资产支付的现金
	AcquisitionAssets any `json:"acquisition_assets"`
	// 存货减少额
	InventoryLiquidating any `json:"inventory_liquidating"`

	// 净利润
	Np any `json:"np"`
	// 营业收入
	Oi any `json:"oi"`
	// 营业成本
	Coe any `json:"coe"`
	// 营业总成本（含各种费用，销售费用、管理费用等）
	CoeTotal any `json:"coe_total"`
	// 每股盈余|基本每股收益
	Eps any `json:"eps"`

	// 货币资金
	MonetaryFund any `json:"monetary_fund"`
	// 交易性金融资产
	TradeFinassetNotfvtpl any `json:"trade_finasset_notfvtpl"`
	// 交易性金融资产（历史遗留）
	TradeFinasset any `json:"trade_finasset"`
	// 衍生金融资产
	DeriveFinasset any `json:"derive_finasset"`

	// 固定资产
	FixedAsset any `json:"fixed_asset"`
	// 在建工程
	Cip any `json:"cip"`

	// 流动资产总额
	CaTotal any `json:"ca_total"`
	// 非流动资产总额
	NcaTotal any `json:"nca_total"`
	// 流动负债总额
	ClTotal any `json:"cl_total"`
	// 非流动负债产总额
	NclTotal any `json:"ncl_total"`
	// 存货
	Inventory any `json:"inventory"`
	// 应收账款
	AccountsRece any `json:"accounts_rece"`
	// 应付账款
	AccountsPayable any `json:"accounts_payable"`

	// 现金流量允当比例
	CashFlowAdequacyRatio any `json:"cash_flow_adequacy_ratio"`
}

// 财务报表表所有列信息
type financialColumns struct {
	// 股票代码
	StockCode string
	// 年份
	Year string
	// 财报季期
	ReportDate string
	// 季期类型（Q1~Q4，分别代表：一季报、半年报、三季报、年报；O，代表：其他）
	ReportType string

	// 年度分红金额
	Dividend string

	// 营业活动现金流量
	Ocf string
	// 投资活动现金流量
	Cfi string
	// 筹资活动现金流量
	Cff string
	// 分配股利、利润或偿付利息支付的现金
	AssignDividendPorfit string
	// 购建固定资产、无形资产和其他长期资产支付的现金
	AcquisitionAssets string
	// 存货减少额
	InventoryLiquidating string

	// 净利润
	Np string
	// 营业收入
	Oi string
	// 营业成本
	Coe string
	// 营业总成本（含各种费用，销售费用、管理费用等）
	CoeTotal string
	// 每股盈余|基本每股收益
	Eps string

	// 货币资金
	MonetaryFund string
	// 交易性金融资产
	TradeFinassetNotfvtpl string
	// 交易性金融资产（历史遗留）
	TradeFinasset string
	// 衍生金融资产
	DeriveFinasset string

	// 固定资产
	FixedAsset string
	// 在建工程
	Cip string

	// 流动资产总额
	CaTotal string
	// 非流动资产总额
	NcaTotal string
	// 流动负债总额
	ClTotal string
	// 非流动负债产总额
	NclTotal string
	// 存货
	Inventory string
	// 应收账款
	AccountsRece string
	// 应付账款
	AccountsPayable string

	// 现金流量允当比例
	CashFlowAdequacyRatio string
}

// 财务报表表信息
type financialTableInfo struct {
	// 表名
	table string
	// 所有列名
	columns financialColumns
}

var FinancialTableInfo = financialTableInfo{
	table: "financial",
	columns: financialColumns{
		StockCode:             "stock_code",
		Year:                  "year",
		ReportDate:            "report_date",
		ReportType:            "report_type",
		Dividend:              "dividend",
		Ocf:                   "ocf",
		Cfi:                   "cfi",
		Cff:                   "cff",
		AssignDividendPorfit:  "assign_dividend_porfit",
		AcquisitionAssets:     "acquisition_assets",
		InventoryLiquidating:  "inventory_liquidating",
		Np:                    "np",
		Oi:                    "oi",
		Coe:                   "coe",
		CoeTotal:              "coe_total",
		Eps:                   "eps",
		MonetaryFund:          "monetary_fund",
		TradeFinassetNotfvtpl: "trade_finasset_notfvtpl",
		TradeFinasset:         "trade_finasset",
		DeriveFinasset:        "derive_finasset",
		FixedAsset:            "fixed_asset",
		Cip:                   "cip",
		CaTotal:               "ca_total",
		NcaTotal:              "nca_total",
		ClTotal:               "cl_total",
		NclTotal:              "ncl_total",
		Inventory:             "inventory",
		AccountsRece:          "accounts_rece",
		AccountsPayable:       "accounts_payable",
		CashFlowAdequacyRatio: "cash_flow_adequacy_ratio",
	},
}

// 数据表名
func (info *financialTableInfo) Table() string {
	return info.table
}

// 字段名（列名）
func (info *financialTableInfo) Columns() financialColumns {
	return info.columns
}
