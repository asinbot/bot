package asin

import (
	"net/url"
	"strconv"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type Params struct {
	action string
	params url.Values
	ctx    *zero.Ctx
}

func getUserInfoWithScore(ctx *zero.Ctx) {
	params := url.Values{}
	qq := getUser(ctx)
	params.Add("qq", strconv.FormatInt(qq, 10))
	data := reuqest(Params{
		action: "UserInfo/getUserInfoWithScore",
		params: params,
		ctx:    ctx,
	})
	if data == nil {
		return
	}
	var sex, score string
	if data["score"].Int() < 0 {
		score = "？？？"
	} else {
		score = data["score"].String()
	}
	if data["sex"].Int() == 0 {
		sex = "未知"
	} else if data["sex"].Int() == 1 {
		sex = "男"
	} else {
		sex = "女"
	}
	ctx.SendChain(message.At(getUser(ctx)), message.Text(
		"\n",
		"姓名："+data["nickname"].String()+"\n",
		"排名："+data["rank"].String()+"\n",
		"积分："+score+"\n",
		"暗币："+data["credit"].String()+"\n",
		"年龄："+data["age"].String()+"\n",
		"性别："+sex+"\n",
		"身高："+data["height"].String()+"\n",
		"体重："+data["weight"].String()+"\n",
		"介绍："+data["introduce"].String()+"\n",
		"加入组织时间："+data["ctime"].String(),
	))
}

func getUserAttrWithFight(ctx *zero.Ctx) {
	params := url.Values{}
	qq := getUser(ctx)
	params.Add("qq", strconv.FormatInt(qq, 10))
	data := reuqest(Params{
		action: "UserAttr/getUserAttrWithFight",
		params: params,
		ctx:    ctx,
	})
	if data == nil {
		return
	}
	ctx.SendChain(message.At(getUser(ctx)), message.Text(
		"\n",
		"姓名："+data["nickname"].String()+"\n",
		"力量："+data["str"].String()+"\n",
		"敏捷："+data["dex"].String()+"\n",
		"体质："+data["con"].String()+"\n",
		"智力："+data["ine"].String()+"\n",
		"感知："+data["wis"].String()+"\n",
		"魅力："+data["cha"].String()+"\n",
		"自由属性点："+data["free"].String()+"\n\n",
		"血量上限（大乱斗）："+data["maxBld"].String()+"\n",
		"攻击力（大乱斗）："+data["atk"].String()+"\n",
		"暴击率（大乱斗）："+data["crit"].String(),
	))
}
