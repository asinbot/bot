package rcon

import (
	"fmt"
	"strings"
)

// Admin represents a Player with elevated privileges.
type Admin struct {
	Player
	Role string
}

var unknownAdmin = Admin{
	Player: Player{
		Name: "unknown admin",
		ID64: "unknown id",
	},
	Role: "unknown role",
}

// Admins will return a slice of active Admins.
func (c *Conn) Admins() ([]Admin, error) {
	result, err := c.send("get", "adminids")
	if err != nil {
		return nil, fmt.Errorf("failed to get information for all admins: %v", err)
	}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse information for all admins")
	}

	admins := []Admin{}

	for _, a := range args[1 : len(args)-1] {
		admins = append(admins, adminFromString(a))
	}

	return admins, nil
}

// Add will add an Admin.
func (c *Conn) AdminAdd(a Admin) error {
	_, err := c.send("adminadd", q(a.ID64), q(a.Role), q(a.Name))
	if err != nil {
		return fmt.Errorf("failed to add admin %s: %v", a.String(), err)
	}

	return nil
}

// Remove will remove an Admin.
func (c *Conn) AdminRemove(a Admin) error {
	_, err := c.send("admindel", a.ID64)
	if err != nil {
		return fmt.Errorf("failed to remove admin %s: %v", a.String(), err)
	}

	return nil
}

// AdminGroups will return existing admin groups.
func (c *Conn) AdminGroups() ([]string, error) {
	result, err := c.send("get", "admingroups")
	if err != nil {
		return nil, fmt.Errorf("failed to get information for all admin groups: %v", err)
	}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse information for all admin groups")
	}

	return args[1 : len(args)-1], nil
}

func (a Admin) String() string {
	return fmt.Sprintf("%s [%s]", a.Player.String(), a.Role)
}

func adminFromString(s string) Admin {
	args := strings.Split(s, " ")
	if len(args) != 3 {
		return Admin{}
	}

	return Admin{
		Player: Player{
			ID64: args[0],
			Name: args[2],
		},
		Role: args[1],
	}
}
