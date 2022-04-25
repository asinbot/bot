package rcon

import (
	"fmt"
	"strings"
)

// Map represents a playable landscape in HLL, having parts of the name constructed into a type.
type Map struct {
	Location string
	Type     string
	Side     string

	MapName
}

// MapName represents a server readable name for a map in HLL.
type MapName string

const (
	MapUtahBeachOffensiveUS          MapName = "utahbeach_offensive_us"
	MapFoyWarfare                    MapName = "foy_warfare"
	MapStMarieDuMontOffensiveUS      MapName = "stmariedumont_off_us"
	MapKurskOffensiveRussia          MapName = "kursk_offensive_rus"
	MapStMereEgliseWarfare           MapName = "stmereeglise_warfare"
	MapCarentanOffensiveUS           MapName = "carentan_offensive_us"
	MapHill400Warfare                MapName = "hill400_warfare"
	MapStMarieDuMontOffensiveGermany MapName = "stmariedumont_off_ger"
	MapHurtgenForestOffensiveGermany MapName = "hurtgenforest_offensive_ger"
	MapStalingradWarfare             MapName = "stalingrad_warfare"
	MapFoyOffensiveGermany           MapName = "foy_offensive_ger"
	MapUtahBeachOffensiveGermany     MapName = "utahbeach_offensive_ger"
	MapCarentanWarfare               MapName = "carentan_warfare"
	MapKurskWarfare                  MapName = "kursk_warfare"
	MapPurpleHeartLaneOffensiveUS    MapName = "purpleheartlane_offensive_us"
	MapStMereEgliseOffensiveUS       MapName = "stmereeglise_offensive_us"
	MapUtahBeachWarfare              MapName = "utahbeach_warfare"
	MapStalingradOffensiveGermany    MapName = "stalingrad_offensive_ger"
	MapHurtgenForestWarfare          MapName = "hurtgenforest_warfare_V2"
	MapStMereEgliseOffensiveGermany  MapName = "stmereeglise_offensive_ger"
	MapHill400OffensiveUS            MapName = "hill400_offensive_US"
	MapOmahaBeachOffensiveUS         MapName = "omahabeach_offensive_us"
	MapPurpleHeartLaneWarfare        MapName = "purpleheartlane_warfare"
	MapKurskOffensiveGermany         MapName = "kursk_offensive_ger"
	MapStalingradOffensiveRussia     MapName = "stalingrad_offensive_rus"
	MapStMarieDuMontWarfare          MapName = "stmariedumont_warfare"
	MapHurtgenForestOffensiveUS      MapName = "hurtgenforest_offensive_US"
)

// Map returns the current map in rotation for a Conn.
func (c *Conn) Map() (Map, error) {
	result, err := c.send("get", "map")
	if err != nil {
		return Map{}, fmt.Errorf("failed to get current map: %v", err)
	}

	return mapFromString(result), nil
}

// Maps returns all possible maps that can be in rotation for a Conn.
func (c *Conn) Maps() ([]Map, error) {
	result, err := c.send("get", "mapsforrotation")
	if err != nil {
		return nil, fmt.Errorf("failed to get maps for rotation: %v", err)
	}

	maps := []Map{}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse information for rotated maps")
	}

	for _, name := range args[1 : len(args)-1] {
		maps = append(maps, mapFromString(name))
	}

	return maps, nil
}

// Rotation returns the current map rotation for a Conn.
func (c *Conn) Rotation() ([]Map, error) {
	result, err := c.send("rotlist")
	if err != nil {
		return nil, fmt.Errorf("failed to get map rotation: %v", err)
	}

	maps := []Map{}

	for _, name := range strings.Split(result, "\n") {
		if name == "" {
			continue
		}

		maps = append(maps, mapFromString(name))
	}

	return maps, nil
}

// RotationAdd adds a map to the current rotation for a Conn.
func (c *Conn) RotationAdd(n MapName) error {
	_, err := c.send("rotadd", n.String())
	if err != nil {
		return fmt.Errorf("failed to add map to rotation: %v", err)
	}

	return nil
}

// RotationAdd adds a map to the current rotation for a Conn.
func (c *Conn) RotationRemove(n MapName) error {
	_, err := c.send("rotdel", n.String())
	if err != nil {
		return fmt.Errorf("failed to add map to rotation: %v", err)
	}

	return nil
}

// SetMap will change the current map in rotation for a Conn.
func (c *Conn) SetMap(n MapName) error {
	_, err := c.send("map", n.String())
	if err != nil {
		return fmt.Errorf("failed to set map as %s: %v", n, err)
	}

	return nil
}

// String returns a prettier standard string for a Map.
func (m Map) String() string {
	s := fmt.Sprintf("%s - %s", m.Location, m.Type)
	if m.Side != "" {
		s = fmt.Sprintf("%s (%s)", s, m.Side)
	}
	return s
}

// String converts a MapName into a standard string.
func (n MapName) String() string {
	return string(n)
}

func mapFromString(s string) Map {
	return mapFromName(MapName(s))
}

func mapFromName(n MapName) Map {
	name := n.String()

	args := strings.Split(name, "_")
	if len(args) == 0 {
		return Map{}
	}

	m := Map{
		Location: location(args[0]),
		MapName:  MapName(name),
	}

	switch len(args) {
	case 2:
		// m.Type = "Warfare"
		m.Type = "冲突"
	case 3:
		// m.Type = "Offensive"
		m.Type = "进攻"
		m.Side = side(args[2])
	}

	if name == "hurtgenforest_warfare_V2" {
		m.Type = "冲突"
	}

	return m
}

func location(name string) string {
	switch name {
	case "utahbeach":
		// return "Utah Beach"
		return "犹他海滩"
	case "omahabeach":
		// return "Omaha Beach"
		return "奥马哈海滩"
	case "foy":
		// return "Foy"
		return "佛依"
	case "stmariedumont":
		// return "St. Marie Du Mont"
		return "圣玛丽杜蒙特"
	case "kursk":
		// return "Kursk"
		return "库尔斯克"
	case "stmereeglise":
		// return "St. Mere Eglise"
		return "圣梅尔埃格利斯"
	case "carentan":
		// return "Carentan"
		return "卡朗唐"
	case "hill400":
		// return "Hill 400"
		return "400号高地"
	case "hurtgenforest":
		// return "Hurtgen Forest"
		return "许特根森林"
	case "stalingrad":
		// return "Stalingrad"
		return "斯大林格勒"
	case "purpleheartlane":
		// return "Purple Heart Lane"
		return "紫心小道"
	default:
		return strings.Title(name)
	}
}

func side(name string) string {
	switch name {
	case "us":
		// return "United States"
		return "美军进攻"
	case "ger":
		// return "Germany"
		return "德军进攻"
	case "rus":
		// return "Russia"
		return "苏联进攻"
	default:
		return ""
	}
}
