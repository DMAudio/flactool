package task

import (
	"dadp.flactool/task/cgoFillers"
	"dadp.flactool/types"
	"math"
	"strconv"
	"strings"
	"sync"
)

type ArgFiller struct {
	fillers map[string]func(string, map[string]interface{}) (string, *types.Exception)
}

func NewArgFiller() *ArgFiller {
	return &ArgFiller{
		fillers: map[string]func(string, map[string]interface{}) (string, *types.Exception){},
	}
}

func (af *ArgFiller) FillArgs(raw string, extraParameters map[string]map[string]interface{}) (string, int, *types.Exception) {
	argList := cgoFillers.GetArgs(raw)
	if argList == nil || argList.Size() == 0 {
		return raw, 0, nil
	} else if argList.Size() > int64(math.MaxInt32) {
		return raw, 0, types.NewException(TMFiller_UnableToParse_Arg, map[string]string{
			"reason": "参数过多",
		}, nil)
	}

	argListSize := int(argList.Size())

	exceptions := map[string]*types.Exception{}
	for argIndex := 0; argIndex <= argListSize-1; argIndex++ {
		argRawItem := argList.Get(argIndex)
		argRaw, argFillerHandler, argFillerArgs := argRawItem.Get(0), argRawItem.Get(1), argRawItem.Get(2)

		var extraArgs map[string]interface{} = nil
		if extraParameters != nil {
			if extraArgsTmp, exist := extraParameters[argFillerHandler]; exist {
				extraArgs = extraArgsTmp
			}
		}

		if result, err := af.ParseArg(argFillerHandler, argFillerArgs, extraArgs); err != nil {
			exceptions["A"+strconv.Itoa(argIndex)] = err
		} else {
			raw = strings.ReplaceAll(raw, argRaw, result)
		}
	}

	if len(exceptions) != 0 {
		return raw, argListSize - len(exceptions), types.NewException(TMFiller_FailedToFill_Args, nil, exceptions)
	}

	if argList := cgoFillers.GetArgs(raw); argList.Size() == 0 {
		return raw, argListSize, nil
	} else if rawProcessed, agrAmount, err := af.FillArgs(raw, extraParameters); err != nil {
		return rawProcessed, argListSize + agrAmount, err
	} else {
		return rawProcessed, argListSize + agrAmount, nil
	}
}

func (af *ArgFiller) ParseArg(fillerName string, parameter string, extraArgs map[string]interface{}) (string, *types.Exception) {
	if handler, exist := af.fillers[fillerName]; !exist {
		return "", types.NewException(TMFiller_Undefined_Handler, map[string]string{
			"handler": fillerName,
		}, nil)
	} else if result, err := handler(parameter, extraArgs); err != nil {
		return "", types.NewException(TMFiller_FailedToExecute_Filler, map[string]string{
			"handler":    fillerName,
			"parameters": parameter,
		}, err)
	} else {
		return result, nil
	}
}

func (af *ArgFiller) Register(fillerName string, filler func(string, map[string]interface{}) (string, *types.Exception)) *types.Exception {
	if _, alreadyExist := af.fillers[fillerName]; alreadyExist {
		return types.NewException(TMFiller_CanNotRegister, map[string]string{
			"reason": "执行者名称冲突",
		}, nil)
	}
	af.fillers[fillerName] = filler
	return nil
}

var globalArgFiller *ArgFiller
var globalArgFillerLock sync.Mutex

func GlobalArgFiller() *ArgFiller {
	if globalArgFiller == nil {
		globalArgFillerLock.Lock()
		defer globalArgFillerLock.Unlock()
		if globalArgFiller == nil {
			globalArgFiller = NewArgFiller()
		}
	}
	return globalArgFiller
}
