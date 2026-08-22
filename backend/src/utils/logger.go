package utils

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

const TIME_FORMAT = "2006/01/02 15:04:05"

func timestamp() string {
	return time.Now().Format(TIME_FORMAT)
}

func formatString(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func LogInfo(format string, args ...any) {
	msg := formatString(format, args...)
	color.Green("[%s] [INFO]  %s", timestamp(), msg)
}

func LogWarn(format string, args ...any) {
	msg := formatString(format, args...)
	color.Yellow("[%s] [WARN]  %s", timestamp(), msg)
}

func LogError(format string, args ...any) {
	msg := formatString(format, args...)
	color.Red("[%s] [ERROR] %s", timestamp(), msg)
}

func LogDebug(format string, args ...any) {
	msg := formatString(format, args...)
	color.Cyan("[%s] [DEBUG] %s", timestamp(), msg)
}
