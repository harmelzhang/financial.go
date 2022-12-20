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
var BeijingMarketPrefixs = []string{"82", "83", "87", "88", "43"}

// ExcludeStockCodePrefix 排除的股票（B股、场内基金）
var ExcludeStockCodePrefix = []string{"1", "2", "5", "9"}

// ----- 数据库配置 -----

const DbHost = "127.0.0.1"  // 服务器地址
const DbPort = 3306         // 端口
const DbUsername = "root"   // 用户名
const DbPassword = "123456" // 密码
const DbName = "financial"  // 数据库名称
const DbMaxIdleConns = 100  // 最大空闲连接数
const DbMaxIdleTime = 2     // 连接最大空闲时长（单位：分）
const DbMaxLifeTime = 1     // 连接最大存活时长（单位：分）

// ----- 爬取进度 -----

const ProgressFileName = "progress.json" // 配置文件路径
const TaskIntervalDay = 7                // 任务周期天数

// ----- HTTP -----

const HttpAccept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"

var UserAgent = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.17763",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36 Edg/90.0.818.46",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.3 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:65.0) Gecko/20100101 Firefox/65.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 7_0_4 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) CriOS/31.0.1650.18 Mobile/11B554a Safari/8536.25",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 8_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12F70 Safari/600.1.4",
	"Mozilla/5.0 (Linux; Android 4.2.1; M040 Build/JOP40D) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.59 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; U; Android 4.4.4; zh-cn; M351 Build/KTU84P) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
}
