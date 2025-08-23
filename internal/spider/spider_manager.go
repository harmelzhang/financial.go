package spider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/panjf2000/ants/v2"
	"harmel.cn/financial/internal/model"
	"harmel.cn/financial/internal/public"
	"harmel.cn/financial/internal/service"
	"harmel.cn/financial/internal/spider/response"
	"harmel.cn/financial/utils/http"
	"harmel.cn/financial/utils/slice"
	"harmel.cn/financial/utils/tools"
	"harmel.cn/financial/utils/xls"
)

// 爬虫管理器
type SpiderManager struct {
	// 抓取模式
	Mode string
	// 进度管理器
	progressManager *ProgressManager
	// 待处理任务通道
	taskChan chan PendingTask
	// 线程池
	pool *ants.Pool
	// 是否有在执行的任务
	hasTaskRunning atomic.Bool
	// 完成通知
	finishNotify chan bool
}

func NewSpiderManager(rootDir string) *SpiderManager {
	pool, err := ants.NewPool(public.SpiderExecutorPoolSize, ants.WithPreAlloc(true))
	if err != nil {
		panic(err)
	}

	return &SpiderManager{
		progressManager: NewProgressManager(rootDir),
		taskChan:        make(chan PendingTask, PENDING_TASKS_INIT_CAPACITY),
		pool:            pool,
		finishNotify:    make(chan bool),
	}
}

// 开启爬虫
func (s *SpiderManager) Start(ctx context.Context, mode string) (err error) {
	s.Mode = mode
	g.Log("spider").Debugf(ctx, "spider is running")

	// 加载历史进度
	err = s.progressManager.Load(ctx)
	if err != nil {
		g.Log("spider").Errorf(ctx, "ProgressManager.Load failed, err is %v", err)
		return err
	}

	// 如果到了五月一日，清空所有任务全部重跑（年报全部出了）
	if time.Now().Format("01-02") == "05-01" {
		s.progressManager.ClearTasks()
	}

	// 如果上次成功了，判断时间是否大于等于配置天数
	if s.progressManager.Done() {
		if time.Now().Unix()-s.progressManager.LastTS() >= public.SpiderTaskIntervalDays*24*3600 {
			s.progressManager.ClearTasks()
		} else {
			g.Log("spider").Debugf(ctx, "task finish, the time since the last successful task completion is less than %d days", public.SpiderTaskIntervalDays)
			return
		}
	}

	// 基础数据
	err = s.fetchIndexSample(ctx)
	if err != nil {
		g.Log("spider").Errorf(ctx, "fetch index sample data failed, err is %v", err)
		return
	}
	err = s.fetchCategory(ctx)
	if err != nil {
		g.Log("spider").Errorf(ctx, "fetch category data failed, err is %v", err)
		return
	}

	// 启动处理任务线程
	go s.doProcTaskWorker(ctx)

	// 等待完成
OUT:
	for {
		select {
		case <-s.finishNotify:
			g.Log("spider").Debug(ctx, "all task execute finish")
			break OUT
		default:
			time.Sleep(2 * time.Second)
			// 如果没有正在执行的线程并且通道里面没有待执行的任务
			if !s.hasTaskRunning.Load() && len(s.taskChan) == 0 {
				tasks := s.progressManager.UnexecutedTasks()
				if len(tasks) == 0 {
					s.progressManager.SetDone()
					err = s.progressManager.Save(ctx)
					if err != nil {
						g.Log("spider").Errorf(ctx, "save process failed, err is %v", err)
					}
					return
				}
				for _, task := range tasks {
					s.taskChan <- task
				}
			}
		}
	}

	return
}

