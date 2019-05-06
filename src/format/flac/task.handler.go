package flac

import (
	"p20190417/types"
	"sync"
)

var globalFlac *Flac
var globalFlacLock sync.Mutex

func GlobalFlac() *Flac {
	if globalFlac == nil {
		globalFlacLock.Lock()
		defer globalFlacLock.Unlock()
		if globalFlac == nil {
			globalFlac = &Flac{}
		}
	}
	return globalFlac
}

const TaskHandler_T4VORB_Key = "FLAC_VORBIS_COMMENT"

func TaskHandler_T4VORB(operation string, args interface{}) (interface{}, *types.Exception) {
	switch operation {
	case "printRefer":
		return TaskHandler_T4VORB_PrintRefer(args)
	case "setRefer":
		return TaskHandler_T4VORB_SetRefer(args)
	case "printTags":
		return TaskHandler_T4VORB_PrintTags(args)
	case "setTags":
		return TaskHandler_T4VORB_SetTags(args)
	case "dumpTags":
		return TaskHandler_T4VORB_dumpTags(args)
	case "importTags":
		return TaskHandler_T4VORB_importTags(args)
	case "deleteTags":
		return TaskHandler_T4VORB_DeleteTags(args)
	case "sortTags":
		return TaskHandler_T4VORB_SortTags(args)
	}

	return nil, nil
}

func TMFlac_Task_Arguments_Format_Exception(expected string, got string) *types.Exception {
	return types.NewException(types.NewMask(
		"ARGUMENTS_FORMAT_ERROR",
		"参数格式错误：预期：{{expected}}，实际：{{got}}",
	), map[string]string{
		"expected": expected,
		"got":      got,
	}, nil)
}
