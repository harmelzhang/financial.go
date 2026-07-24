package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"harmel.cn/financial/internal/cmd"
	"harmel.cn/financial/internal/public"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	ctx := gctx.New()
	gCfg := g.Cfg()

	// 指数样本类型
	for code, name := range gCfg.MustGet(ctx, "indexSample").Map() {
		public.IndexSampleType[code] = fmt.Sprint(name)
	}
	// 市场标识前缀
	for key, value := range gCfg.MustGet(ctx, "marketPrefix").Map() {
		anyValues := value.([]any)
		values := make([]string, 0, len(anyValues))
		for _, v := range anyValues {
			values = append(values, fmt.Sprint(v))
		}
		switch key {
		case "shanghai":
			public.ShanghaiMarketPrefixs = values
		case "shenzhen":
			public.ShenzhenMarketPrefixs = values
		case "beijing":
			public.BeijingMarketPrefixs = values
		}
	}

	cmd, err := gcmd.NewFromObject(cmd.CommandMain{})
	if err != nil {
		panic(err)
	}
	cmd.Run(ctx)
}