// 最新指数样本信息
func (s *SpiderManager) fetchIndexSample(ctx context.Context) error {
	for typeCode := range public.IndexSampleType {
		g.Log("spider").Debugf(ctx, "start fetch %s index sample data", typeCode)

		// 请求数据
		url := fmt.Sprintf(public.UrlIndexSample, typeCode)
		client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
		body, _, err := client.Get(nil)
		if err != nil {
			g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
			continue
		}

		// 读取Excel
		items, err := xls.ReadXls(body, 0, 1)
		if err != nil {
			g.Log("spider").Errorf(ctx, "read xls failed, err is %v", err)
			continue
		}

		// 删除旧数据 & 插入新数据
		service.IndexSampleService.DeleteByType(ctx, typeCode)
		for _, item := range items {
			stockCode := item[4]
			indexSample := &model.IndexSample{
				TypeCode:  typeCode,
				StockCode: stockCode,
			}
			err = service.IndexSampleService.Insert(ctx, indexSample)
			if err != nil {
				g.Log("spider").Errorf(ctx, "insert index sample failed, TypeCode is %s StockCode is %s err is %v", typeCode, stockCode, err)
			}
		}

		g.Log("spider").Debugf(ctx, "fetch %s index sample data success", typeCode)
	}
	return nil
}

// 最新行业分类信息（含分类下的股票）
func (s *SpiderManager) fetchCategory(ctx context.Context) error {
	for typeName, typeValue := range public.CategoryType {
		g.Log("spider").Debugf(ctx, "start fetch %s catagory data", typeName)

		// 查询行业分类
		url := fmt.Sprintf(public.UrlCategory, typeValue)
		client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
		body, _, err := client.Get(nil)
		if err != nil {
			g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
			continue
		}

		categoryRes, err := http.ParseResponse[response.CategoryResult](body)
		if err != nil {
			g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
			continue
		}
		if categoryRes.Code == "200" && categoryRes.Success {
			// 删除数据库中指定类型的分类数据（同时会级联删除行业下股票信息）
			err = service.CategoryService.DeleteByType(ctx, typeName)
			if err != nil {
				g.Log("spider").Warningf(ctx, "delete category type %s data failed, err is %v", typeName, err)
				continue
			}
			// 递归插入新数据
			s.recursionCategorys(ctx, typeName, categoryRes.Data.MapList["4"])
		} else {
			g.Log("spider").Errorf(ctx, "fetch %s category data response error, code is %s", typeName, categoryRes.Code)
			continue
		}

		// 查询行业下的所有股票代码
		url = fmt.Sprintf(public.UrlCategoryStock, typeValue)
		client = http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
		body, _, err = client.Get(nil)
		if err != nil {
			g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
			continue
		}

		stockCodeRes, err := http.ParseResponse[response.StockCodeResult](body)
		if err != nil {
			g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
			continue
		}
		if stockCodeRes.Code == "200" && stockCodeRes.Success {
			// 插入新数据
			for _, stock := range stockCodeRes.Data.List {
				var categoryCode string
				if stock.CicsLeve1Code != "" {
					// 中证
					if stock.CicsLeve4Code == "99999999" {
						continue
					}
					categoryCode = stock.CicsLeve4Code
				} else {
					// 证监会
					if stock.CsrcLeve2Code == "" {
						// FIX 证券会暂时没对新三板股票进行分类，后续待优化
						continue
					}
					categoryCode = stock.CsrcLeve1Code + stock.CsrcLeve2Code
				}
				categoryStockCode := &model.CategoryStockCode{
					CategoryCode: categoryCode,
					StockCode:    stock.Code,
				}
				err := service.CategoryStockCodeService.Insert(ctx, categoryStockCode)
				if err != nil {
					g.Log("spider").Warningf(ctx, "insert categroy stock code failed, err is %v", err)
				}
				//  丢入任务列表
				task := PendingTask{Id: stock.Code}
				exist := s.progressManager.PutTask(task)
				if !exist {
					s.taskChan <- task
				}
			}
		} else {
			g.Log("spider").Errorf(ctx, "fetch %s category stock code data response error, code is %s", typeName, categoryRes.Code)
			continue
		}

		g.Log("spider").Debugf(ctx, "fetch %s catagory data success", typeName)
	}

	return nil
}

