package flac

import (
	"dadp.flactool/task"
	"dadp.flactool/types"
	"dadp.flactool/utils"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var t4VORB_Body *MetaBlockT4VORB
var t4VORB_Body_Lock sync.Mutex

func TaskHandler_T4VORB_GetBody() (*MetaBlockT4VORB, *types.Exception) {
	if t4VORB_Body == nil {
		t4VORB_Body_Lock.Lock()
		defer t4VORB_Body_Lock.Unlock()
		if t4VORB_Body != nil {
			return t4VORB_Body, nil
		}
		var err *types.Exception
		var globalFlac *Flac
		if globalFlac, err = GlobalFlacInit(); err != nil {
			return nil, err
		}

		for blockIndex, block := range globalFlac.GetBlocks() {
			if block.blockType != MetaBlockType_VORBIS_COMMENT {
				continue
			}
			if blockBody, ok := block.GetBody().(*MetaBlockT4VORB); !ok {
				return nil, types.NewException(TMFlac_CanNotAssert_METABLOCKAsSpecificType, map[string]string{
					"index": strconv.Itoa(blockIndex),
					"type":  "MetaBlockT4VORB",
				}, nil)
			} else {
				t4VORB_Body = blockBody
				return t4VORB_Body, nil
			}
		}
		return nil, types.NewException(TMFlac_MetaT4_NotFound, nil, nil)
	}

	return t4VORB_Body, nil
}

func TaskHandler_T4VORB_PrintRefer(args interface{}) (interface{}, *types.Exception) {
	if body, err := TaskHandler_T4VORB_GetBody(); err != nil {
		return nil, err
	} else {
		return body.GetRefer(), nil
	}
}

func TaskHandler_T4VORB_SetRefer(args interface{}) (interface{}, *types.Exception) {
	if args == nil {
		return nil, nil
	} else if body, err := TaskHandler_T4VORB_GetBody(); err != nil {
		return nil, err
	} else if newRefer, ok := args.(string); !ok {
		return nil, types.Exception_Mismatched_Format("Type:string", "Type:"+reflect.TypeOf(args).String())
	} else if newReferProcessed, _, err := task.GlobalArgFilter().FillArgs(newRefer, nil); err != nil {
		return nil, err
	} else {
		body.SetRefer(newReferProcessed)
		return nil, nil
	}
}

func TaskHandler_T4VORB_PrintTags(args interface{}) (interface{}, *types.Exception) {
	if body, err := TaskHandler_T4VORB_GetBody(); err != nil {
		return nil, err
	} else {
		message := types.NewBuffer()
		if commentList, err := body.DumpCommentList(); err != nil {
			return nil, types.NewException(TMFlac_CanNotDump_MetaT4CommentList, nil, err)
		} else {
			for _, commentRecord := range commentList {
				_, _ = message.WriteString(fmt.Sprintf("%s=%s\n", commentRecord[0], commentRecord[1]))
			}
		}

		return strings.TrimSpace(message.String()), nil
	}
}

func TaskHandler_T4VORB_SetTags(args interface{}) (interface{}, *types.Exception) {
	if argList, err := types.InterfaceToStringSlice(args, types.TypeString_Error); err != nil {
		return nil, err
	} else {
		for _, argLine := range argList {
			if argParsed, _, err := task.GlobalArgFilter().FillArgs(argLine, nil); err != nil {
				return nil, err
			} else if err := TaskHandler_T4VORB_SetTagLine(argParsed); err != nil {
				return nil, TaskHandler_T4VORB_SetTagLine(argParsed)
			}
		}
	}
	return nil, nil
}

func TaskHandler_T4VORB_SetTagLine(line string) *types.Exception {
	if body, err := TaskHandler_T4VORB_GetBody(); err != nil {
		return err
	} else if lineSplit := strings.SplitN(line, "=", 2); len(lineSplit) != 2 {
		return types.Exception_Mismatched_Format("(key)=(value)", line)
	} else {
		body.SetComments(strings.TrimSpace(lineSplit[0]), strings.TrimSpace(lineSplit[1]), types.SSLM_Append)
	}
	return nil
}

func TaskHandler_T4VORB_dumpTags(args interface{}) (interface{}, *types.Exception) {
	if args == nil {
		return nil, nil
	} else if tagListPath, ok := args.(string); !ok {
		return nil, types.Exception_Mismatched_Format("Type:string", "Type:"+reflect.TypeOf(args).String())
	} else if tagListPathParsed, _, err := task.GlobalArgFilter().FillArgs(tagListPath, nil); err != nil {
		return nil, err
	} else if tagListContent, err := TaskHandler_T4VORB_PrintTags(nil); err != nil {
		return nil, err
	} else if err := utils.FileWriteString(tagListPathParsed, tagListContent.(string)); err != nil {
		return nil, err
	}

	return nil, nil
}

func TaskHandler_T4VORB_importTags(args interface{}) (interface{}, *types.Exception) {
	if args == nil {
		return nil, nil
	} else if tagListPath, ok := args.(string); !ok {
		return nil, types.Exception_Mismatched_Format("Type:string", "Type:"+reflect.TypeOf(args).String())
	} else if tagListPathParsed, _, err := task.GlobalArgFilter().FillArgs(tagListPath, nil); err != nil {
		return nil, err
	} else if fileContent, err := utils.FileReadBytes(tagListPathParsed); err != nil {
		return nil, err
	} else {
		for _, tagLine := range strings.Split(string(fileContent), "\n") {
			if err := TaskHandler_T4VORB_SetTagLine(tagLine); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func TaskHandler_T4VORB_loadTags(args interface{}) (interface{}, *types.Exception) {
	if args == nil {
		return nil, nil
	} else if body, err := TaskHandler_T4VORB_GetBody(); err != nil {
		return nil, err
	} else if tagListPath, ok := args.(string); !ok {
		return nil, types.Exception_Mismatched_Format("Type:string", "Type:"+reflect.TypeOf(args).String())
	} else if tagListPathParsed, _, err := task.GlobalArgFilter().FillArgs(tagListPath, nil); err != nil {
		return nil, err
	} else if fileContent, err := utils.FileReadBytes(tagListPathParsed); err != nil {
		return nil, err
	} else if strings.TrimSpace(string(fileContent)) == "" {
		body.SetCommentMap(&types.SSListedMap{})
	} else {
		commentMap := &types.SSListedMap{}
		for _, tagLine := range strings.Split(string(fileContent), "\n") {
			if lineSplit := strings.SplitN(tagLine, "=", 2); len(lineSplit) != 2 {
				return nil, types.Exception_Mismatched_Format("(key)=(value)", tagLine)
			} else {
				commentMap.Set(strings.TrimSpace(lineSplit[0]), strings.TrimSpace(lineSplit[1]), types.SSLM_Append)
			}
		}
		body.SetCommentMap(commentMap)
	}

	return nil, nil
}

func TaskHandler_T4VORB_DeleteTags(args interface{}) (interface{}, *types.Exception) {
	if body, err := TaskHandler_T4VORB_GetBody(); err != nil {
		return nil, err
	} else if argList, err := types.InterfaceToStringSlice(args); err != nil {
		return nil, err
	} else {
		for _, arg := range argList {
			body.DeleteComment(arg)
		}
	}
	return nil, nil
}

func TaskHandler_T4VORB_SortTags(args interface{}) (interface{}, *types.Exception) {
	if body, err := TaskHandler_T4VORB_GetBody(); err != nil {
		return nil, err
	} else if sortBy, err := types.InterfaceToStringSlice(args); err != nil {
		return nil, err
	} else {
		body.SortComment(sortBy)
	}

	return nil, nil
}
