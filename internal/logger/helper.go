package logger

import (
	"github.com/anpotashev/mpdgo/internal/commands"
	"strings"
)

func Truncate(msg string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}
	r := []rune(msg)
	if len(r) <= maxLength {
		return msg
	}
	return string(r[:maxLength]) + "…"
}

func JoinAndTruncateSingleCommands(elems []commands.SingleCommand, sep string, maxLength int) string {
	if maxLength <= 0 {
		return ""
	}
	var b strings.Builder
	length := 0
	for i, e := range elems {
		if i > 0 {
			length += len([]rune(sep))
			if length > maxLength {
				break
			}
			b.WriteString(sep)
		}

		s := e.String()
		runes := []rune(s)

		if length+len(runes) > maxLength {
			remain := maxLength - length
			if remain > 0 {
				b.WriteString(string(runes[:remain]))
			}
			b.WriteString("…")
			return b.String()
		}

		b.WriteString(s)
		length += len(runes)
	}
	return b.String()
}
