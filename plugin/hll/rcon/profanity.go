package rcon

import (
	"fmt"
	"strings"
)

func (c *Conn) Profanities() ([]string, error) {
	result, err := c.send("get", "profanity")
	if err != nil {
		return nil, fmt.Errorf("failed to get profanities: %v", err)
	}

	args := strings.Split(result, "\t")
	if len(args) == 0 {
		return nil, fmt.Errorf("failed to parse profanities")
	}

	words := []string{}

	for _, word := range args[1 : len(args)-1] {
		words = append(words, word)
	}

	return words, nil
}

func (c *Conn) SetProfanities(words ...string) error {
	_, err := c.send("BanProfanity", strings.Join(words, ","))
	if err != nil {
		return fmt.Errorf("failed to set profanities: %v", err)
	}

	return nil
}

func (c *Conn) UnsetProfanities(words ...string) error {
	_, err := c.send("UnbanProfanity", strings.Join(words, ","))
	if err != nil {
		return fmt.Errorf("failed to remove profanities: %v", err)
	}

	return nil
}
