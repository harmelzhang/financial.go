package public

var (
	// 任务间隔天数
	SpiderTaskIntervalDays int64 = 0
	// 线程池执行器初始化个数
	SpiderExecutorPoolSize = 3
	// 请求超时时长（秒）
	SpiderTimtout = 3
	// 一次查询报告的数量
	QueryReportPageSize = 5
)

var (
	// 指数样本类型
	IndexSampleType = make(map[string]string)

	// 上交所股票前缀
	ShanghaiMarketPrefixs []string
	// 深交所股票前缀
	ShenzhenMarketPrefixs []string
	// 北交所股票前缀
	BeijingMarketPrefixs []string
)

var (
	// 行业分类类型（CSRC：证监会、CICS：中证）
	CategoryType = map[string]string{"CSRC": "1", "CICS": "2"}
)

const (
	// 查询所有数据
	SpiderModeAll = "all"
	// 查询新报告数据
	SpiderModeDiff = "diff"
)

const (
	// 一季报
	ReportTypeQ1 = "Q1"
	// 半年报
	ReportTypeH1 = "H1"
	// 三季报
	ReportTypeQ3 = "Q3"
	// 年报
	ReportTypeFY = "FY"
	// 其他报告
	ReportTypeO = "O"
)
