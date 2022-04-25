package monitor

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	control "github.com/FloatTech/zbputils/control"
)

const (
	servicename = "monitor"
)

func init() {
	engine := control.Register(servicename, &control.Options{
		DisableOnDefault: false,
		Help:             "monitor - 聊天监控\n",
	})
	engine.OnMessage().SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.UserID == 3480326047 {
				ctx.SendPrivateMessage(1063614727, message.Text("[消息监听] 许晓柔在群里说："+ctx.MessageString()))
			}
			// if ctx.Event.UserID == 3439075097 {
			// 	ctx.SendPrivateMessage(1063614727, message.Text("[消息监听] 崔恩熙在群里说："+ctx.MessageString()))
			// }
		})
}
