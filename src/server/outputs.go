package server

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintStandardMessage(source string, message string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("[%s] %s\n", green(source), message)
}

func PrintWarn(text string, args ...any) {
	yellow := color.New(color.FgYellow).SprintFunc()
	message := fmt.Sprintf(text, args...)
	fmt.Printf("[%s] %s\n", yellow("WARNING"), message)
}

func PrintError(text string, args ...any) {
	red := color.New(color.FgRed).SprintFunc()
	message := fmt.Sprintf(text, args...)
	fmt.Printf("[%s] %s\n", red("ERROR"), message)
}
