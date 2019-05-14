package flac

import (
	"p20190417/types"
)

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
		return nil, nil
	}
}

const TaskHandler_T6PICT_Key = MetaBlockTypeStr_PICTURE

func TaskHandler_T6PICT(operation string, args interface{}) (interface{}, *types.Exception) {
	switch operation {
	case "dumpPic":
		return TaskHandler_T6PICT_dumpPic(args)
	case "setPic":
		return TaskHandler_T6PICT_setPic(args)
	default:
		return nil, nil
	}
}

const TaskHandler_BLOCKS_Key = "BLOCKS"

func TaskHandler_BLOCKS(operation string, args interface{}) (interface{}, *types.Exception) {
	switch operation {
	case "sortBlocks":
		return TaskHandler_MAIN_SortBlocks(args)
	//case "deleteBlocks":
	//	return TaskHandler_MAIN_DeleteBlocks(args)
	default:
		return nil, nil
	}
}

