package types

import (
	"reflect"
	"strconv"
)

func Mismatched_Format_Exception(expected string, got string) *Exception {
	return NewException(NewMask(
		"Mismatched_FORMAT_ERROR",
		"格式错误：预期：{{expected}}，实际：{{got}}",
	), map[string]string{
		"expected": expected,
		"got":      got,
	}, nil)
}

func SliceItem_Parsed_Failed(position int, cause *Exception) *Exception {
	return NewException(NewMask(
		"SLICE_ITEM_PARSED_FAILED",
		"数组的第 {{position}} 个元素解析失败",
	), map[string]string{
		"position": strconv.Itoa(position),
	}, cause)
}

func InterfaceToStringSlice(input interface{}) ([]string, *Exception) {
	if inputSliceParsed, ok := input.([]interface{}); !ok {
		return nil, Mismatched_Format_Exception("Kind:Slice",
			"Kind:"+reflect.TypeOf(input).Kind().String(),
		)
	} else {
		result := make([]string, 0)
		for inputIndex, inputRaw := range inputSliceParsed {
			if inputRaw == nil {
				continue
			}
			if inputParsed, ok := inputRaw.(string); !ok {
				return nil, SliceItem_Parsed_Failed(inputIndex,
					Mismatched_Format_Exception("string",
						reflect.TypeOf(inputRaw).String(),
					))
			} else {
				result = append(result, inputParsed)
			}
		}
		return result, nil
	}
}
