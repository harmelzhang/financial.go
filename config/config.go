package config

// ----- 指数 -----

// 指数样本信息获取地址
const indexUrlPrefix = "https://csi-web-dev.oss-cn-shanghai-finance-1-pub.aliyuncs.com"
const indexUrlLocation = "/static/html/csindex/public/uploads/file/autofile/cons/%scons.xls"

// FetchIndexUrlTemplate 查询指数样本信息地址模板
const FetchIndexUrlTemplate = indexUrlPrefix + indexUrlLocation

// ----- 行业分类 -----

// FetchCategoryUrl 查询行业分类地址
const FetchCategoryUrl = "http://quotes.money.163.com/old"

// FetchStockCodesUrl 查询行业下所有的股票代码地址
const FetchStockCodesUrl = "https://quotes.money.163.com/hs/service/diyrank.php?query=PLATE_IDS:%s&fields=SYMBOL&sort=SYMBOL&order=asc&page=%d&count=%d"

// FetchStockCodesCount 每页查询多少条数据
const FetchStockCodesCount = 100

// ----- 股票 -----

// FetchStockBaseInfoNeteaseUrl 查询股票基本信息地址（网易）
const FetchStockBaseInfoNeteaseUrl = "http://quotes.money.163.com/f10/gszl_%s.html"

// FetchStockBaseInfoEastmoneyUrl 查询股票基本信息地址（东方财富）
const FetchStockBaseInfoEastmoneyUrl = "https://emweb.securities.eastmoney.com/PC_HSF10/CompanySurvey/PageAjax?code=%s%s"

// ShanghaiMarketPrefixs 上交所股票前缀
var ShanghaiMarketPrefixs = []string{"60", "68"}

// ShenzhenMarketPrefixs 深交所股票前缀
var ShenzhenMarketPrefixs = []string{"00", "30"}

// BeijingMarketPrefixs 北交所股票前缀
var BeijingMarketPrefixs = []string{"82", "83", "87", "88"}

// ----- 数据库配置 -----

const DB_HOST = "127.0.0.1"  // 服务器地址
const DB_PORT = 3306         // 端口
const DB_USERNAME = "root"   // 用户名
const DB_PASSWORD = "123456" // 密码
const DB_NAME = "financial"  // 数据库名称
