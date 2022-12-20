package main

import (
	"financial/config"
	categoryConfig "financial/config/category"
	_ "financial/config/indexSample"
	"financial/models"
	"financial/utils/tools"
	"log"
	"time"
)

func main() {
	progress := models.NewProgress()

	// 初始化配置文件
	if tools.FileIsExist(config.ProgressFileName) {
		progress.Load(config.ProgressFileName)
	} else {
		progress.Write(config.ProgressFileName)
	}

	// 如果到了五月一日，全部重跑（年报全部出了）
	if time.Now().Format("01-02") == "05-01" {
		progress = models.NewProgress()
	}

	// 如果上次成功了，判断时间是否大于配置天数
	if progress.Done {
		if time.Now().Unix()-progress.Time >= config.TaskIntervalDay*24*3600 {
			progress = models.NewProgress()
		} else {
			log.Printf("任务结束 : 离上次任务成功结束时间小于%d天", config.TaskIntervalDay)
			return
		}
	}

	for _, category := range categoryConfig.GetCategorys() {
		if category.Exist() {
			category.UpdateDb()
		} else {
			category.IntoDb()
		}

		stocks := category.GetStocks()
		for _, stock := range stocks {
			// 跳过B股和场内基金
			if tools.IndexOf(config.ExcludeStockCodePrefix, stock.Code.String()[0:1]) != -1 {
				continue
			}

			if tools.IndexOf(progress.Info[category.Id], stock.Code.String()) != -1 {
				continue
			}

			stock.BuildStockInfo()
			stock.ReplaceDb()
			progress.Info[category.Id] = append(progress.Info[category.Id], stock.Code.String())
			progress.Time = time.Now().Unix()
			progress.Write(config.ProgressFileName)
		}
	}

	progress.Done = true
	progress.Time = time.Now().Unix()
	progress.Write(config.ProgressFileName)
	log.Println("任务结束！")
}
