package models

import (
	"bytes"
	"encoding/json"
	"financial/config"
	"financial/utils/db"
	"financial/utils/http"
	"financial/utils/tools"
	"fmt"
	"github.com/antchfx/htmlquery"
	"log"
	"strings"
)

// Stock 股票信息
type Stock struct {
	Code                *Value   // 股票代码
	StockName           *Value   // 股票名称
	StockNamePinyin     *Value   // 股票名称（拼音）
	CompanyName         *Value   // 公司名称
	Organization        *Value   // 组织形式（民营、国营、中外合资...）
	Region              *Value   // 地域（省份）
	Address             *Value   // 办公地址
	WebSite             *Value   // 公司网站
	MainBusiness        *Value   // 主营业务
	BusinessScope       *Value   // 经营范围
	DateOfIncorporation *Value   // 成立日期
	ListingDate         *Value   // 上市日期
	MainUnderwriter     *Value   // 主承销商
	Sponsor             *Value   // 上市保荐人
	AccountingFirm      *Value   // 会计师事务所
	MarketPlace         *Value   // 交易市场（上海、深圳、北京）
	Category            Category // 所属行业分类
}

// GetStockMarketPlace 查询股票交易市场名称和简称（SH、SZ、BJ）
func (stock *Stock) GetStockMarketPlace() (string, string) {
	name, shortName := "", ""
	stockCodePrefix := stock.Code.String()[0:2]
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
	log.Printf("查询股票 [%s %s] 基本信息", stock.Category.Name, stock.Code.String())

	marketPlaceName, marketPlaceShortName := stock.GetStockMarketPlace()

	url := fmt.Sprintf(config.FetchStockBaseInfoNeteaseUrl, stock.Code.String())
	if marketPlaceShortName == "BJ" {
		url = fmt.Sprintf(config.FetchStockBaseInfoEastmoneyUrl, marketPlaceShortName, stock.Code.String())
	}

	data := http.Get(url)
	root, err := htmlquery.Parse(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("解析HTML出错 : %s", err)
	}

	fetchEastmoneyData := func() {
		jsonData := make(map[string]interface{})
		err := json.Unmarshal(data, &jsonData)
		if err != nil {
			log.Fatalf("解析JSON数据失败 : %s", err)
		}

		baseInfo := jsonData["jbzl"].([]interface{})[0].(map[string]interface{})    // 基本资料
		publishInfo := jsonData["fxxg"].([]interface{})[0].(map[string]interface{}) // 发行相关

		trimMapKeyVal := func(m map[string]interface{}, key string) interface{} {
			if m[key] == nil {
				return nil
			}
			val := strings.Trim(m[key].(string), "")
			if key == "FOUND_DATE" || key == "LISTING_DATE" {
				val = strings.Split(val, " ")[0]
			}
			if val == "" {
				return nil
			}
			return val
		}

		stock.StockName = NewValue(trimMapKeyVal(baseInfo, "STR_NAMEA"))
		stock.CompanyName = NewValue(trimMapKeyVal(baseInfo, "ORG_NAME"))
		// stock.Organization = ""
		stock.Region = NewValue(trimMapKeyVal(baseInfo, "PROVINCE"))
		stock.Address = NewValue(trimMapKeyVal(baseInfo, "ADDRESS"))
		stock.WebSite = NewValue(trimMapKeyVal(baseInfo, "ORG_WEB"))
		// stock.MainBusiness = ""
		stock.BusinessScope = NewValue(trimMapKeyVal(baseInfo, "BUSINESS_SCOPE"))
		stock.DateOfIncorporation = NewValue(trimMapKeyVal(publishInfo, "FOUND_DATE"))
		stock.ListingDate = NewValue(trimMapKeyVal(publishInfo, "LISTING_DATE"))
		// stock.MainUnderwriter = ""
		// stock.Sponsor = ""
		stock.AccountingFirm = NewValue(trimMapKeyVal(baseInfo, "ACCOUNTFIRM_NAME"))
	}

	fetchNeteaseData := func() {
		// 获取指定 XPTAH 路径节点的文本元素
		fetchNodeData := func(xpath string) interface{} {
			result := ""
			node := htmlquery.Find(root, xpath)[0].FirstChild
			if node != nil {
				result = node.Data
			}
			result = strings.Trim(result, " ")
			if result == "--" {
				result = ""
			}
			if result == "" {
				return nil
			}
			return result
		}

		stock.StockName = NewValue(fetchNodeData("/html/body/div[2]/div[4]/table/tbody/tr[2]/td[2]"))
		stock.CompanyName = NewValue(fetchNodeData("/html/body/div[2]/div[4]/table/tbody/tr[3]/td[2]"))
		stock.Organization = NewValue(fetchNodeData("/html/body/div[2]/div[4]/table/tbody/tr[1]/td[2]"))
		stock.Region = NewValue(fetchNodeData("/html/body/div[2]/div[4]/table/tbody/tr[1]/td[4]"))
		stock.Address = NewValue(fetchNodeData("/html/body/div[2]/div[4]/table/tbody/tr[2]/td[4]"))
		stock.WebSite = NewValue(fetchNodeData("/html/body/div[2]/div[4]/table/tbody/tr[8]/td[2]"))
		stock.MainBusiness = NewValue(fetchNodeData("/html/body/div[2]/div[4]/table/tbody/tr[9]/td[2]"))
		stock.BusinessScope = NewValue(fetchNodeData("/html/body/div[2]/div[4]/table/tbody/tr[10]/td[2]"))
		stock.DateOfIncorporation = NewValue(fetchNodeData("/html/body/div[2]/div[5]/table/tbody/tr[1]/td[2]"))
		stock.ListingDate = NewValue(fetchNodeData("/html/body/div[2]/div[5]/table/tbody/tr[2]/td[2]"))
		stock.MainUnderwriter = NewValue(fetchNodeData("/html/body/div[2]/div[5]/table/tbody/tr[16]/td[2]"))
		stock.Sponsor = NewValue(fetchNodeData("/html/body/div[2]/div[5]/table/tbody/tr[17]/td[2]"))
		stock.AccountingFirm = NewValue(fetchNodeData("/html/body/div[2]/div[5]/table/tbody/tr[18]/td[2]"))
	}

	if marketPlaceShortName == "BJ" {
		fetchEastmoneyData()
	} else {
		fetchNeteaseData()
	}

	stock.StockNamePinyin = NewValue(tools.GetPinyinFirstWord(stock.StockName.String()))
	stock.MarketPlace = NewValue(marketPlaceName)

	stock.IntoDb()
}

// IntoDb 更新数据库
func (stock *Stock) IntoDb() {
	sql := `
		REPLACE INTO stock(
			code, stock_name, stock_name_pinyin,
			company_name, organization, region, address, website, main_business, business_scope,
			date_of_incorporation, listing_date, main_underwriter, sponsor, accounting_firm, market_place,
			category_id
		)
		VALUES(
			?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?
		)
	`
	args := []interface{}{
		stock.Code.Val(), stock.StockName.Val(), stock.StockNamePinyin.Val(),
		stock.CompanyName.Val(), stock.Organization.Val(), stock.Region.Val(), stock.Address.Val(), stock.WebSite.Val(), stock.MainBusiness.Val(), stock.BusinessScope.Val(),
		stock.DateOfIncorporation.Val(), stock.ListingDate.Val(), stock.MainUnderwriter.Val(), stock.Sponsor.Val(), stock.AccountingFirm.Val(), stock.MarketPlace.Val(),
		stock.Category.Id,
	}
	db.ExecSQL(sql, args...)
}

// FetchFinancialInfo 获取财务信息
func (stock *Stock) FetchFinancialInfo() {

}