// 递归查询分类
func (s *SpiderManager) recursionCategorys(ctx context.Context, typeName string, categorys []response.Category) {
	if len(categorys) == 0 {
		return
	}

	for order, category := range categorys {
		mCategory := &model.Category{
			Type:         typeName,
			Code:         category.Id,
			Name:         category.Name,
			Level:        category.Level,
			DisplayOrder: order + 1,
		}
		if category.ParentId != "" {
			mCategory.ParentCode = category.ParentId
		}
		// 插入数据库
		err := service.CategoryService.Insert(ctx, mCategory)
		if err != nil {
			g.Log("spider").Warningf(ctx, "insert category data failed, err is %v", err)
			continue
		}
		if len(category.Children) != 0 {
			s.recursionCategorys(ctx, typeName, category.Children)
		}
	}
}

// 处理任务
func (s *SpiderManager) doProcTaskWorker(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			g.Log("spiser").Critical(ctx, "doProcTaskWorker panic: %v", err)
		}
	}()

	for {
		task := <-s.taskChan
		err := s.executeTask(ctx, task)
		if err != nil {
			g.Log("spider").Errorf(ctx, "execute task failed, err is %v", err)
		}
	}
}

// 根据股票代码和报告期查询索引位置
func (s *SpiderManager) findFinancialIndex(stockCode, reportDate string, financials []*model.Financial) int {
	for idx, financial := range financials {
		if financial.StockCode == stockCode && financial.ReportDate == reportDate {
			return idx
		}
	}
	return -1
}

// 执行实际任务
func (s *SpiderManager) executeTask(ctx context.Context, task PendingTask) (err error) {
	err = s.pool.Submit(func() {
		s.hasTaskRunning.Store(true)
		defer func() {
			s.hasTaskRunning.Store(false)
		}()

		// 是否处理完
		isFinished, err := s.progressManager.TaskStatus(ctx, task.Id)
		if err != nil {
			g.Log("spider").Errorf(ctx, "query task status failed, err is %v", err)
			return
		}
		if isFinished {
			return
		}

		g.Log("spider").Debugf(ctx, "start execute task %s", task.Id)

		// 基本信息
		stock, err := s.fetchStockBaseInfo(ctx, task.Id)
		if err != nil {
			g.Log("spider").Errorf(ctx, "fetch stock %s base info failed, err is %v", stock.Code, err)
			return
		}

		// 查询所有报告期
		reportDates, err := s.queryAllReportData(ctx, stock)
		if err != nil {
			g.Log("spider").Errorf(ctx, "fetch stock %s report date info failed, err is %v", stock.Code, err)
			return
		}

		// 初始化操作
		financials := make([]*model.Financial, 0, len(reportDates))
		for _, reportDate := range reportDates {
			ymd := strings.Split(reportDate, "-")
			financial := &model.Financial{
				StockCode:  stock.Code,
				Year:       ymd[0],
				ReportDate: reportDate,
			}
			switch ymd[1] {
			case "03":
				financial.ReportType = public.ReportTypeQ1
			case "06":
				financial.ReportType = public.ReportTypeH1
			case "09":
				financial.ReportType = public.ReportTypeQ3
			case "12":
				financial.ReportType = public.ReportTypeFY
			default:
				financial.ReportType = public.ReportTypeO
			}
			financials = append(financials, financial)
		}

		// 分页查询财报
		reportDatePages, totalPages := slice.ArraySlice(reportDates, public.QueryReportPageSize)
		for i, reportDates := range reportDatePages {
			g.Log("spider").Debugf(ctx, "fetch stock %s report info page %d/%d", stock.Code, i+1, totalPages)
			queryDates := strings.Join(reportDates, ",")
			// 现金流量表
			err = s.fetchCashFlowSheet(ctx, stock, queryDates, financials)
			if err != nil {
				return
			}
			// 资产负债表
			err = s.fetchBalanceSheet(ctx, stock, queryDates, financials)
			if err != nil {
				return
			}
			// 利润表
			err = s.fetchIncomeSheet(ctx, stock, queryDates, financials)
			if err != nil {
				return
			}
		}
		// 分红数据
		err = s.fetchDividendData(ctx, stock, financials)
		if err != nil {
			return
		}

		// 计算现金流量允当比率（年报）
		s.calcCashFlowAdequacyRatio(ctx, financials)

		// 插入或更新数据库
		for _, financial := range financials {
			err = service.FinancialService.Replace(ctx, financial)
			if err != nil {
				g.Log("spider").Errorf(ctx, "insert financial data failed, err is %v", err)
				return
			}
		}

		// 比率
		err = s.calcFinancialRatios(ctx, stock)
		if err != nil {
			g.Log("spider").Errorf(ctx, "calc %s financial ratios failed, err is %v", stock.Code, err)
			return
		}

		// 标记完成
		s.progressManager.MarkTask(ctx, task.Id, true)

		// 写入磁盘
		err = s.progressManager.Save(ctx)
		if err != nil {
			g.Log("spider").Errorf(ctx, "save process failed, err is %v", err)
			return
		}

		g.Log("spider").Debugf(ctx, "task %s execute success", task.Id)

		// 通知完成
		if s.progressManager.Done() {
			s.finishNotify <- true
		}
	})
	return
}

