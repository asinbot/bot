package rcon

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	Name string
	ID64 string
}

type PlayerKD struct {
	Kills  string
	Deaths string
}

type PlayerScore struct {
	// 击杀分
	C string
	// 进攻分
	O string
	// 防御分
	D string
	// 支援分
	S string
}

type PlayerInfo struct {
	Name    string
	ID64    string
	Team    string
	Role    string
	Unit    string
	Loadout string
	KD      PlayerKD
	Score   PlayerScore
	Level   string
}

type Ban struct {
	Player
	Admin
	Reason string
	time.Duration
	time.Time
}

func (c *Conn) BannedTemporarily() ([]Ban, error) {
	result, err := c.send("get", "tempbans")
	if err != nil {
		return nil, fmt.Errorf("failed to get permanent bans: %v", err)
	}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse temporary ban information")
	}

	admins, err := c.Admins()
	if err != nil {
		return nil, fmt.Errorf("failed to parse temporary ban admins: %v", err)
	}

	bans := []Ban{}

	for _, ban := range args[1 : len(args)-1] {
		b, err := parseBanned(ban, admins)
		if err != nil {
			return nil, err
		}

		bans = append(bans, b)
	}

	return bans, nil
}

func (b Ban) String() string {
	return fmt.Sprintf("%s [%s] from: %s until %s by: %s", b.Player.String(), b.Reason, b.Time.Format(time.Stamp), time.Now().Add(b.Duration).Format(time.Stamp), b.Admin)
}

func (c *Conn) BannedPermanently() ([]Ban, error) {
	result, err := c.send("get", "permabans")
	if err != nil {
		return nil, fmt.Errorf("failed to get permanent bans: %v", err)
	}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse temporary ban information")
	}

	admins, err := c.Admins()
	if err != nil {
		return nil, fmt.Errorf("failed to parse temporary ban admins: %v", err)
	}

	bans := []Ban{}

	for _, ban := range args[1 : len(args)-1] {
		b, err := parseBanned(ban, admins)
		if err != nil {
			return nil, err
		}

		bans = append(bans, b)
	}

	return bans, nil
}

// BanPermanently will remove an active player and block server access indefinitely.
func (c *Conn) BanPermanently(p Player, reason, admin string) error {
	_, err := c.send("permaban", q(p.ID64), q(reason), q(admin))
	if err != nil {
		return fmt.Errorf("failed to set permanent ban %s: %v", p, err)
	}

	return nil
}

// BanRemove will remove a Player's temp or perma ban and re-allow server access.
func (c *Conn) BanRemove(p Player) error {
	_, err := c.send("pardontempban", q(p.ID64))
	if err != nil {
		if err != ErrResultFailed {
			return fmt.Errorf("failed to remove ban for %s: %v", p, err)
		}

		_, err := c.send("pardonpermaban", q(p.ID64))
		if err != nil {
			return fmt.Errorf("failed to remove ban for %s: %v", p, err)
		}
	}

	return nil
}

// BanTemporarily will remove an active player and block server access for the specified hours.
func (c *Conn) BanTemporarily(p Player, hours int, reason, admin string) error {
	_, err := c.send("tempban", q(p.ID64), strconv.Itoa(hours), q(reason), q(admin))
	if err != nil {
		return fmt.Errorf("failed to set temporary ban %s: %v", p, err)
	}

	return nil
}

// Kick will remove an active player.
func (c *Conn) Kick(p Player, reason string) error {
	_, err := c.send("kick", q(p.Name), q(reason))
	if err != nil {
		return fmt.Errorf("failed to kick %s: %v", p, err)
	}

	return nil
}

// Punish will punish an active player.
func (c *Conn) Punish(p Player, reason string) error {
	_, err := c.send("punish", q(p.Name), q(reason))
	if err != nil {
		return fmt.Errorf("failed to punish %s: %v", p, err)
	}

	return nil
}

// Player return a Player for a given username.
func (c *Conn) Player(username string) (Player, error) {
	result, err := c.send("playerinfo", username)
	if err != nil {
		return Player{}, fmt.Errorf("failed to get player information for %s: %v", username, err)
	}

	args := strings.Split(result, "\n")
	if len(args) < 2 {
		return Player{}, fmt.Errorf("invalid player information for %s", username)
	}

	fmt.Print(result)

	name := strings.Split(args[0], ":")
	id := strings.Split(args[1], ":")

	if len(name) < 2 || len(id) < 2 {
		return Player{}, fmt.Errorf("invalid player information for %s", username)
	}

	return Player{
		Name: strings.TrimSpace(name[1]),
		ID64: strings.TrimSpace(id[1]),
	}, nil
}

