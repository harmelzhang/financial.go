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

// FetchStockCodesUrl 查询行业下所有的股票代码
var FetchStockCodesUrl = "https://quotes.money.163.com/hs/service/diyrank.php?query=PLATE_IDS:%s&fields=SYMBOL&sort=SYMBOL&order=asc&page=%d&count=%d"

// FetchStockCodesCount 每页查询多少条数据
var FetchStockCodesCount = 100
