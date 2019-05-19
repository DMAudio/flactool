package task

import (
	"os"
	"p20190417/types"
	"strings"
)

func Filter_Env(input string, extraArgs map[string]interface{}) (string, *types.Exception) {
	input = strings.TrimSpace(input)
	return os.Getenv(input), nil
}

func Filter_FmtFileName(input string, extraArgs map[string]interface{}) (string, *types.Exception) {
	input = strings.ReplaceAll(input,"\u003A", "\uFF1A")
	input = strings.ReplaceAll(input,"\u002F", "\uFF0F")
	input = strings.ReplaceAll(input,"\u005C", "\uFF3C")
	input = strings.ReplaceAll(input,"\u003F", "\uFF1F")
	input = strings.ReplaceAll(input,"\u0022", "\u0027\u0027")
	input = strings.ReplaceAll(input,"\u002A", "\uFF0A")
	input = strings.ReplaceAll(input,"\u003C", "\uFF1C")
	input = strings.ReplaceAll(input,"\u003E", "\uFF1E")
	input = strings.ReplaceAll(input,"\u007C", "\uFF5C")
	return input, nil
}
