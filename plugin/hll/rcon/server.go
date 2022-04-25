package rcon

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// VoteKickThreshold represents a list item for vote kick threshold updates.
type VoteKickThreshold struct {
	Players   int
	Threshold int
}

// Name returns the name of the server.
func (c *Conn) Name() (string, error) {
	result, err := c.send("get", "name")
	if err != nil {
		return "", fmt.Errorf("failed to get server name: %v", err)
	}

	return result, nil
}

// IdleTime returns the maximum time a player can be idle in a server.
func (c *Conn) IdleTime() (time.Duration, error) {
	result, err := c.send("get", "idletime")
	if err != nil {
		return -1, fmt.Errorf("failed to get idle auto-kick time: %v", err)
	}

	t, err := strconv.Atoi(result)
	if err != nil {
		return -1, fmt.Errorf("failed to parse idle auto-kick time: %v", err)
	}

	return time.Duration(t) * time.Minute, nil
}

// MaxPing returns the maxiumum RTT time (milliseconds) a player can have with a server.
func (c *Conn) MaxPing() (time.Duration, error) {
	result, err := c.send("get", "highping")
	if err != nil {
		return -1, fmt.Errorf("failed to get max ping auto-kick threshold: %v", err)
	}

	p, err := strconv.Atoi(result)
	if err != nil {
		return -1, fmt.Errorf("failed to parse max ping auto-kick threshold: %v", err)
	}

	return time.Duration(p) * time.Millisecond, nil
}

// AutoBalance will return the server configuration for auto-balancing teams.
func (c *Conn) AutoBalance() (bool, error) {
	result, err := c.send("get", "autobalanceenabled")
	if err != nil {
		return false, fmt.Errorf("failed to set vote kick configuration: %v", err)
	}

	return result == "on", nil
}

// SetAutoBalance will update the server configuration for auto-balancing teams.
func (c *Conn) SetAutoBalance(enabled bool) error {
	_, err := c.send("setautobalanceenabled", map[bool]string{true: "on", false: "off"}[enabled]) // Ternary!
	if err != nil {
		return fmt.Errorf("failed to set auto balance configuration: %v", err)
	}

	return nil
}

// SetAutoBalanceThreshold will update the delta required for the server to autobalance teams.
func (c *Conn) SetAutoBalanceThreshold(diff int) error {
	_, err := c.send("setautobalancethreshold", strconv.Itoa(diff))
	if err != nil {
		return fmt.Errorf("failed to set auto balance threshold configuration: %v", err)
	}

	return nil
}

// SwitchTeamCooldown will return the time (minutes) before a player can switch teams.
func (c *Conn) SwitchTeamCooldown() (time.Duration, error) {
	result, err := c.send("get", "teamswitchcooldown")
	if err != nil {
		return -1, fmt.Errorf("failed to set switch team cooldown configuration: %v", err)
	}

	s, err := strconv.Atoi(result)
	if err != nil {
		return -1, fmt.Errorf("failed to parse switch team cooldown configuration: %v", err)
	}

	return time.Duration(s) * time.Minute, nil
}

// AutoBalanceThreshold will return the delta required for the server to autobalance teams.
func (c *Conn) AutoBalanceThreshold() (int, error) {
	result, err := c.send("get", "autobalancethreshold")
	if err != nil {
		return -1, fmt.Errorf("failed to get auto balance threshold configuration: %v", err)
	}

	t, err := strconv.Atoi(result)
	if err != nil {
		return -1, fmt.Errorf("failed to parse auto balance threshold configuration: %v", err)
	}

	return t, nil
}

// SetSwitchTeamCooldown will update the time (minutes) before a player can switch teams.
func (c *Conn) SetSwitchTeamCooldown(m time.Duration) error {
	_, err := c.send("setteamswitchcooldown", strconv.Itoa(int(m.Minutes())))
	if err != nil {
		return fmt.Errorf("failed to set switch team cooldown configuration: %v", err)
	}

	return nil
}

// SetIdleTime updates the maximum time (minutes) a player can be idle in a server (0 to disable).
func (c *Conn) SetIdleTime(m time.Duration) error {
	_, err := c.send("setkickidletime", strconv.Itoa(int(m.Minutes())))
	if err != nil {
		return fmt.Errorf("failed to set idle auto-kick time: %v", err)
	}

	return nil
}

