package utils

import (
	"testing"

	"github.com/fatih/color"
)

func TestLogger(t *testing.T) {
	color.NoColor = false
	LogInfo("这是log info")
	LogWarn("这是log warn")
	LogError("这是log error")
	LogDebug("这是log debug")
}