// 查询股票市场
func (s *SpiderManager) queryStockMarketPlace(stockCode string) (string, string) {
	name, shortName := "", ""
	stockCodePrefix := stockCode[0:2]
	if slice.IndexOf(public.ShanghaiMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "上海", "SH"
	} else if slice.IndexOf(public.ShenzhenMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "深圳", "SZ"
	} else if slice.IndexOf(public.BeijingMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "北京", "BJ"
	}
	return name, shortName
}

// 基本信息
func (s *SpiderManager) fetchStockBaseInfo(ctx context.Context, stockCode string) (stock *model.Stock, err error) {
	if s.Mode == public.SpiderModeDiff {
		stock, err = service.StockService.FindStockByCode(ctx, stockCode)
		if err != nil {
			g.Log("spider").Errorf(ctx, "find stock %s by code failed, err is %v", stockCode, err)
			return
		}
		if stock != nil {
			return
		}
	}

	marketName, marketShortName := s.queryStockMarketPlace(stockCode)

	// 公司类型
	url := fmt.Sprintf(public.UrlStockCompanyType, stockCode)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}

	companyType, companyTypeCode := "普通", "4"
	companyTypeRes, err := http.ParseResponse[response.CompanyTypeResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
	}
	if companyTypeRes.Success && companyTypeRes.Code == 0 {
		if companyTypeRes.Result.Count != 0 {
			companyType = companyTypeRes.Result.Data[0].Type
			companyTypeCode = companyTypeRes.Result.Data[0].TypeCode
		}
	}

	// 主营业务
	url = fmt.Sprintf(public.UrlStockMainBusiness, stockCode)
	client = http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err = client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
	}

	mainBusinessResult, err := http.ParseResponse[response.MainBusinessResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}

	mainBusiness := ""
	if mainBusinessResult.Code == 0 && mainBusinessResult.Success {
		mainBusiness = mainBusinessResult.Result.Data[0].Info
	} else {
		g.Log("spider").Errorf(ctx, "fetch %s main business data response error, code is %d", stockCode, mainBusinessResult.Code)
	}

	// 主要信息
	url = fmt.Sprintf(public.UrlStockBaseInfo, marketShortName, stockCode)
	client = http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err = client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}

	baseInfoRes, err := http.ParseResponse[response.StockBaseInfoResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	baseInfo := baseInfoRes.BaseInfo[0]
	listingInfo := baseInfoRes.ListingInfo[0]

	stock = &model.Stock{
		Code:            stockCode,
		Name:            baseInfo.Name,
		NamePinYin:      tools.PinyinFirstWord(baseInfo.Name),
		BeforeName:      baseInfo.BeforeName,
		CompanyName:     baseInfo.CompanyName,
		CompanyType:     companyType,
		CompanyTypeCode: companyTypeCode,
		CompanyProfile:  strings.TrimSpace(baseInfo.CompanyProfile),
		Region:          baseInfo.Region,
		Address:         baseInfo.Address,
		Website:         baseInfo.Website,
		MainBusiness:    mainBusiness,
		BusinessScope:   baseInfo.BusinessScope,
		CreateDate:      listingInfo.CreateDate,
		ListingDate:     listingInfo.ListingDate,
		LawFirm:         baseInfo.LawFirm,
		AccountingFirm:  baseInfo.AccountingFirm,
		MarketPlace:     marketName,
	}
	if strings.TrimSpace(baseInfo.BeforeName) == "" {
		stock.BeforeName = nil
	}
	if stock.BeforeName != nil {
		stock.BeforeName = strings.ReplaceAll(fmt.Sprint(stock.BeforeName), "→", "、")
	}
	err = service.StockService.Replace(ctx, stock)
	if err != nil {
		g.Log("spider").Errorf(ctx, "replace db stock failed, err is %v", err)
		return
	}

	return
}

