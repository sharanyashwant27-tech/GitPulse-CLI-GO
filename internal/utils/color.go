package utils

import "github.com/fatih/color"

// Successf prints a green success line to stdout.
func Successf(format string, a ...any) {
	color.New(color.FgGreen, color.Bold).Printf(format, a...)
}

// Warnf prints a yellow warning line to stdout.
func Warnf(format string, a ...any) {
	color.New(color.FgYellow).Printf(format, a...)
}

// Errorf prints a red error line to stdout.
func Errorf(format string, a ...any) {
	color.New(color.FgRed, color.Bold).Printf(format, a...)
}

// Cyanf prints a cyan info line to stdout.
func Cyanf(format string, a ...any) {
	color.New(color.FgCyan).Printf(format, a...)
}
