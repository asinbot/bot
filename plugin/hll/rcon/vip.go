package rcon

import (
	"fmt"
	"strconv"
	"strings"
)

// VIP represents a Player with some elevated privileges.
type VIP struct {
	Player
}

// SetVIPSlots will update the current number of joinable VIPs.
func (c *Conn) SetVIPSlots(slots int) error {
	_, err := c.send("setnumvipslots", strconv.Itoa(slots))
	if err != nil {
		return fmt.Errorf("failed to set vip slots: %v", err)
	}

	return nil
}

var x = make(chan bool)

// Admins will return a slice of active Admins.
func (c *Conn) VIPs() ([]VIP, error) {
	result, err := c.send("get", "vipids")
	if err != nil {
		return nil, fmt.Errorf("failed to get information for all vips: %v", err)
	}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse information for all vips")
	}

	vips := []VIP{}

	for _, v := range args[1 : len(args)-1] {
		nameID := strings.Split(v, " ")

		vips = append(vips, VIP{
			Player: Player{
				Name: strings.Trim(strings.TrimSpace(nameID[1]), `"`),
				ID64: strings.TrimSpace(nameID[0]),
			},
		})
	}

	return vips, nil
}

// Add will add a new VIP.
func (c *Conn) VIPAdd(v VIP) error {
	_, err := c.send("vipadd", q(v.ID64), q(v.Name))
	if err != nil {
		return fmt.Errorf("failed to add vip %s: %v", v.String(), err)
	}

	return nil
}

// Remove will remove a VIP.
func (c *Conn) VIPRemove(v VIP) error {
	_, err := c.send("vipdel", v.ID64)
	if err != nil {
		return fmt.Errorf("failed to remove vip %s: %v", v.String(), err)
	}

	return nil
}

// VIPSlots will return the current number of joinable VIPs.
func (c *Conn) VIPSlots() (int, error) {
	result, err := c.send("get", "numvipslots")
	if err != nil {
		return -1, fmt.Errorf("failed to get vip slots: %v", err)
	}

	s, err := strconv.Atoi(result)
	if err != nil {
		return -1, fmt.Errorf("failed to parse information for vip slots: %v", err)
	}

	return s, nil
}