func (c *Conn) Playerinfo(username string) (PlayerInfo, error) {
	result, err := c.send("playerinfo", username)
	if err != nil {
		return PlayerInfo{}, fmt.Errorf("failed to get player information for %s: %v", username, err)
	}

	args := strings.Split(result, "\n")
	if len(args) < 2 {
		return PlayerInfo{}, fmt.Errorf("invalid player information for %s", username)
	}

	fmt.Print(result)

	nameIdx, idIdx, teamIdx, roleIdx, unitIdx, loadoutIdx, killsIdx, scoreIdx, levelIdx := 0, 1, 2, 3, 4, 5, 6, 7, 8

	name := strings.Split(args[nameIdx], ":")
	id := strings.Split(args[idIdx], ":")

	team, err := parseTeam(args[teamIdx])
	if err != nil {
		roleIdx -= 1
		unitIdx -= 1
		loadoutIdx -= 1
		killsIdx -= 1
		scoreIdx -= 1
		levelIdx -= 1
	}

	role, err := parseRole(args[roleIdx])
	if err != nil {
		unitIdx -= 1
		loadoutIdx -= 1
		killsIdx -= 1
		scoreIdx -= 1
		levelIdx -= 1
	}

	unit, err := parseUnit(args[unitIdx])
	if err != nil {
		loadoutIdx -= 1
		killsIdx -= 1
		scoreIdx -= 1
		levelIdx -= 1
	}

	loadout, err := parseLoadout(args[loadoutIdx])
	if err != nil {
		killsIdx -= 1
		scoreIdx -= 1
		levelIdx -= 1
	}

	KD, err := parseKD(args[killsIdx])
	if err != nil {
		scoreIdx -= 1
		levelIdx -= 1
	}

	score, err := parseScore(args[scoreIdx])
	if err != nil {
		levelIdx -= 1
	}

	level, _ := parseLevel(args[levelIdx])

	if len(name) < 2 || len(id) < 2 {
		return PlayerInfo{}, fmt.Errorf("invalid player information for %s", username)
	}

	return PlayerInfo{
		Name:    strings.TrimSpace(name[1]),
		ID64:    strings.TrimSpace(id[1]),
		Team:    team,
		Role:    role,
		Unit:    unit,
		Loadout: loadout,
		KD:      KD,
		Score:   score,
		Level:   level,
	}, nil
}

// Players returns all active Players.
func (c *Conn) Players() ([]Player, error) {
	result, err := c.send("get", "playerids")
	if err != nil {
		return nil, fmt.Errorf("failed to get information for all players: %v", err)
	}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse information for all players")
	}

	players := []Player{}

	for _, p := range args[1 : len(args)-1] {
		nameID := strings.Split(p, ":")

		players = append(players, Player{
			Name: strings.TrimSpace(nameID[0]),
			ID64: strings.TrimSpace(nameID[1]),
		})
	}

	return players, nil
}

func (c *Conn) PlayerInfos() {
	result, _ := c.send("get", "players")

	fmt.Print(result)
}

func (c *Conn) SetSwitchTeamNow(p Player) error {
	_, err := c.send("switchteamnow", p.ID64)
	if err != nil {
		return fmt.Errorf("failed to set switch player now for %s: %v", p.String(), err)
	}

	return nil
}

func (c *Conn) SetSwitchTeamOnDeath(p Player) error {
	_, err := c.send("switchteamondeath", p.ID64)
	if err != nil {
		return fmt.Errorf("failed to set switch player on death for %s: %v", p.String(), err)
	}

	return nil
}

func (p Player) String() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.ID64)
}

var matchBannedName = regexp.MustCompile(`: nickname "(.*?)" banned"`)
var matchBanned = regexp.MustCompile(`(.*?) : nickname "(.*?)" banned for (.*?) hours on (.*?) for "(.*?)" by admin "(.*?)"`)