// 查询所有报告期
func (s *SpiderManager) queryAllReportData(ctx context.Context, stock *model.Stock) (reportDates []string, err error) {
	_, shortMarketName := s.queryStockMarketPlace(stock.Code)

	fetchReportDates := make([]string, 0)
	appendReportDate := func(reportDateRes *response.ReportDateResult) {
		for _, item := range reportDateRes.Data {
			date := strings.Split(item.Date, " ")[0]
			if slice.IndexOf(fetchReportDates, date) == -1 {
				fetchReportDates = append(fetchReportDates, date)
			}
		}
	}

	dbReportDatas := make([]string, 0)
	if s.Mode == public.SpiderModeDiff {
		dbReportDatas, err = service.FinancialService.GetReportDates(ctx, stock.Code)
		if err != nil {
			g.Log("spider").Errorf(ctx, "get report dates failed, stock code is %s err is %v", stock.Code, err)
			return nil, err
		}
	}

	// 资产负债表
	url := fmt.Sprintf(public.UrlBalanceSheetReport, stock.CompanyTypeCode, shortMarketName, stock.Code)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}
	reportDateRes, err := http.ParseResponse[response.ReportDateResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	appendReportDate(reportDateRes)

	// 利润表
	url = fmt.Sprintf(public.UrlIncomeSheetReport, stock.CompanyTypeCode, shortMarketName, stock.Code)
	client = http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err = client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}
	reportDateRes, err = http.ParseResponse[response.ReportDateResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	appendReportDate(reportDateRes)

	// 现金流量表
	url = fmt.Sprintf(public.UrlCashFlowSheetReport, stock.CompanyTypeCode, shortMarketName, stock.Code)
	client = http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err = client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return
	}
	reportDateRes, err = http.ParseResponse[response.ReportDateResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return
	}
	appendReportDate(reportDateRes)

	for _, reportDate := range fetchReportDates {
		if slice.IndexOf(dbReportDatas, reportDate) != -1 {
			continue
		}
		reportDates = append(reportDates, reportDate)
	}

	return
}

// 现金流量表
func (s *SpiderManager) fetchCashFlowSheet(ctx context.Context, stock *model.Stock, queryDates string, financials []*model.Financial) error {
	_, marketShortName := s.queryStockMarketPlace(stock.Code)
	url := fmt.Sprintf(public.UrlCashFlowSheet, stock.CompanyTypeCode, queryDates, marketShortName, stock.Code)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return err
	}
	financialRes, err := http.ParseResponse[response.FinancialResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return err
	}
	if financialRes.Type == "1" || financialRes.Status == 1 {
		g.Log("spider").Warningf(ctx, "fetch %s cash flow sheet data response error, type is %s status is %d, url is %s", stock.Code, financialRes.Type, financialRes.Status, url)
		return err
	}

	for _, sheet := range financialRes.Data {
		reportDate := strings.Split(sheet.ReportDate, " ")[0]
		idx := s.findFinancialIndex(stock.Code, reportDate, financials)
		if idx == -1 {
			continue
		}
		financial := financials[idx]

		financial.Ocf = sheet.Ocf
		financial.Cfi = sheet.Cfi
		financial.Cff = sheet.Cff
		financial.AssignDividendPorfit = sheet.AssignDividendPorfit
		financial.AcquisitionAssets = sheet.AcquisitionAssets
		financial.InventoryLiquidating = sheet.InventoryLiquidating
	}
	return nil
}

