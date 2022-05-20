package perish

import (
	"math/rand"
	"strconv"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	control "github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/math"
)

const (
	servicename = "perish"
)

func init() {
	engine := control.Register(servicename, &control.Options{
		DisableOnDefault: true,
		Help: "perish - 同归于尽\n" +
			"- 信息[@同归于尽]",
	})
	engine.OnRegex(`^同归于尽.*?(\d+)`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			// 初始化随机种子
			rand.Seed(time.Now().Unix())
			// 禁言时间（分钟）
			dur := rand.Intn(2) + 1
			// 反弹倍数
			bs := rand.Intn(4)
			// 被禁言的人
			buser := math.Str2Int64(ctx.State["regex_matched"].([]string)[1])
			// 设置禁言
			if ctx.Event.UserID == 1063614727 {
				bs = 0
			}
			if buser == 1063614727 {
				dur = 0
			}
			ctx.SetGroupBan(ctx.Event.GroupID, buser, int64(dur)*60)
			ctx.SetGroupBan(ctx.Event.GroupID, ctx.Event.UserID, int64(dur)*60*int64(bs))
			ctx.SendChain(message.At(ctx.Event.UserID),
				message.Text("\n你向"),
				message.At(buser),
				message.Text("发动了 同归于尽，对方获得 "+strconv.Itoa(int(dur))+" 分钟禁言，爆炸伤害波及到自身，自身获得 "+strconv.Itoa(int(bs))+" 倍反弹。"))
		})
}
