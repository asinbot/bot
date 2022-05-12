package hll

import (
	"fmt"
	"math"
	"strconv"

	"github.com/FloatTech/zbputils/img/text"
	"github.com/asinbot/bot/plugin/hll/rcon"
	"github.com/fogleman/gg"
	log "github.com/sirupsen/logrus"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

type Params struct {
	url string
	ctx *zero.Ctx
}

type groupServer struct {
	groupid int
	ids     []int
}

var gs = map[int64][]int{
	// jump 群
	426108037: {14584900},
	865821038: {14584900},
	// jump 新训练营
	798429948: {14584900},
	// icf 群
	741874107: {12246570, 14585453, 14581814},
	686811081: {12246570, 14585453, 14581814},
	// GSG
	909936359: {14251919, 14713002},
	// Dsquad
	758464100: {14544600},
	875570590: {14544600},
}

func getServerStatus(ctx *zero.Ctx) {
	data := request(Params{"https://api.battlemetrics.com/servers?page%5Bsize%5D=100&sort=&fields%5Bserver%5D=rank%2Cname%2Cplayers%2CmaxPlayers%2Caddress%2Cip%2Cport%2Ccountry%2Clocation%2Cdetails%2Cstatus&relations%5Bserver%5D=game%2CserverGroup&filter%5Bgame%5D=hll&filter%5Bsearch%5D=%2245.135.231.198%22+%7C+CN-J.U.M.P+RealWar+%7CHongKong%7CYY%3A28915649%7CQQ%3A426108037", ctx})
	if data != nil {
		msg := "【服务器列表】\n"
		num := 0
		serverList := data["data"].Array()
		sids := []int{13506120}
		if gs[ctx.Event.GroupID] != nil {
			sids = gs[ctx.Event.GroupID]
		}
		for _, server := range serverList {
			// if server.Get("id").Int() == 12319216 || server.Get("id").Int() == 13282192 {
			id16, _ := strconv.Atoi(server.Get("id").String())
			if inArray(id16, sids) {
				if num != 0 {
					msg += "\n=========\n"
				}
				info := server.Get("attributes").Map()
				msg += "服务器名：" + info["name"].String() + "\n"
				// msg += "状态：" + info["status"].String() + "\n"
				msg += "状态：" + getServerState(info["status"].String()) + "\n"
				msg += "地图：" + info["details"].Get("map").String() + "\n"
				msg += "人数：" + info["players"].String() + "/" + info["maxPlayers"].String() + "\n"
				msg += "服务器地址：" + getServerCountry(info["country"].String()) + "\n"
				msg += "是否需要密码：" + getServerNeedPwd(info["details"].Get("password").String())
				num++
			}
		}
		// ctx.SendChain(message.Text(msg))
		data, err := text.RenderToBase64(msg, text.FontFile, 400, 14)
		if err != nil {
			log.Errorln("[hll]:", err)
		}
		ctx.SendChain(message.Image("base64://" + helper.BytesToString(data)))
	}
}

func getSteamId64(ctx *zero.Ctx) {
	data := request(Params{"http://steamrep.com/search?q=", ctx})
	if data != nil {

	}
}

func getOnlinePlayers(ctx *zero.Ctx) {
	data2 := request(Params{"https://api.battlemetrics.com/servers?page%5Bsize%5D=100&sort=&fields%5Bserver%5D=rank%2Cname%2Cplayers%2CmaxPlayers%2Caddress%2Cip%2Cport%2Ccountry%2Clocation%2Cdetails%2Cstatus&relations%5Bserver%5D=game%2CserverGroup&filter%5Bgame%5D=hll&filter%5Bsearch%5D=%2245.135.231.198%22+%7C+CN-J.U.M.P+RealWar+%7CHongKong%7CYY%3A28915649%7CQQ%3A426108037", ctx})
	if data2 != nil {
		playerNum := 0
		serverList := data2["data"].Array()
		sids := []int{13506120}
		if gs[ctx.Event.GroupID] != nil {
			sids = gs[ctx.Event.GroupID]
		}
		data := request(Params{"https://api.battlemetrics.com/players?version=%5E0.1.0&page%5Bsize%5D=100&sort=-lastSeen&fields%5Bidentifier%5D=type%2Cidentifier%2ClastSeen&fields%5Bserver%5D=name&filter%5Bpublic%5D=true&filter%5Bservers%5D=" + strconv.Itoa(sids[0]) + "&filter%5Bonline%5D=true&filter%5BplayerFlags%5D=", ctx})
		if data != nil {
			for _, server := range serverList {
				// if server.Get("id").Int() == 12319216 || server.Get("id").Int() == 13282192 {
				id16, _ := strconv.Atoi(server.Get("id").String())
				if inArray(id16, sids) {
					playerNum, _ = strconv.Atoi(server.Get("attributes.players").String())
				}
			}

			playerList := data["data"].Array()

			n := 2
			if n < playerNum {
				n = playerNum
			}
			if n > 5 {
				n = 5
			}

			dc := gg.NewContext(150*n+10, int(math.Ceil(float64(playerNum)/5))*15+25)
			dc.SetRGB(1, 1, 1)
			dc.Clear()
			dc.SetRGB(0, 0, 0)
			err := dc.LoadFontFace(text.FontFile, float64(12))
			if err != nil {
				log.Errorln("[txt2img]", err)
				return
			}
			dc.DrawString("当前服务器在线玩家为（"+strconv.Itoa(playerNum)+"人）：", 5, 15)
			num := 0
			var x, y float64
			for _, server := range playerList {
				if playerNum <= 0 {
					break
				}
				dc.DrawString(server.Get("attributes.name").String(), x*150+5, (y+1)*15+20)
				x++
				if x > 4 {
					x = 0
					y++
				}
				playerNum--
				num++
			}
			data, err := canvas2Base64(dc)
			if err != nil {
				log.Errorln("[hll]:", err)
			}
			ctx.SendChain(message.Image("base64://" + helper.BytesToString(data)))
		}
	}
}

func getOnlinePlayers2(ctx *zero.Ctx) {
	conn := rcon.Get("103.161.224.92:28035")
	p, _ := conn.Players()
	playerNum := len(p)

	n := math.Min(math.Max(float64(playerNum), 2), 5)

	dc := gg.NewContext(150*int(n)+10, int(math.Ceil(float64(playerNum)/5))*15+25)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	err := dc.LoadFontFace(text.FontFile, float64(12))
	if err != nil {
		log.Errorln("[txt2img]", err)
		return
	}
	dc.DrawString("当前服务器在线玩家为（"+strconv.Itoa(playerNum)+"人）：", 5, 15)
	// num := 0
	var x, y float64
	for _, player := range p {
		// conn.Player(player.Name)
		dc.DrawString(player.Name, x*150+5, (y+1)*15+20)
		// dc.DrawString(player.ID64, 1*150+5, (y+1)*15+20)
		x++
		if x > 4 {
			x = 0
			y++
		}
	}
	data, err := canvas2Base64(dc)
	if err != nil {
		log.Errorln("[hll]:", err)
	}
	ctx.SendChain(message.Image("base64://" + helper.BytesToString(data)))
}

func getMapRotation(ctx *zero.Ctx) {
	conn := rcon.Get("103.161.224.92:28035")
	maps, _ := conn.Rotation()
	cmap, _ := conn.Map()
	mapsNum := len(maps)

	n := 2

	dc := gg.NewContext(150*int(n)+10, int(math.Ceil(float64(mapsNum)))*15+25)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	err := dc.LoadFontFace(text.FontFile, float64(12))
	if err != nil {
		log.Errorln("[txt2img]", err)
		return
	}
	dc.DrawString("服务器地图列表为：", 5, 15)
	// num := 0
	var x, y float64
	for _, m := range maps {
		// conn.Player(player.Name)
		mapName := m.Location + " - " + m.Type
		if m.Side != "" {
			mapName += " - " + m.Side
		}
		if m == cmap {
			mapName += "（当前地图）"
		}
		dc.DrawString(string(mapName), x*150+5, (y+1)*15+20)
		y++
	}
	data, err := canvas2Base64(dc)
	if err != nil {
		log.Errorln("[hll]:", err)
	}
	ctx.SendChain(message.Image("base64://" + helper.BytesToString(data)))
}

func getPlayerInfo(ctx *zero.Ctx, playerName string) {
	conn := rcon.Get("103.161.224.92:28035")
	player, err := conn.Playerinfo(playerName)

	fmt.Print(player)

	if err != nil {
		ctx.SendChain(message.Text("查询该玩家失败或此玩家不在服务器中"))
		return
	}

	msg := ""
	msg += "昵称：" + player.Name + "\n"
	msg += "等级：" + player.Level + "\n"
	// msg += "状态：" + info["status"].String() + "\n"
	msg += "阵营：" + player.Team + "\n"
	msg += "小队：" + player.Unit + "\n"
	msg += "职业：" + player.Role + "\n"
	msg += "KD：" + player.KD.Kills + "/" + player.KD.Deaths + "\n"
	msg += "----------" + "\n"
	msg += "击杀分：" + player.Score.C + " / 进攻分：" + player.Score.O + "\n"
	msg += "防御分：" + player.Score.D + " / 支援分：" + player.Score.S
	data, err := text.RenderToBase64(msg, text.FontFile, 250, 14)
	if err != nil {
		log.Errorln("[hll]:", err)
	}
	ctx.SendChain(message.Image("base64://" + helper.BytesToString(data)))
}