// parseBanned will return a slice of items that are in quotes from a string.
func parseBanned(s string, admins []Admin) (Ban, error) {
	b := Ban{}

	matches := matchBanned.FindAllStringSubmatch(s, -1)
	if len(matches) != 1 {
		return b, fmt.Errorf("failed to parse temporary ban information")
	}

	match := matches[0]

	b.ID64 = match[1]
	b.Name = match[2]

	hours, err := strconv.Atoi(match[3])
	if err != nil {
		return Ban{}, fmt.Errorf("failed to parse temporary ban hours: %v", err)
	}

	b.Duration = time.Duration(hours) * time.Hour

	b.Time, err = time.Parse("2006.01.02-15.04.05", match[4])
	if err != nil {
		return Ban{}, fmt.Errorf("failed to parse temporary ban time: %v", err)
	}

	b.Reason = match[5]

	b.Admin = unknownAdmin

	for i := range admins {
		if admins[i].Name == match[6] || admins[i].ID64 == match[6] {
			b.Admin = admins[i]
			break
		}
	}

	return b, nil
}

// 解析玩家团队
func parseTeam(s string) (string, error) {
	team := strings.Split(s, ":")
	if team[0] != "Team" {
		return "", fmt.Errorf("No Team")
	} else {
		return tfTeam(strings.TrimSpace(team[1])), nil
	}
}

// 解析玩家职业
func parseRole(s string) (string, error) {
	role := strings.Split(s, ":")
	if role[0] != "Role" {
		return "", fmt.Errorf("No Role")
	} else {
		return tfRole(strings.TrimSpace(role[1])), nil
	}
}

// 解析玩家小队
func parseUnit(s string) (string, error) {
	unit := strings.Split(s, ":")
	if unit[0] != "Unit" {
		return "", fmt.Errorf("No Unit")
	} else {
		u := strings.Split(strings.TrimSpace(unit[1]), " - ")
		u2 := strings.Split(strings.TrimSpace(u[1]), "")
		return strings.TrimSpace(u2[0]), nil
	}
}

//
func parseLoadout(s string) (string, error) {
	loadout := strings.Split(s, ":")
	if loadout[0] != "Loadout" {
		return "", fmt.Errorf("No Loadout")
	} else {
		return strings.TrimSpace(loadout[1]), nil
	}
}

// 解析玩家击杀死亡情况
func parseKD(s string) (PlayerKD, error) {
	kd := strings.Split(s, "-")
	k := strings.Split(kd[0], ":")
	if k[0] != "Kills" {
		return PlayerKD{}, fmt.Errorf("No Kills")
	} else {
		return PlayerKD{
			Kills:  strings.TrimSpace(k[1]),
			Deaths: strings.TrimSpace(strings.Split(kd[1], ":")[1]),
		}, nil
	}
}

// 解析玩家得分情况
func parseScore(s string) (PlayerScore, error) {
	score := strings.Split(s, ":")
	if score[0] != "Score" {
		return PlayerScore{}, fmt.Errorf("No Score")
	} else {
		s := strings.Split(strings.TrimSpace(score[1]), ", ")
		return PlayerScore{
			C: strings.TrimSpace(strings.Split(s[0], " ")[1]),
			O: strings.TrimSpace(strings.Split(s[1], " ")[1]),
			D: strings.TrimSpace(strings.Split(s[2], " ")[1]),
			S: strings.TrimSpace(strings.Split(s[3], " ")[1]),
		}, nil
	}
}

// 解析玩家等级
func parseLevel(s string) (string, error) {
	level := strings.Split(s, ":")
	if level[0] != "Level" {
		return "", fmt.Errorf("No Level")
	} else {
		return strings.TrimSpace(level[1]), nil
	}
}

// 翻译团队名
func tfTeam(team string) string {
	switch team {
	case "Allies":
		return "同盟国"
	case "Axis":
		return "轴心国"
	default:
		return team
	}
}

// 翻译职业名
func tfRole(role string) string {
	switch role {
	case "Commander":
		return "指挥官"
	case "Officer":
		return "军官"
	case "Rifleman":
		return "步枪兵"
	case "Assault":
		return "突击兵"
	case "AutomaticRifleman":
		return "自动步枪兵"
	case "Medic":
		return "医疗兵"
	case "Support":
		return "支援兵"
	case "MachineGunner":
		return "机枪兵"
	case "AntiTank":
		return "反坦克兵"
	case "Engineer":
		return "工兵"
	case "TankCommander":
		return "车组指挥官"
	case "Crewman":
		return "车组成员"
	case "Spotter":
		return "观察手"
	case "Sniper":
		return "狙击手"
	default:
		return role
	}
}
