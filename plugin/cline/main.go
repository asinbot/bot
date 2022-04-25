package cline

import (
	"io/ioutil"
	"net/http"
	"strconv"

	control "github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	servicename = "cline"
)

func init() {
	// engine := control.Register(servicename, &control.Options
	engine := control.Register(servicename, &control.Options{
		DisableOnDefault: true,
		Help:             "cline - 动画排行",
	})
	engine.OnFullMatch("动画排行").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			resp, err := http.Get("http://cline.pohun.com/server/api.php")
			if err != nil {
				ctx.SendChain(message.Text("排行榜获取失败"))
				return
			}
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				ctx.SendChain(message.Text("排行榜获取失败"))
				return
			}
			json := gjson.ParseBytes(data)
			// fmt.Print(json, "\n")
			msg := "本季度动画播放量排行榜（仅显示前20条）"
			for idx, v := range json.Array() {
				msg += "\n" + strconv.Itoa(idx+1) + ". " + v.Get("title").String()
				if idx == 19 {
					break
				}
			}
			msg += "\n\n更多动画排行榜请访问：http://cline.pohun.com"
			ctx.SendChain(message.Text(msg))
		})
}
