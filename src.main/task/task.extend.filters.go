package task

import (
	"dadp.flactool/types"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type ArgFilter struct {
	filters map[string]func(string, map[string]interface{}) (string, *types.Exception)
	argRegexp *regexp.Regexp
}

func NewArgFilter() *ArgFilter {
	argRegexp, _ := regexp.Compile("{@([a-zA-Z0-9]+):([^{}]*)}")
	return &ArgFilter{
		filters: map[string]func(string, map[string]interface{}) (string, *types.Exception){},
		argRegexp: argRegexp,
	}
}

func (af *ArgFilter) FillArgs(raw string, extraArgsCollection map[string]map[string]interface{}) (string, int, *types.Exception) {
	argList := af.argRegexp.FindAllStringSubmatch(raw, -1)
	if !(argList != nil && len(argList) > 0) {
		return raw, 0, nil
	}
	for _, arg := range argList {
		if len(arg) == 0 {
			continue
		}
		if len(arg) != 3 {
			return raw, 0, types.NewException(TMFilter_UnableToParse_Arg, map[string]string{
				"reason": "参数 %s 格式不完整",
			}, nil)
		}
		if arg[1] = strings.TrimSpace(arg[1]); arg[1] == "" {
			return raw, 0, types.NewException(TMFilter_UnableToParse_Arg, map[string]string{
				"reason": fmt.Sprintf("参数 %s 格式不完整: 处理者不能为空", arg[0]),
			}, nil)
		}
		if arg[2] = strings.TrimSpace(arg[2]); arg[2] == "" {
			return raw, 0, types.NewException(TMFilter_UnableToParse_Arg, map[string]string{
				"reason": fmt.Sprintf("参数 %s 格式不完整: 参数不能为空", arg[0]),
			}, nil)
		}
	}

	exceptions := map[string]*types.Exception{}
	for argIndex, argRawSlice := range argList {
		argRaw, argFilterHandler, argFilterArgs := argRawSlice[0], argRawSlice[1], argRawSlice[2]

		var extraArgs map[string]interface{} = nil
		if extraArgsCollection != nil {
			if extraArgsTmp, exist := extraArgsCollection[argFilterHandler]; exist {
				extraArgs = extraArgsTmp
			}
		}

		if result, err := af.ParseArg(argFilterHandler, argFilterArgs, extraArgs); err != nil {
			exceptions["A"+strconv.Itoa(argIndex)] = err
		} else {
			raw = strings.ReplaceAll(raw, argRaw, result)
		}
	}

	if len(exceptions) != 0 {
		return raw, len(argList) - len(exceptions), types.NewException(TMFilter_FailedToFill_Args, nil, exceptions)
	}

	if argList := af.argRegexp.FindAllStringSubmatch(raw, -1); len(argList) == 0 {
		return raw, len(argList), nil
	} else if rawProcessed, agrAmount, err := af.FillArgs(raw, extraArgsCollection); err != nil {
		return rawProcessed, len(argList) + agrAmount, err
	} else {
		return rawProcessed, len(argList) + agrAmount, nil
	}

}

func (af *ArgFilter) ParseArg(filterName string, args string, extraArgs map[string]interface{}) (string, *types.Exception) {
	if filter, exist := af.filters[filterName]; !exist {
		return "", types.NewException(TMFilter_Undefined_Handler, map[string]string{
			"handler": filterName,
		}, nil)
	} else if result, err := filter(args, extraArgs); err != nil {
		return "", types.NewException(TMFilter_FailedToExecute_Filter, map[string]string{
			"handler": filterName,
			"args":    args,
		}, err)
	} else {
		return result, nil
	}
}

func (af *ArgFilter) Register(filterName string, handler func(string, map[string]interface{}) (string, *types.Exception)) *types.Exception {
	if _, alreadyExist := af.filters[filterName]; alreadyExist {
		return types.NewException(TMFilter_CanNotRegister, map[string]string{
			"reason": "执行者名称冲突",
		}, nil)
	}
	af.filters[filterName] = handler
	return nil
}

var globalArgFilter *ArgFilter
var globalArgFilterLock sync.Mutex

func GlobalArgFilter() *ArgFilter {
	if globalArgFilter == nil {
		globalArgFilterLock.Lock()
		defer globalArgFilterLock.Unlock()
		if globalArgFilter == nil {
			globalArgFilter = NewArgFilter()
		}
	}
	return globalArgFilter
}
