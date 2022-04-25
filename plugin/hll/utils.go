package hll

import (
	"encoding/base64"
	"image/jpeg"
	"io/ioutil"
	"net/http"

	"github.com/FloatTech/zbputils/binary"
	"github.com/fogleman/gg"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// 网络请求
func request(p Params) map[string]gjson.Result {
	resp, err := http.Get(p.url)
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			resp.Body.Close()
			json := gjson.ParseBytes(body)
			return json.Map()
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

func inArray(id int, ids []int) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func getServerState(state string) string {
	if state == "online" {
		return "在线"
	}
	return "离线"
}

func getServerCountry(country string) string {
	switch country {
	case "HK":
		return "香港"
	case "SG":
		return "新加坡"
	case "JP":
		return "日本"
	default:
		return country
	}
}

func getServerNeedPwd(isPwd string) string {
	if isPwd == "false" {
		return "否"
	}
	return "是"
}

func canvas2Base64(canvas *gg.Context) (base64Bytes []byte, err error) {
	buffer := binary.SelectWriter()
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	var opt jpeg.Options
	opt.Quality = 70
	if err = jpeg.Encode(encoder, canvas.Image(), &opt); err != nil {
		return nil, err
	}
	encoder.Close()
	base64Bytes = buffer.Bytes()
	binary.PutWriter(buffer)
	return
}
