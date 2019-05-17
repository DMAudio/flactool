package flac

import (
	"p20190417/task"
	"p20190417/types"
)

func TaskHandler_Exception_UnsupportedTask(operation string) *types.Exception {
	return types.NewException(task.TMConfig_Unsupported_TaskOperation, map[string]string{
		"operation": operation,
	}, nil)
}

const TaskHandler_T4VORB_Key = MetaBlockTypeStr_VORBIS_COMMENT

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
	default:
		return nil, TaskHandler_Exception_UnsupportedTask(operation)
	}
}

const TaskHandler_T6PICT_Key = MetaBlockTypeStr_PICTURE

func TaskHandler_T6PICT(operation string, args interface{}) (interface{}, *types.Exception) {
	switch operation {
	case "dumpPic":
		return TaskHandler_T6PICT_dumpPic(args)
	case "setPic":
		return TaskHandler_T6PICT_setPic(args)
	case "addPic":
		return TaskHandler_T6PICT_addPic(args)
	case "setPicType":
		return TaskHandler_T6PICT_setPicType(args)
	case "setDesc":
		return TaskHandler_T6PICT_setDesc(args)
	case "getDesc":
		return TaskHandler_T6PICT_getDesc(args)
	default:
		return nil, TaskHandler_Exception_UnsupportedTask(operation)
	}
}

const TaskHandler_BLOCKS_Key = "BLOCKS"

func TaskHandler_BLOCKS(operation string, args interface{}) (interface{}, *types.Exception) {
	switch operation {
	case "sortBlocks":
		return TaskHandler_MAIN_SortBlocks(args)
	case "deleteBlocks":
		return TaskHandler_MAIN_DeleteBlocks(args)
	default:
		return nil, TaskHandler_Exception_UnsupportedTask(operation)
	}
}