// SetMaxPing updates the maxiumum RTT time (milliseconds) a player can have with a server (0 to disable).
func (c *Conn) SetMaxPing(ms time.Duration) error {
	_, err := c.send("sethighping", strconv.Itoa(int(ms.Milliseconds())))
	if err != nil {
		return fmt.Errorf("failed to set max ping auto-kick time: %v", err)
	}

	return nil
}

// SetQueueLength will update the current number of players allowed to queue.
func (c *Conn) SetQueueLength(length int) error {
	_, err := c.send("setmaxqueuedplayers", strconv.Itoa(length))
	if err != nil {
		return fmt.Errorf("failed to set max queued players: %v", err)
	}

	return nil
}

// QueueLength will return the current number of players allowed to queue.
func (c *Conn) QueueLength() (int, error) {
	result, err := c.send("get", "maxqueuedplayers")
	if err != nil {
		return -1, fmt.Errorf("failed to get queue length: %v", err)
	}

	q, err := strconv.Atoi(result)
	if err != nil {
		return -1, fmt.Errorf("failed to parse queue length configuration: %v", err)
	}

	return q, nil
}

// SetBroadcast will update the current broadcast message.
func (c *Conn) SetBroadcast(message string) error {
	_, err := c.send("broadcast", q(message))
	if err != nil {
		return fmt.Errorf("failed to set broadcast message: %v", err)
	}

	return nil
}

// Slots will return the current number of slots.
func (c *Conn) Slots() (numerator, denominator int, err error) {
	result, err := c.send("get", "slots")
	if err != nil {
		return 0, -1, fmt.Errorf("failed to get slots: %v", err)
	}

	numDen := strings.Split(result, "/")
	if len(numDen) != 2 {
		return 0, -1, fmt.Errorf("failed to parse slots configuration: %v", err)
	}

	numerator, err = strconv.Atoi(numDen[0])
	if err != nil {
		return 0, -1, fmt.Errorf("failed to parse numerator for slots configuration: %v", err)
	}

	denominator, err = strconv.Atoi(numDen[1])
	if err != nil {
		return 0, -1, fmt.Errorf("failed to parse denominator for slots configuration: %v", err)
	}

	return
}

// VoteKick will return the ability for players to vote kick eachother in a server.
func (c *Conn) VoteKick() (bool, error) {
	result, err := c.send("get", "votekickenabled")
	if err != nil {
		return false, fmt.Errorf("failed to set votekick configuration: %v", err)
	}

	return result == "on", nil
}

// VoteKickThreshold will return the current votekick threshold
func (c *Conn) VoteKickThreshold() (string, error) {
	result, err := c.send("get", "votekickthreshold")
	if err != nil {
		return "", fmt.Errorf("failed to get vote kick threshold: %v", err)
	}

	return result, nil
}

// SetVoteKick will update the ability for players to vote kick eachother in a server.
func (c *Conn) SetVoteKick(enabled bool) error {
	_, err := c.send("setvotekickenabled", map[bool]string{true: "on", false: "off"}[enabled]) // Ternary!
	if err != nil {
		return fmt.Errorf("failed to set votekick configuration: %v", err)
	}

	return nil
}

// SetVoteKickThreshold will update the current votekick thresholds.
func (c *Conn) SetVoteKickThreshold(pairs ...VoteKickThreshold) error {
	threshold := ""
	for i, pair := range pairs {
		threshold += fmt.Sprintf("%d,%d", pair.Players, pair.Threshold)
		if i != len(pairs)-1 {
			threshold += ","
		}
	}

	_, err := c.send("setvotekickthreshold", threshold)
	if err != nil {
		return fmt.Errorf("failed to set votekick threshold configuration: %v", err)
	}

	return nil
}

// ResetVoteKickThreshold .
func (c *Conn) ResetVoteKickThreshold() error {
	_, err := c.send("resetvotekickthreshold")
	if err != nil {
		return fmt.Errorf("failed to reset votekick threshold configuration: %v", err)
	}

	return nil
}
