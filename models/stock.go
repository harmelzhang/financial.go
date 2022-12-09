package models

import (
	"bytes"
	"encoding/json"
	"financial/config"
	"financial/utils/http"
	"financial/utils/tools"
	"fmt"
	"github.com/antchfx/htmlquery"
	"log"
	"strings"
)

// Stock 股票信息
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
	BusinessScope       string   // 经营范围
	DateOfIncorporation string   // 成立日期
	ListingDate         string   // 上市日期
	MainUnderwriter     string   // 主承销商
	Sponsor             string   // 上市保荐人
	AccountingFirm      string   // 会计师事务所
	MarketPlace         string   // 交易市场（上海、深圳、北京）
	Category            Category // 所属行业分类
}

// GetStockMarketPlace 查询股票交易市场名称和简称（SH、SZ、BJ）
func (stock *Stock) GetStockMarketPlace() (string, string) {
	name, shortName := "", ""
	stockCodePrefix := stock.Code[0:2]
	if tools.IndexOf(config.ShanghaiMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "上海", "SH"
	} else if tools.IndexOf(config.ShenzhenMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "深圳", "SZ"
	} else if tools.IndexOf(config.BeijingMarketPrefixs, stockCodePrefix) != -1 {
		name, shortName = "北京", "BJ"
	}
	return name, shortName
}

// BuildStockInfo 查询股票信息
func (stock *Stock) BuildStockInfo() {
	log.Printf("查询股票 [%s] 基本信息", stock.Code)

	marketPlaceName, marketPlaceShortName := stock.GetStockMarketPlace()

	url := fmt.Sprintf(config.FetchStockBaseInfoNeteaseUrl, stock.Code)
	if marketPlaceShortName == "BJ" {
		url = fmt.Sprintf(config.FetchStockBaseInfoEastmoneyUrl, marketPlaceShortName, stock.Code)
	}

	data := http.Get(url)
	root, err := htmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("解析HTML出错 : %s", err)
	}

	// 获取指定 XPTAH 路径节点的文本元素
	fetchData := func(xpath string) string {
		result := ""
		node := htmlquery.Find(root, xpath)[0].FirstChild
		if node != nil {
			result = node.Data
		}
		result = strings.Trim(result, " ")
		if result == "--" {
			result = ""
		}
		return result
	}

	fetchEastmoneyData := func() {
		jsonData := make(map[string]interface{})
		err := json.Unmarshal(data, &jsonData)
		if err != nil {
			log.Fatalf("解析JSON数据失败 : %s", err)
		}

		baseInfo := jsonData["jbzl"].([]interface{})[0].(map[string]interface{})    // 基本资料
		publishInfo := jsonData["fxxg"].([]interface{})[0].(map[string]interface{}) // 发行相关

		trimMapKeyVal := func(m map[string]interface{}, key string) string {
			if m[key] == nil {
				return ""
			}
			val := strings.Trim(m[key].(string), "")
			if key == "FOUND_DATE" || key == "LISTING_DATE" {
				val = strings.Split(val, " ")[0]
			}
			return val
		}

		stock.StockName = trimMapKeyVal(baseInfo, "STR_NAMEA")
		stock.CompanyName = trimMapKeyVal(baseInfo, "ORG_NAME")
		stock.Organization = ""
		stock.Region = trimMapKeyVal(baseInfo, "PROVINCE")
		stock.Address = trimMapKeyVal(baseInfo, "ADDRESS")
		stock.WebSite = trimMapKeyVal(baseInfo, "ORG_WEB")
		stock.MainBusiness = ""
		stock.BusinessScope = trimMapKeyVal(baseInfo, "BUSINESS_SCOPE")
		stock.DateOfIncorporation = trimMapKeyVal(publishInfo, "FOUND_DATE")
		stock.ListingDate = trimMapKeyVal(publishInfo, "LISTING_DATE")
		stock.MainUnderwriter = ""
		stock.Sponsor = ""
		stock.AccountingFirm = trimMapKeyVal(baseInfo, "ACCOUNTFIRM_NAME")
	}

	fetchNeteaseData := func() {
		stock.StockName = fetchData("/html/body/div[2]/div[4]/table/tbody/tr[2]/td[2]")
		stock.CompanyName = fetchData("/html/body/div[2]/div[4]/table/tbody/tr[3]/td[2]")
		stock.Organization = fetchData("/html/body/div[2]/div[4]/table/tbody/tr[1]/td[2]")
		stock.Region = fetchData("/html/body/div[2]/div[4]/table/tbody/tr[1]/td[4]")
		stock.Address = fetchData("/html/body/div[2]/div[4]/table/tbody/tr[2]/td[4]")
		stock.WebSite = fetchData("/html/body/div[2]/div[4]/table/tbody/tr[8]/td[2]")
		stock.MainBusiness = fetchData("/html/body/div[2]/div[4]/table/tbody/tr[9]/td[2]")
		stock.BusinessScope = fetchData("/html/body/div[2]/div[4]/table/tbody/tr[10]/td[2]")
		stock.DateOfIncorporation = fetchData("/html/body/div[2]/div[5]/table/tbody/tr[1]/td[2]")
		stock.ListingDate = fetchData("/html/body/div[2]/div[5]/table/tbody/tr[2]/td[2]")
		stock.MainUnderwriter = fetchData("/html/body/div[2]/div[5]/table/tbody/tr[16]/td[2]")
		stock.Sponsor = fetchData("/html/body/div[2]/div[5]/table/tbody/tr[17]/td[2]")
		stock.AccountingFirm = fetchData("/html/body/div[2]/div[5]/table/tbody/tr[18]/td[2]")
	}

	if marketPlaceShortName == "BJ" {
		fetchEastmoneyData()
	} else {
		fetchNeteaseData()
	}

	stock.StockNamePinyin = tools.GetPinyinFirstWord(stock.StockName)
	stock.MarketPlace = marketPlaceName

	stock.IntoDb()
}

// IntoDb 更新数据库
func (stock *Stock) IntoDb() {

}

// FetchFinancialInfo 获取财务信息
func (stock *Stock) FetchFinancialInfo() {

}
