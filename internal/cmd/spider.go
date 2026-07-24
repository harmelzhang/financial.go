package cmd

import (
	"context"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"harmel.cn/financial/internal/public"
	"harmel.cn/financial/internal/spider"
)

type CommandMainSpiderInput struct {
	g.Meta `name:"spider" brief:"start spider server"`
	// 爬虫抓取数据模式，默认是差异抓取（只抓取新的年份），all 表示全部重新抓取
	Mode string `name:"mode" brief:"spider fetch mode" default:"diff"`
}

type CommandMainSpiderOutput struct{}

func (c *CommandMain) Spider(ctx context.Context, in CommandMainSpiderInput) (out *CommandMainSpiderOutput, err error) {
	gCfg := g.Cfg()
	spiderLogger := g.Log("spider")

	rootDir, err := os.Getwd()
	if err != nil {
		spiderLogger.Errorf(ctx, "get root dir failed, err is %v", err)
		return
	}

	// 系统变量
	public.SpiderTaskIntervalDays = gCfg.MustGet(ctx, "spider.taskIntervalDays").Int64()
	public.SpiderExecutorPoolSize = gCfg.MustGet(ctx, "spider.executorPoolSize").Int()
	public.SpiderTimtout = gCfg.MustGet(ctx, "spider.timeout").Int()
	public.QueryReportPageSize = gCfg.MustGet(ctx, "spider.queryReportPageSize").Int()

	// 启动爬虫管理器
	spiderManager := spider.NewSpiderManager(rootDir)
	spiderLogger.Debugf(ctx, "start spider mode is %s", in.Mode)
	err = spiderManager.Start(ctx, in.Mode)

	return
}
