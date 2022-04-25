package hll

import (
	"fmt"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	control "github.com/FloatTech/zbputils/control"
)

const (
	servicename = "hll"
)

type Test struct {
	lastUpdateTime string
}

func init() {
	// wsc := NewWsClientManager("ws.battlemetrics.com", "/?audit_log=id%3Dd24c7c96-a521-4136-b468-abfc28849fcc", 10000)
	// wsc.start()
	engine := control.Register(servicename, &control.Options{
		DisableOnDefault: true,
		Help: "hll - 人间地狱\n" +
			"- 查看服务器状态 | 查询服务器状态 | 服务器状态\n" +
			"- 查看在线玩家 | 查询在线玩家 | 在线玩家\n",
	})
	engine.OnFullMatchGroup([]string{"查看服务器状态", "查询服务器状态", "服务器状态", "server status"}).SetBlock(true).
		Handle(getServerStatus)
	engine.OnFullMatchGroup([]string{"查看在线玩家", "查询在线玩家", "在线玩家"}).SetBlock(true).
		// Handle(func(ctx *zero.Ctx) {
		// 	if ctx.Event.GroupID == 426108037 || ctx.Event.GroupID == 865821038 || ctx.Event.GroupID == 798429948 {
		// 		getOnlinePlayers(ctx)
		// 	}
		// })
		// Handle(getOnlinePlayers)
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID == 875570590 {
				// getOnlinePlayers(ctx)
				getOnlinePlayers2(ctx)
			} else {
				getOnlinePlayers(ctx)
			}
		})
	engine.OnFullMatchGroup([]string{"查看地图列表", "查询地图列表", "地图列表"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID == 875570590 {
				// getOnlinePlayers(ctx)
				getMapRotation(ctx)
			}
		})
	engine.OnKeywordGroup([]string{"如何加入", "怎么加入", "加战队", "加入战队"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID == 426108037 || ctx.Event.GroupID == 865821038 || ctx.Event.GroupID == 798429948 {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("加入战队需要参加新训哦~详情请联系管理员。"))
			}
		})
	engine.OnKeywordGroup([]string{"求带", "萌新"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID == 426108037 || ctx.Event.GroupID == 865821038 {
				// ctx.SendChain(message.At(ctx.Event.UserID), message.Text("萌新想要一起玩可以来YY 28915649，看哪个子频道有人可以去问问。每天晚上也会有人组织带新的。新训营群号：798429948"))
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("萌新想要一起玩可以来YY 230863，看哪个子频道有人可以去问问。每天晚上也会有人组织带新的。新训营群号：798429948"))
			}
		})
	engine.OnKeywordGroup([]string{"vip", "VIP"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID == 426108037 || ctx.Event.GroupID == 865821038 || ctx.Event.GroupID == 798429948 {
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("想要获取服务器VIP的话，可以选择成为战队正式队员，或者经常暖服获取，再或者联系管理员购买。"))
			}
		})
	engine.OnRegex(`^查询玩家 (.*)$`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			msg := ctx.State["regex_matched"].([]string)[1]
			fmt.Print(msg)
			getPlayerInfo(ctx, msg)
			// rely, err := tl(msg[0])
			// if err != nil {
			// 	ctx.SendChain(message.Text("ERROR: ", err))
			// }
			// info := gjson.ParseBytes(rely)
			// repo := info.Get("data.0")
			// process.SleepAbout1sTo2s()
			// ctx.SendChain(message.Text(repo.Get("value.0")))
		})
	// engine.OnKeyword("测试").SetBlock(true).
	// 	Handle(func(ctx *zero.Ctx) {
	// 	})
}
