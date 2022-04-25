# rcon
Bindings for the Hell Let Loose RCON tool written purely in Go.

This project is constructed to act like a third-party library for other applications. These might include a REST server, a command line interface, etc.


`import "github.com/verocity-gaming/rcon"`

# Client

## Connect to a server
```
c, err := rcon.New(addr, password)
if err != nil {
	panic(err)
}
defer c.Close()
```

# Maps

## Get the current map
```
m, err := c.Map()
if err != nil {
	panic(err)
}

println(m.Location, m.Type, m.Side)
```

## Set the current map
```
err = c.SetMap(rcon.MapCarentanOffensiveUS)
if err != nil {
	panic(err)
}
```

## Get the current map rotation
```
r, err := c.Rotation()
if err != nil {
	panic(err)
}

for _, m := range r {
	println(m.String())
}
```

# Players
```
p, err := c.Players()
if err != nil {
	panic(err)
}

for _, player := range p {
	println(player.String())
}
```

# VIPs
```
v, err := c.VIPs()
if err != nil {
	panic(err)
}

for _, vip := range v {
	println(vip.String())
}
```

# Conn

```
type Conn struct {
        // Has unexported fields.
}
    Conn represents a connection to a HLL RCON server. A Conn supports multiple
    thread-safe connections.

func New(addr string, password string) (*Conn, error)
func (c *Conn) AdminAdd(a Admin) error
func (c *Conn) AdminGroups() ([]string, error)
func (c *Conn) AdminRemove(a Admin) error
func (c *Conn) Admins() ([]Admin, error)
func (c *Conn) AutoBalance() (bool, error)
func (c *Conn) AutoBalanceThreshold() (int, error)
func (c *Conn) BanPermanently(p Player, reason, admin string) error
func (c *Conn) BanRemove(p Player) error
func (c *Conn) BanTemporarily(p Player, hours int, reason, admin string) error
func (c *Conn) Close() error
func (c *Conn) IdleTime() (time.Duration, error)
func (c *Conn) Kick(p Player, reason string) error
func (c *Conn) Map() (Map, error)
func (c *Conn) Maps() ([]Map, error)
func (c *Conn) MaxPing() (time.Duration, error)
func (c *Conn) Name() (string, error)
func (c *Conn) PermanentlyBanned() ([]Ban, error)
func (c *Conn) Player(username string) (Player, error)
func (c *Conn) Players() ([]Player, error)
func (c *Conn) Profanities() ([]string, error)
func (c *Conn) Punish(p Player, reason string) error
func (c *Conn) ResetVoteKickThreshold() error
func (c *Conn) Rotation() ([]Map, error)
func (c *Conn) RotationAdd(n MapName) error
func (c *Conn) RotationRemove(n MapName) error
func (c *Conn) SetAutoBalance(enabled bool) error
func (c *Conn) SetAutoBalanceThreshold(diff int) error
func (c *Conn) SetBroadcast(message string) error
func (c *Conn) SetIdleTime(m time.Duration) error
func (c *Conn) SetMap(n MapName) error
func (c *Conn) SetMaxPing(ms time.Duration) error
func (c *Conn) SetProfanities(words ...string) error
func (c *Conn) SetQueueLength(length int) error
func (c *Conn) SetSwitchTeamCooldown(m time.Duration) error
func (c *Conn) SetSwitchTeamNow(p Player) error
func (c *Conn) SetSwitchTeamOnDeath(p Player) error
func (c *Conn) SetVIPSlots(slots int) error
func (c *Conn) SetVoteKick(enabled bool) error
func (c *Conn) SetVoteKickThreshold(pairs ...VoteKickThreshold) error        
func (c *Conn) Slots() (numerator, denominator int, err error)
func (c *Conn) SwitchTeamCooldown() (time.Duration, error)
func (c *Conn) TemporarilyBanned() ([]Ban, error)
func (c *Conn) UnsetProfanities(words ...string) error
func (c *Conn) VIPAdd(v VIP) error
func (c *Conn) VIPRemove(v VIP) error
func (c *Conn) VIPSlots() (int, error)
func (c *Conn) VIPs() ([]VIP, error)
func (c *Conn) VoteKick() (bool, error)
```