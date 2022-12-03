package config

// ----- 指数 -----

// 指数样本信息获取地址
var indexUrlPrefix = "https://csi-web-dev.oss-cn-shanghai-finance-1-pub.aliyuncs.com"
var indexUrlLocation = "/static/html/csindex/public/uploads/file/autofile/cons/%scons.xls"

// FetchIndexUrlTemplate 查询指数样本信息地址模板
var FetchIndexUrlTemplate = indexUrlPrefix + indexUrlLocation

// ----- 行业分类 -----

// FetchCategoryUrl 查询行业分类地址
var FetchCategoryUrl = "http://quotes.money.163.com/old"

// FetchStockCodesUrl 查询行业下所有的股票代码地址
var FetchStockCodesUrl = "https://quotes.money.163.com/hs/service/diyrank.php?query=PLATE_IDS:%s&fields=SYMBOL&sort=SYMBOL&order=asc&page=%d&count=%d"

// FetchStockCodesCount 每页查询多少条数据
var FetchStockCodesCount = 100

// ----- 股票 -----

// FetchStockBaseInfoNeteaseUrl 查询股票基本信息地址（网易）
var FetchStockBaseInfoNeteaseUrl = "http://quotes.money.163.com/f10/gszl_%s.html"

// FetchStockBaseInfoEastmoneyUrl 查询股票基本信息地址（东方财富）
var FetchStockBaseInfoEastmoneyUrl = "https://emweb.securities.eastmoney.com/PC_HSF10/CompanySurvey/PageAjax?code=%s%s"

var ShanghaiMarketPrefixs = []string{"60", "68"}
var ShenzhenMarketPrefixs = []string{"00", "30"}
var BeijingMarketPrefixs = []string{"82", "83", "87", "88"}
