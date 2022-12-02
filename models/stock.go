package models

// Stock 股票
type Stock struct {
	Code                string   // 股票代码
	StockName           string   // 股票名称
	StockNamePinyin     string   // 股票名称（拼音）
	CompanyName         string   // 公司名称
	Organization        string   // 组织形式（民营、国营...）
	Region              string   // 地域（省份）
	Address             string   // 办公地址
	WebSite             string   // 公司网站
	MainBusiness        string   // 主营业务
	BusinessScope       string   // 业务范围
	DateOfIncorporation string   // 成立日期
	ListingDate         string   // 上市日期
	MainUnderwriter     string   // 主承销商
	Sponsor             string   // 上市保荐人
	AccountingFirm      string   // 会计师事务所
	MarketPlace         string   // 交易市场（上海、深圳、北京）
	Category            Category // 所属行业分类
}
