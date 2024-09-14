package logging

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestSetLoggingLevelPerColor(t *testing.T) {
	expectedColorPerLevel := map[string]color.Attribute{
		// Valid cases
		"DEBUG": color.FgCyan,
		"INFO":  color.FgGreen,
		"WARN":  color.FgYellow,
		"ERROR": color.FgRed,
		// Sample invalid case hits default
		"INVALID": color.FgWhite,
	}
	for level, color := range expectedColorPerLevel {
		assert.Equal(t, color, setOutputColorPerLevel(level))
	}
}
