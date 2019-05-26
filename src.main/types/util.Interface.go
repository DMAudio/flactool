package types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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

func InterfaceToInterfaceSlice(input interface{}) ([]interface{}, *Exception) {
	result := make([]interface{}, 0)
	if input == nil {
		return result, nil
	}

	if inputSliceParsed, ok := input.([]interface{}); !ok {
		return nil, Exception_Mismatched_Format("Kind:Slice",
			"Kind:"+reflect.TypeOf(input).Kind().String(),
		)
	} else {
		for _, inputRaw := range inputSliceParsed {
			if inputRaw == nil {
				continue
			}
			result = append(result, inputRaw)
		}
		return result, nil
	}
}

type SliceOpt uint8

const (
	TypeString_Split SliceOpt = 1 << (8 - 1 - iota) // split string into slice by ','
	TypeString_Warp                                 // warp string as [1]string
	TypeString_Error                                // throw error when trying to parse string
)

func InterfaceToStringSlice(input interface{}, option SliceOpt) ([]string, *Exception) {
	result := make([]string, 0)

	if input == nil {
		return result, nil
	}

	switch input.(type) {
	case string:
		switch {
		case option&TypeString_Split != 0:
			inputSliceParsed := strings.Split(input.(string), ",")
			for _, inputRaw := range inputSliceParsed {
				inputRaw = strings.TrimSpace(inputRaw)
				if inputRaw == "" {
					continue
				}
				result = append(result, inputRaw)
			}
			return result, nil
		case option&TypeString_Warp != 0:
			return []string{input.(string)}, nil
		default:
			fallthrough
		case option&TypeString_Error != 0:
			return nil, Exception_Mismatched_Format(
				"Kind:Slice",
				"Type:String",
			)
		}
	case []interface{}:
		if inputSliceParsed, err := InterfaceToInterfaceSlice(input); err != nil {
			return nil, err
		} else {
			for inputIndex, inputRaw := range inputSliceParsed {
				inputRaw = strings.TrimSpace(inputRaw.(string))
				if inputRaw == "" {
					continue
				}
				if inputParsed, ok := inputRaw.(string); !ok {
					return nil, SliceItem_Parsed_Failed(
						inputIndex,
						Exception_Mismatched_Format("Type:string", reflect.TypeOf(inputRaw).String()),
					)
				} else {
					result = append(result, inputParsed)
				}
			}
		}
		return result, nil
	default:
		return nil, Exception_Mismatched_Format(
			"Kind:Slice or Type:string",
			"Kind:"+reflect.TypeOf(input).Kind().String()+", "+
				"Type:"+reflect.TypeOf(input).String(),
		)
	}
}

func InterfaceToStringMap(input interface{}) (map[string]string, *Exception) {
	if inputParsed, ok := input.(map[interface{}]interface{}); !ok {
		return nil, Exception_Mismatched_Format("Kind:map", "Kind:"+reflect.TypeOf(input).Kind().String())
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
					Exception_Mismatched_Format("Type:string or Type:int", "Type:"+reflect.TypeOf(keyRaw).String()),
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
					Exception_Mismatched_Format("Type:string or Type:int", "Type:"+reflect.TypeOf(valueRaw).String()),
				)
			}
			result[keyParsed] = valueParsed
		}
		return result, nil
	}
}
