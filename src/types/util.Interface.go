package types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

func MapKey_Parsed_Failed(key interface{}, cause *Exception) *Exception {
	return NewException(NewMask(
		"SLICE_Key_PARSED_FAILED",
		"键值对中的键{{key}}解析失败",
	), map[string]string{
		"position": fmt.Sprintf("%v", key),
	}, cause)
}

func MapItem_Parsed_Failed(key string, cause *Exception) *Exception {
	return NewException(NewMask(
		"SLICE_ITEM_PARSED_FAILED",
		"键值对中键为 {{key}} 的元素解析失败",
	), map[string]string{
		"key": key,
	}, cause)
}

func InterfaceToStringSlice(input interface{}) ([]string, *Exception) {
	result := make([]string, 0)

	switch input.(type) {
	case string:
		inputSliceParsed := strings.Split(input.(string), ",")
		for _, inputRaw := range inputSliceParsed {
			inputRaw = strings.TrimSpace(inputRaw)
			if inputRaw == "" {
				continue
			}
			result = append(result, inputRaw)
		}
		return result, nil
	case []interface{}:
		if inputSliceParsed, ok := input.([]interface{}); !ok {
			return nil, Mismatched_Format_Exception("Kind:Slice",
				"Kind:"+reflect.TypeOf(input).Kind().String(),
			)
		} else {
			for inputIndex, inputRaw := range inputSliceParsed {
				if inputRaw == nil {
					continue
				}
				inputRaw = strings.TrimSpace(inputRaw.(string))
				if inputRaw == "" {
					continue
				}
				if inputParsed, ok := inputRaw.(string); !ok {
					return nil, SliceItem_Parsed_Failed(
						inputIndex,
						Mismatched_Format_Exception("string", reflect.TypeOf(inputRaw).String()),
					)
				} else {
					result = append(result, inputParsed)
				}
			}
		}
		return result, nil
	default:
		return nil, Mismatched_Format_Exception("Kind:Slice or String", reflect.TypeOf(input).String())
	}

}

func InterfaceToStringMap(input interface{}) (map[string]string, *Exception) {
	if inputParsed, ok := input.(map[interface{}]interface{}); !ok {
		return nil, Mismatched_Format_Exception("Kind:map", "Kind:"+reflect.TypeOf(input).Kind().String())
	} else {
		result := map[string]string{}
		for keyRaw, valueRaw := range inputParsed {
			var keyParsed string
			if keyRaw == nil {
				continue
			}
			switch keyRaw.(type) {
			case string:
				keyParsed = keyRaw.(string)
			case int:
				keyParsed = strconv.FormatInt(int64(keyRaw.(int)), 10)
			default:
				return nil, MapKey_Parsed_Failed(keyRaw,
					Mismatched_Format_Exception("string / int", reflect.TypeOf(keyRaw).String()),
				)
			}

			var valueParsed string
			if valueRaw == nil {
				continue
			}
			switch valueRaw.(type) {
			case string:
				valueParsed = valueRaw.(string)
			case int:
				valueParsed = strconv.FormatInt(int64(valueRaw.(int)), 10)
			default:
				return nil, MapItem_Parsed_Failed(keyParsed,
					Mismatched_Format_Exception("string / int", reflect.TypeOf(valueRaw).String()),
				)
			}
			result[keyParsed] = valueParsed
		}
		return result, nil
	}
}
