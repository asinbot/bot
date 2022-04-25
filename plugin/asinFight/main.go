package asinFight

import (
	control "github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
)

const (
	servicename = "asinFight"
	pori        = 2
	actionname  = "【刺客大乱斗】"
)

func init() {
	engine := control.Register(servicename, &control.Options{
		DisableOnDefault: true,
	})
	// 监听心跳
	engine.OnMetaEvent().SetBlock(true).
		Handle(func(ctx *zero.Ctx) {

		})
}