// 资产负债表
func (s *SpiderManager) fetchBalanceSheet(ctx context.Context, stock *model.Stock, queryDates string, financials []*model.Financial) error {
	_, marketShortName := s.queryStockMarketPlace(stock.Code)
	url := fmt.Sprintf(public.UrlBalanceSheet, stock.CompanyTypeCode, queryDates, marketShortName, stock.Code)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return err
	}
	financialRes, err := http.ParseResponse[response.FinancialResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return err
	}
	if financialRes.Type == "1" || financialRes.Status == 1 {
		g.Log("spider").Warningf(ctx, "fetch %s balance sheet data response error, type is %s status is %d, url is %s", stock.Code, financialRes.Type, financialRes.Status, url)
		return err
	}

	for _, sheet := range financialRes.Data {
		reportDate := strings.Split(sheet.ReportDate, " ")[0]
		idx := s.findFinancialIndex(stock.Code, reportDate, financials)
		if idx == -1 {
			continue
		}
		financial := financials[idx]

		financial.MonetaryFund = sheet.MonetaryFund
		financial.TradeFinassetNotfvtpl = sheet.TradeFinassetNotfvtpl
		financial.TradeFinasset = sheet.TradeFinasset
		financial.DeriveFinasset = sheet.DeriveFinasset

		financial.FixedAsset = sheet.FixedAsset
		financial.Cip = sheet.Cip

		financial.CaTotal = sheet.CaTotal
		financial.NcaTotal = sheet.NcaTotal
		financial.ClTotal = sheet.ClTotal
		financial.NclTotal = sheet.NclTotal
		financial.Inventory = sheet.Inventory
		financial.AccountsRece = sheet.AccountsRece
		financial.AccountsPayable = sheet.AccountsPayable
	}
	return nil
}

// 利润表
func (s *SpiderManager) fetchIncomeSheet(ctx context.Context, stock *model.Stock, queryDates string, financials []*model.Financial) error {
	_, marketShortName := s.queryStockMarketPlace(stock.Code)
	url := fmt.Sprintf(public.UrlIncomeSheet, stock.CompanyTypeCode, queryDates, marketShortName, stock.Code)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return err
	}
	financialRes, err := http.ParseResponse[response.FinancialResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return err
	}
	if financialRes.Type == "1" || financialRes.Status == 1 {
		g.Log("spider").Warningf(ctx, "fetch %s balance sheet data response error, type is %s status is %d, url is %s", stock.Code, financialRes.Type, financialRes.Status, url)
		return err
	}

	for _, sheet := range financialRes.Data {
		reportDate := strings.Split(sheet.ReportDate, " ")[0]
		idx := s.findFinancialIndex(stock.Code, reportDate, financials)
		if idx == -1 {
			continue
		}
		financial := financials[idx]

		financial.Np = sheet.Np
		financial.Oi = sheet.Oi
		financial.Coe = sheet.Coe
		financial.CoeTotal = sheet.CoeTotal
		financial.Eps = sheet.Eps
	}
	return nil
}

// 分红数据
func (s *SpiderManager) fetchDividendData(ctx context.Context, stock *model.Stock, financials []*model.Financial) error {
	url := fmt.Sprintf(public.UrlDividend, stock.Code)
	client := http.New(url, time.Duration(public.SpiderTimtout)*time.Second)
	body, _, err := client.Get(nil)
	if err != nil {
		g.Log("spider").Errorf(ctx, "request url failed, err is %v", err)
		return err
	}
	dividendRes, err := http.ParseResponse[response.DividendResult](body)
	if err != nil {
		g.Log("spider").Errorf(ctx, "parse response failed, err is %v", err)
		return err
	}
	if dividendRes.Code == 0 && dividendRes.Success {
		for _, dividend := range dividendRes.Result.Data {
			reportDate := dividend.Year + "-12-31"
			idx := s.findFinancialIndex(stock.Code, reportDate, financials)
			if idx == -1 {
				continue
			}
			financial := financials[idx]
			financial.Dividend = dividend.Money
		}
	} else {
		g.Log("spider").Errorf(ctx, "fetch %s dividend data response error, code is %d message is %s", stock.Code, dividendRes.Code, dividendRes.Message)
		return err
	}
	return nil
}

