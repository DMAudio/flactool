package task

import (
	"fmt"
	"p20190417/types"
	"sync"
)

type Handler struct {
	list map[string]func(string, interface{}) (interface{}, *types.Exception)
}

var TMTask_Output = types.NewMask(
	"TASK_OUTPUT",
	"\n{{value}}",
)

func (h *Handler) Register(handler string, handlerFunc func(string, interface{}) (interface{}, *types.Exception)) *types.Exception {
	if _, exist := h.list[handler]; exist {
		return types.NewException(TMTask_CanNotRegister, map[string]string{
			"reason": "执行者名称冲突",
		}, nil)
	}

	h.list[handler] = handlerFunc

	return nil
}

func (h *Handler) Execute(handler string, operation string, args interface{}) *types.Exception {
	if handlerFunc, exist := h.list[handler]; !exist {
		return types.NewException(TMTask_CanNotExecute, map[string]string{
			"reason": "执行者不存在",
		}, nil)
	} else if output, err := handlerFunc(operation, args); err != nil {
		return err
	} else if output != nil {
		types.Throw(types.NewException(TMTask_Output, map[string]string{
			"value": fmt.Sprintf("%v", output),
		}, nil), types.RsInfo)
	}

	return nil
}

var globalHandler *Handler
var globalHandlerLock sync.Mutex

func GlobalHandler() *Handler {
	if globalHandler == nil {
		globalHandlerLock.Lock()
		defer globalHandlerLock.Unlock()
		if globalHandler == nil {
			globalHandler = &Handler{
				list: map[string]func(string, interface{}) (interface{}, *types.Exception){},
			}
		}
	}
	return globalHandler
}
