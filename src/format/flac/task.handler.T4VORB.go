package flac

import (
	"fmt"
	"io/ioutil"
	"os"
	"p20190417/types"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
)

var t4VORB_Body *MetaBlockT4VORB
var t4VORB_Body_Lock sync.Mutex

func TaskHandler_T4VORB_GetBody() (*MetaBlockT4VORB, *types.Exception) {
	if t4VORB_Body == nil {
		t4VORB_Body_Lock.Lock()
		defer t4VORB_Body_Lock.Unlock()
		if t4VORB_Body == nil {
			globalFlac := GlobalFlac()
			if globalFlac == nil || !globalFlac.Initialized() {
				return nil, types.NewException(TMFlac_UninitializedObject, nil, nil)
			}

			for _, block := range globalFlac.MetaBlocks {
				if block.Type != MetaBlockType_VORBIS_COMMENT {
					continue
				}
				if blockBody, ok := block.Body.(*MetaBlockT4VORB); ok {
					t4VORB_Body = blockBody
					return t4VORB_Body, nil
				}
			}
			return nil, types.NewException(TMFlac_MetaT4_NotFound, nil, nil)
		}
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
	if body, err := TaskHandler_T4VORB_GetBody(); err != nil {
		return nil, err
	} else if newRefer, ok := args.(string); !ok {
		return nil, types.Mismatched_Format_Exception("string",
			reflect.TypeOf(args).String(),
		)
	} else {
		body.SetRefer(newRefer)
		return nil, nil
	}
}

func TaskHandler_T4VORB_PrintTags(args interface{}) (interface{}, *types.Exception) {
	if body, err := TaskHandler_T4VORB_GetBody(); err != nil {
		return nil, err
	} else {
		message := types.NewBuffer()
		if commentList, err := body.Comments().DumpList(); err != nil {
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
	if argList, err := types.InterfaceToStringSlice(args); err != nil {
		return nil, err
	} else {
		for _, argParsed := range argList {
			if err := TaskHandler_T4VORB_SetTagLine(argParsed); err != nil {
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
		return types.Mismatched_Format_Exception("(key)=(value)",
			line,
		)
	} else {
		body.Comments().Set(lineSplit[0], lineSplit[1], types.SSLM_Append)
	}
	return nil
}

func TaskHandler_T4VORB_dumpTags(args interface{}) (interface{}, *types.Exception) {
	if tagListPath, ok := args.(string); !ok {
		return nil, types.Mismatched_Format_Exception("string",
			reflect.TypeOf(args).String(),
		)
	} else if tagListAbsPath, err := filepath.Abs(strings.TrimSpace(tagListPath)); err != nil {
		return nil, types.NewException(TMFlac_Task_CanNotParse_FileAbsolutePath, nil, nil)
	} else if tagListContent, err := TaskHandler_T4VORB_PrintTags(nil); err != nil {
		return nil, err
	} else if fileObj, err := os.OpenFile(tagListAbsPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return nil, types.NewException(TMFlac_Task_CanNotOpen_File, nil, nil)
	} else {
		defer func() {
			_ = fileObj.Close()
		}()
		if _, err := fileObj.WriteString(tagListContent.(string)); err != nil {
			return nil, types.NewException(TMFlac_Task_CanNotWriteTo_File, nil, nil)
		}
	}

	return nil, nil
}

func TaskHandler_T4VORB_importTags(args interface{}) (interface{}, *types.Exception) {
	if tagListPath, ok := args.(string); !ok {
		return nil, types.Mismatched_Format_Exception("string",
			reflect.TypeOf(args).String(),
		)
	} else if tagListAbsPath, err := filepath.Abs(strings.TrimSpace(tagListPath)); err != nil {
		return nil, types.NewException(TMFlac_Task_CanNotParse_FileAbsolutePath, nil, nil)
	} else if fileObj, err := os.OpenFile(tagListAbsPath, os.O_RDONLY, 0644); err != nil {
		return nil, types.NewException(TMFlac_Task_CanNotOpen_File, nil, nil)
	} else {
		defer func() {
			_ = fileObj.Close()
		}()
		if tagListContents, err := ioutil.ReadAll(fileObj); err != nil {
			return nil, types.NewException(TMFlac_Task_CanNotReadFrom_File, nil, nil)
		} else {
			for _, tagLine := range strings.Split(string(tagListContents), "\n") {
				if err := TaskHandler_T4VORB_SetTagLine(tagLine); err != nil {
					return nil, err
				}
			}
		}
	}
	return nil, nil
}

func TaskHandler_T4VORB_DeleteTags(args interface{}) (interface{}, *types.Exception) {
	if body, err := TaskHandler_T4VORB_GetBody(); err != nil {
		return nil, err
	} else if argList, ok := args.([]interface{}); !ok {
		return nil, types.Mismatched_Format_Exception("Kind:Slice",
			"Kind:"+reflect.TypeOf(args).Kind().String(),
		)
	} else {
		for _, argRaw := range argList {
			if argParsed, ok := argRaw.(string); !ok {
				return nil, types.Mismatched_Format_Exception("string",
					reflect.TypeOf(argRaw).String(),
				)
			} else {
				body.Comments().Delete(argParsed)
			}
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
		body.Comments().Sort(sortBy)
	}

	return nil, nil
}
