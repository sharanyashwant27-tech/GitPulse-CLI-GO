package cmd

import (
	"strings"
)

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func pad(s string, n int) string {
	if len(s) >= n {
		return truncate(s, n)
	}
	return s + strings.Repeat(" ", n-len(s))
}
