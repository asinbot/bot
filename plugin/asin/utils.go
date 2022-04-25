package asin

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 获取参数
func getArgs(ctx *zero.Ctx) []string {
	return strings.Split(ctx.State["args"].(string), " ")
}

// 随机回复文本
func randText(text ...string) message.MessageSegment {
	return message.Text(text[rand.Intn(len(text))])
}

// 获取用户（如果有at就跟着at走，没有就是发信人）
func getUser(ctx *zero.Ctx) int64 {
	if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
		qq, _ := strconv.ParseInt(ctx.Event.Message[1].Data["qq"], 10, 64)
		return qq
	}
	return ctx.Event.UserID
}

// 网络请求
func reuqest(p Params) map[string]gjson.Result {
	resp, err := http.PostForm("http://api.asin.pohun.cn/api/"+p.action, p.params)
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			resp.Body.Close()
			json := gjson.ParseBytes(body)
			if json.Get("errCode").Int() != 200 {
				p.ctx.SendChain(message.Text(json.Get("errMsg").String()))
				return nil
			}
			return json.Get("data").Map()
		}
		if p.ctx != nil {
			p.ctx.SendChain(message.Text("读取数据失败"))
		}
		return nil
	}
	if p.ctx != nil {
		p.ctx.SendChain(message.Text("读取数据失败"))
	}
	return nil
}
