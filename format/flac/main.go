package flac

import (
	"gitlab.com/KTGWKenta/DADP.FlacTool/task"
	"gitlab.com/MGEs/Com.Base/types"
	"strconv"
)

func Init() *types.Exception {
	if err := task.GlobalHandler().Register(TaskHandler_T4VORB_Key, TaskHandler_T4VORB); err != nil {
		return err
	}

	if err := task.GlobalHandler().Register(TaskHandler_T6PICT_Key, TaskHandler_T6PICT); err != nil {
		return err
	}

	if err := task.GlobalHandler().Register(TaskHandler_BLOCKS_Key, TaskHandler_BLOCKS); err != nil {
		return err
	}

	if err := task.GlobalArgFiller().Register("flac", ArgFiller); err != nil {
		return err
	}

	return nil
}

func SIMap_GetFlac(env map[string]interface{}) (*Flac, *types.Exception) {
	if fObjRaw, exist := env["flac"]; !exist || fObjRaw == nil {
		return nil, types.NewException(TMFlac_UndefinedObject, nil, nil)
	} else if fObj, ok := fObjRaw.(*Flac); !ok {
		return nil, types.NewException(TMFlac_CanNotAssert_FlacObject, nil, nil)
	} else if !fObj.Initialized() {
		return nil, types.NewException(TMFlac_UninitializedObject, nil, nil)
	} else {
		return fObj, nil
	}
}

func TaskHandler_GetT4VORBBody(fObj *Flac) (*MetaBlockT4VORB, *types.Exception) {
	for blockIndex, block := range fObj.GetBlocks() {
		if block.blockType != MetaBlockType_VORBIS_COMMENT {
			continue
		}
		if blockBody, ok := block.GetBody().(*MetaBlockT4VORB); !ok {
			return nil, types.NewException(TMFlac_CanNotAssert_METABLOCKAsSpecificType, map[string]string{
				"index": strconv.Itoa(blockIndex),
				"type":  "MetaBlockT4VORB",
			}, nil)
		} else {
			return blockBody, nil
		}
	}
	return nil, types.NewException(TMFlac_MetaT4_NotFound, nil, nil)
}

func TaskHandler_Exception_UnsupportedTask(operation string) *types.Exception {
	return types.NewException(task.TMConfig_Unsupported_TaskOperation, map[string]string{
		"operation": operation,
	}, nil)
}

const TaskHandler_T4VORB_Key = MetaBlockTypeStr_VORBIS_COMMENT

func TaskHandler_T4VORB(operation string, env map[string]interface{}, args interface{}) (interface{}, *types.Exception) {
	if fObj, err := SIMap_GetFlac(env); err != nil {
		return nil, err
	} else if body, err := TaskHandler_GetT4VORBBody(fObj); err != nil {
		return nil, err
	} else {
		switch operation {
		case "printRefer":
			return TaskHandler_T4VORB_PrintRefer(fObj, body, args)
		case "setRefer":
			return TaskHandler_T4VORB_SetRefer(fObj, body, args)
		case "printTags":
			return TaskHandler_T4VORB_PrintTags(fObj, body, args)
		case "setTags":
			return TaskHandler_T4VORB_SetTags(fObj, body, args)
		case "loadTags":
			return TaskHandler_T4VORB_loadTags(fObj, body, args)
		case "dumpTags":
			return TaskHandler_T4VORB_dumpTags(fObj, body, args)
		case "importTags":
			return TaskHandler_T4VORB_importTags(fObj, body, args)
		case "deleteTags":
			return TaskHandler_T4VORB_DeleteTags(fObj, body, args)
		case "sortTags":
			return TaskHandler_T4VORB_SortTags(fObj, body, args)
		default:
			return nil, TaskHandler_Exception_UnsupportedTask(operation)
		}
	}
}

const TaskHandler_T6PICT_Key = MetaBlockTypeStr_PICTURE

func TaskHandler_T6PICT(operation string, env map[string]interface{}, args interface{}) (interface{}, *types.Exception) {
	if fObj, err := SIMap_GetFlac(env); err != nil {
		return nil, err
	} else {
		switch operation {
		case "dumpPic":
			return TaskHandler_T6PICT_dumpPic(fObj, args)
		case "setPic":
			return TaskHandler_T6PICT_setPic(fObj, args)
		case "addPic":
			return TaskHandler_T6PICT_addPic(fObj, args)
		case "setPicType":
			return TaskHandler_T6PICT_setPicType(fObj, args)
		case "setDesc":
			return TaskHandler_T6PICT_setDesc(fObj, args)
		case "getDesc":
			return TaskHandler_T6PICT_getDesc(fObj, args)
		default:
			return nil, TaskHandler_Exception_UnsupportedTask(operation)
		}
	}
}

const TaskHandler_BLOCKS_Key = "BLOCKS"

func TaskHandler_BLOCKS(operation string, env map[string]interface{}, args interface{}) (interface{}, *types.Exception) {
	if fObj, err := SIMap_GetFlac(env); err != nil {
		return nil, err
	} else {
		switch operation {
		case "sortBlocks":
			return TaskHandler_MAIN_SortBlocks(fObj, args)
		case "deleteBlocks":
			return TaskHandler_MAIN_DeleteBlocks(fObj, args)
		default:
			return nil, TaskHandler_Exception_UnsupportedTask(operation)
		}
	}
}
