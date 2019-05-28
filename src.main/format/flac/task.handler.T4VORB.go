package flac

import (
	"dadp.flactool/task"
	"dadp.flactool/types"
	"dadp.flactool/utils"
	"fmt"
	"reflect"
	"strings"
)

func TaskHandler_T4VORB_PrintRefer(flac *Flac, T4Body *MetaBlockT4VORB, args interface{}) (interface{}, *types.Exception) {
	return T4Body.GetRefer(), nil
}

func TaskHandler_T4VORB_SetRefer(flac *Flac, T4Body *MetaBlockT4VORB, args interface{}) (interface{}, *types.Exception) {
	if args == nil {
		return nil, nil
	} else if newRefer, ok := args.(string); !ok {
		return nil, types.Exception_Mismatched_Format("Type:string", "Type:"+reflect.TypeOf(args).String())
	} else if newReferProcessed, _, err := task.GlobalArgFiller().FillArgs(newRefer, ToArgFillerParameter(flac)); err != nil {
		return nil, err
	} else {
		T4Body.SetRefer(newReferProcessed)
		return nil, nil
	}
}

func TaskHandler_T4VORB_PrintTags(flac *Flac, T4Body *MetaBlockT4VORB, args interface{}) (interface{}, *types.Exception) {
	message := types.NewBuffer()
	if commentList, err := T4Body.DumpCommentList(); err != nil {
		return nil, types.NewException(TMFlac_CanNotDump_MetaT4CommentList, nil, err)
	} else {
		for _, commentRecord := range commentList {
			_, _ = message.WriteString(fmt.Sprintf("%s=%s\n", commentRecord[0], commentRecord[1]))
		}
	}

	return strings.TrimSpace(message.String()), nil
}

func TaskHandler_T4VORB_SetTags(flac *Flac, T4Body *MetaBlockT4VORB, args interface{}) (interface{}, *types.Exception) {
	if argList, err := types.InterfaceToStringSlice(args, types.TypeString_Error); err != nil {
		return nil, err
	} else {
		for _, argLine := range argList {
			if argParsed, _, err := task.GlobalArgFiller().FillArgs(argLine, ToArgFillerParameter(flac)); err != nil {
				return nil, err
			} else if err := TaskHandler_T4VORB_SetTagLine(flac, T4Body, argParsed); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func TaskHandler_T4VORB_SetTagLine(flac *Flac, T4Body *MetaBlockT4VORB, line string) *types.Exception {
	if lineSplit := strings.SplitN(line, "=", 2); len(lineSplit) != 2 {
		return types.Exception_Mismatched_Format("(key)=(value)", line)
	} else {
		T4Body.SetComments(strings.TrimSpace(lineSplit[0]), strings.TrimSpace(lineSplit[1]), types.SSLM_Append)
	}
	return nil
}

func TaskHandler_T4VORB_dumpTags(flac *Flac, T4Body *MetaBlockT4VORB, args interface{}) (interface{}, *types.Exception) {
	if args == nil {
		return nil, nil
	} else if tagListPath, ok := args.(string); !ok {
		return nil, types.Exception_Mismatched_Format("Type:string", "Type:"+reflect.TypeOf(args).String())
	} else if tagListPathParsed, _, err := task.GlobalArgFiller().FillArgs(tagListPath, ToArgFillerParameter(flac)); err != nil {
		return nil, err
	} else if tagListContent, err := TaskHandler_T4VORB_PrintTags(flac, T4Body, nil); err != nil {
		return nil, err
	} else if err := utils.FileWriteString(tagListPathParsed, tagListContent.(string)); err != nil {
		return nil, err
	}

	return nil, nil
}

func TaskHandler_T4VORB_importTags(flac *Flac, T4Body *MetaBlockT4VORB, args interface{}) (interface{}, *types.Exception) {
	if args == nil {
		return nil, nil
	} else if tagListPath, ok := args.(string); !ok {
		return nil, types.Exception_Mismatched_Format("Type:string", "Type:"+reflect.TypeOf(args).String())
	} else if tagListPathParsed, _, err := task.GlobalArgFiller().FillArgs(tagListPath, ToArgFillerParameter(flac)); err != nil {
		return nil, err
	} else if fileContent, err := utils.FileReadBytes(tagListPathParsed); err != nil {
		return nil, err
	} else {
		for _, tagLine := range strings.Split(string(fileContent), "\n") {
			if err := TaskHandler_T4VORB_SetTagLine(flac, T4Body, tagLine); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func TaskHandler_T4VORB_loadTags(flac *Flac, T4Body *MetaBlockT4VORB, args interface{}) (interface{}, *types.Exception) {
	if args == nil {
		return nil, nil
	} else if tagListPath, ok := args.(string); !ok {
		return nil, types.Exception_Mismatched_Format("Type:string", "Type:"+reflect.TypeOf(args).String())
	} else if tagListPathParsed, _, err := task.GlobalArgFiller().FillArgs(tagListPath, ToArgFillerParameter(flac)); err != nil {
		return nil, err
	} else if fileContent, err := utils.FileReadBytes(tagListPathParsed); err != nil {
		return nil, err
	} else if strings.TrimSpace(string(fileContent)) == "" {
		T4Body.SetCommentMap(&types.SSListedMap{})
	} else {
		commentMap := &types.SSListedMap{}
		for _, tagLine := range strings.Split(string(fileContent), "\n") {
			if lineSplit := strings.SplitN(tagLine, "=", 2); len(lineSplit) != 2 {
				return nil, types.Exception_Mismatched_Format("(key)=(value)", tagLine)
			} else {
				commentMap.Set(strings.TrimSpace(lineSplit[0]), strings.TrimSpace(lineSplit[1]), types.SSLM_Append)
			}
		}
		T4Body.SetCommentMap(commentMap)
	}

	return nil, nil
}

func TaskHandler_T4VORB_DeleteTags(flac *Flac, T4Body *MetaBlockT4VORB, args interface{}) (interface{}, *types.Exception) {
	if argList, err := types.InterfaceToStringSlice(args, types.TypeString_Split); err != nil {
		return nil, err
	} else {
		for _, arg := range argList {
			T4Body.DeleteComment(arg)
		}
	}
	return nil, nil
}

func TaskHandler_T4VORB_SortTags(flac *Flac, T4Body *MetaBlockT4VORB, args interface{}) (interface{}, *types.Exception) {
	if sortBy, err := types.InterfaceToStringSlice(args, types.TypeString_Split); err != nil {
		return nil, err
	} else {
		T4Body.SortComment(sortBy)
	}

	return nil, nil
}
