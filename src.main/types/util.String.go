package types

import "strings"

func StringRepresentsVariadicType(input string) (bool, string) {
	if inputTrimmed := strings.TrimSuffix(input, "..."); inputTrimmed != input {
		return true, input
	} else {
		return false, input
	}
}