// TODO 计算现金流量允当比率（年报）
func (s *SpiderManager) calcCashFlowAdequacyRatio(ctx context.Context, financials []*model.Financial) {
	if len(financials) == 0 {
		return
	}

	// 过滤出报告类型为年报的数据
	annualFinancials := make([]*model.Financial, 0)
	for _, financial := range financials {
		if financial.ReportType == public.ReportTypeFY {
			annualFinancials = append(annualFinancials, financial)
		}
	}
	// 没有年报直接跳过
	if len(annualFinancials) == 0 {
		return
	}

	// 查询数据库
	dbFinancials, err := service.FinancialService.GetByType(ctx, financials[0].StockCode, public.ReportTypeFY)
	if err != nil {
		g.Log("spider").Errorf(ctx, "get %s financial data by type failed, err is %v", financials[0].StockCode, err)
		return
	}

	// 构建五年数据
	for _, financial := range annualFinancials {
		reportDatas := make([]*model.Financial, 0, 5)
		iYear, _ := strconv.Atoi(financial.Year)
		for i := iYear; i >= iYear-4; i-- {
			ymd := fmt.Sprintf("%d-12-31", i)
			// 先从新的数据里面找，找不到再从数据库找
			var found bool
			for _, financial := range annualFinancials {
				if financial.ReportDate == ymd {
					reportDatas = append(reportDatas, financial)
					found = true
					break
				}
			}
			if !found {
				// 从数据库找
				for _, dbFinancial := range dbFinancials {
					if dbFinancial.ReportDate == ymd {
						dbFinancial.Ocf = dbFinancial.Ocf.(*gvar.Var).Float64()
						dbFinancial.AcquisitionAssets = dbFinancial.AcquisitionAssets.(*gvar.Var).Float64()
						dbFinancial.AssignDividendPorfit = dbFinancial.AssignDividendPorfit.(*gvar.Var).Float64()
						dbFinancial.InventoryLiquidating = dbFinancial.InventoryLiquidating.(*gvar.Var).Float64()
						reportDatas = append(reportDatas, dbFinancial)
						break
					}
				}
			}
		}
		if len(reportDatas) != 5 {
			continue
		}
		/**
		现金流量允当比率 = 最近五年营业活动净现金流/最近五年(资本支出+现金股利+存货增加)
		计算：营业活动现金流量 / (购建固定资产、无形资产和其他长期资产支付的现金 + 分配股利、利润或偿付利息支付的现金 - 存货减少额)
		*/
		var numerator, denominator float64
		// 如果有一年没数据就跳过
		hasNull := false
		for _, data := range reportDatas {
			if data == nil || data.Ocf == nil {
				hasNull = true
				break
			}
			numerator += data.Ocf.(float64)
			if data.AcquisitionAssets != nil {
				denominator += data.AcquisitionAssets.(float64)
			}
			if data.AssignDividendPorfit != nil {
				denominator += data.AssignDividendPorfit.(float64)
			}
			if data.InventoryLiquidating != nil {
				denominator -= data.InventoryLiquidating.(float64)
			}
		}

		if denominator == 0 {
			continue
		}
		if !hasNull {
			financial.CashFlowAdequacyRatio = fmt.Sprintf("%.2f", numerator/denominator*10000/100)
		}
	}
}

