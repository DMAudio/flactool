package types

import "strings"

func WrappedMessage(contentPrefix string, content string, withWarpper bool) string {
	result := NewBuffer()
	contentPrefix = strings.Trim(contentPrefix, "\n")

	splitContent := strings.Split(strings.TrimSpace(content), "\n")
	_, _ = result.WriteStrings(contentPrefix, splitContent[0], "\n")

	if !withWarpper || len(contentPrefix) <= 1 {
		linePrefix := strings.Repeat(" ", len(contentPrefix))
		for i := 1; i <= len(splitContent)-1; i++ {
			line := splitContent[i]
			_, _ = result.WriteStrings(linePrefix, line, "\n")
		}
	} else {
		linePrefix := strings.Repeat(" ", len(contentPrefix)-1)
		for i := 1; i <= len(splitContent)-1; i++ {
			line := splitContent[i]
			if i < len(splitContent)-1 {
				_, _ = result.WriteStrings("│", linePrefix, line, "\n")
			} else {
				_, _ = result.WriteStrings("└", linePrefix, line, "\n")
			}
		}
	}
	return strings.TrimSpace(result.String())
}

type ThrowableString string

func (m ThrowableString) GetMessage(bool) string {
	return string(m)
}
