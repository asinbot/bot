package rcon

import (
	"fmt"
	"strings"
)

// q will quote a string s and return "s".
func q(s string) string {
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		return s
	}

	return fmt.Sprintf(`"%s"`, s)
}
