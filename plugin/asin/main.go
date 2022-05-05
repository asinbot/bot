package asin

import (
	zero "github.com/wdvxdr1123/ZeroBot"

	control "github.com/FloatTech/zbputils/control"
)

const (
	servicename = "asin"
)

func init() {
	engine := control.Register(servicename, &control.Options{
		DisableOnDefault: true,
		Help: "asin - 刺客系统\n" +
			"- 信息 [@xxx]\n" +
			"- 属性 [@xxx]",
	})
	engine.OnFullMatch("信息").SetBlock(true).
		Handle(getUserInfoWithScore)
	engine.OnFullMatch("属性").SetBlock(true).
		Handle(getUserAttrWithFight)
	engine.OnFullMatch("属性加点").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {

		})
}
