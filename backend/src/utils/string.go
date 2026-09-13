package utils

import (
	"unicode/utf8"
)

func TruncateString(s string, maxLen int) string {
	runes := []rune(s)

	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}

	return s
}

func TruncateBytes(s string, maxBytes int) string {
	if len(s) <= maxBytes {
		return s
	}

	trancated := s[:maxBytes]

	// 检查是否是合法的 utf-8字符串
	if utf8.ValidString(trancated) {
		return trancated
	}

	// 非法，向前回退
	for !utf8.ValidString(trancated) {
		maxBytes--
		trancated = s[:maxBytes]
	}

	return trancated
}
