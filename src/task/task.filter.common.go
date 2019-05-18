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
