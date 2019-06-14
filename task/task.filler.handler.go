package task

import (
	"gitlab.com/MGEs/Com.Base/types"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func Filler_Env(input string, extraArgs map[string]interface{}) (string, *types.Exception) {
	input = strings.TrimSpace(input)
	return os.Getenv(input), nil
}

func Filler_FmtFileName(input string, extraArgs map[string]interface{}) (string, *types.Exception) {
	input = strings.ReplaceAll(input, "\u003A", "\uFF1A")
	input = strings.ReplaceAll(input, "\u002F", "\uFF0F")
	input = strings.ReplaceAll(input, "\u005C", "\uFF3C")
	input = strings.ReplaceAll(input, "\u003F", "\uFF1F")
	input = strings.ReplaceAll(input, "\u0022", "\u0027\u0027")
	input = strings.ReplaceAll(input, "\u002A", "\uFF0A")
	input = strings.ReplaceAll(input, "\u003C", "\uFF1C")
	input = strings.ReplaceAll(input, "\u003E", "\uFF1E")
	input = strings.ReplaceAll(input, "\u007C", "\uFF5C")
	return input, nil
}

var TMArg_FailedTo_Parse_URISequence = types.NewMask(
	"FailedTo_Parse_URISequence",
	"无法解析编码的URI序列",
)

func Filler_DecodeURI(input string, extraArgs map[string]interface{}) (string, *types.Exception) {
	input = strings.ToUpper(input)
	if matched, err := regexp.MatchString("^([0-9A-F]{2})+$", input); err != nil {
		return "", types.NewException(TMFiller_FailedTo_CompileRegex, nil, err)
	} else if !matched {
		return "", types.Exception_Mismatched_Format("Regex:$([0-9A-Z]{2})+^", input)
	}

	encodeBuffer := types.NewBuffer()
	if regex, err := regexp.Compile("[0-9A-F]{2}"); err != nil {
		return "", types.NewException(TMFiller_FailedTo_CompileRegex, nil, err)
	} else {
		for _, charCode := range regex.FindAllString(input, -1) {
			encodeBuffer.WriteStringsIE("%", charCode)
		}
	}

	if decodedValue, err := url.QueryUnescape(encodeBuffer.String()); err != nil {
		return "", types.NewException(TMArg_FailedTo_Parse_URISequence, nil, err)
	} else {
		return decodedValue, nil
	}
}