// 计算财务比率
func (s *SpiderManager) calcFinancialRatios(ctx context.Context, stock *model.Stock) error {
	sql := `
		UPDATE financial
		SET
		    oi = IF(oi = 0, NULL, oi),
		    coe = IF(coe = 0, NULL, coe),
		    np = IF(np = 0, NULL, np),

		    asset_total = IF(ca_total IS NULL AND nca_total IS NULL, NULL, IFNULL(ca_total, 0) + IFNULL(nca_total, 0)),
		    asset_total = IF(asset_total = 0, NULL, asset_total),

		    liability_total = IF(cl_total IS NULL AND ncl_total IS NULL, NULL, IFNULL(cl_total, 0) + IFNULL(ncl_total, 0)),
		    np_ratio = ROUND(np / oi * 100, 2),
		    dividend_ratio = ROUND(dividend / np * 100, 2),
		    oi_ratio = ROUND((oi - coe) / oi * 100, 2),
		    operating_profit_ratio = ROUND((oi - coe_total) / oi * 100, 2),
		    operating_safety_ratio = IF(oi_ratio = 0, NUll, ROUND(operating_profit_ratio / oi_ratio * 100, 2)),
		    cash_equivalent_ratio = ROUND((monetary_fund + IFNULL(IFNULL(trade_finasset, trade_finasset_notfvtpl), 0) + IFNULL(derive_finasset, 0)) / asset_total * 100, 2),
		    cash_ratio = IF(cl_total = 0, NULL, ROUND(monetary_fund / cl_total * 100, 2)),
		    ca_ratio = ROUND(ca_total / asset_total * 100, 2),
		    cl_ratio = ROUND(cl_total / asset_total * 100, 2),
		    ncl_ratio = ROUND(ncl_total / asset_total * 100, 2),
		    debt_ratio = ROUND((cl_total + ncl_total) / asset_total * 100, 2),
		    long_term_funds_ratio = IF((fixed_asset + cip) = 0, NULL, ROUND((ncl_total + (asset_total - liability_total)) / (fixed_asset + cip) * 100, 2)),
		    equity_ratio = ROUND(100 - debt_ratio, 2),
		    equity_multiplier = IF((asset_total - liability_total) = 0, NUll, ROUND(asset_total / (asset_total - liability_total), 2)),
		    capitalization_ratio = IF((asset_total - liability_total) = 0, NULL, ROUND((cl_total + ncl_total) / (asset_total - liability_total) * 100, 2)),
		    inventory_ratio = ROUND(inventory / asset_total * 100, 2),
		    accounts_rece_ratio = ROUND(accounts_rece / asset_total * 100, 2),
		    accounts_payable_ratio = ROUND(accounts_payable / asset_total * 100, 2),
		    current_ratio = IF(cl_total = 0, NULL, ROUND(ca_total / cl_total * 100, 2)),
		    quick_ratio = IF(cl_total = 0, NULL, ROUND((ca_total - inventory) / cl_total * 100, 2)),
		    roe = IF((asset_total - cl_total - ncl_total) = 0, NUll, ROUND(np / (asset_total - cl_total - ncl_total) * 100, 2)),
		    roa = ROUND(np / asset_total * 100, 2),
		    accounts_rece_turnover_ratio = ROUND(oi / IF(accounts_rece = 0, NULL, accounts_rece), 2),
		    average_cash_receipt_days = ROUND(360 / IF(accounts_rece_turnover_ratio = 0, NULL, accounts_rece_turnover_ratio), 2),
		    inventory_turnover_ratio = ROUND(coe / IF(inventory = 0, NULL, inventory), 2),
		    average_sales_days = ROUND(360 / IF(inventory_turnover_ratio = 0, NULL, inventory_turnover_ratio), 2),
		    immovables_turnover_ratio = IF((fixed_asset + cip) = 0, NULL, ROUND(oi / (fixed_asset + cip), 2)),
		    total_asset_turnover_ratio = ROUND(oi / asset_total, 2),
		    cash_flow_ratio = IF(cl_total = 0, NULL, ROUND(ocf / cl_total * 100, 2)),
		    cash_reinvestment_ratio = IF((asset_total - cl_total) = 0, NULL, ROUND((ocf - assign_dividend_porfit) / (asset_total - cl_total) * 100, 2)),
		    profit_cash_ratio = ROUND(ocf / np * 100, 2)
		WHERE stock_code = ?
	`
	_, err := g.DB().Exec(ctx, sql, stock.Code)
	return err
}
