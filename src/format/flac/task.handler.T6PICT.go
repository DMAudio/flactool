package flac

import (
	"p20190417/task"
	"p20190417/types"
	"p20190417/utils"
	"strconv"
	"strings"
)

func TaskHandler_T6PICT_Each(args interface{}, processer func(*MetaBlock, *MetaBlockT6PICT, string) *types.Exception) *types.Exception {
	var argsRaw []string
	var err *types.Exception
	if argsRaw, err = types.InterfaceToStringSlice(args); err != nil {
		return err
	}

	var globalFlac *Flac
	if globalFlac, err = GlobalFlacInit(); err != nil {
		return err
	}

	PatternPathPairs := make([][2]string, 0)
	for _, argRaw := range argsRaw {
		if searchPattern, filePath, err := SplitFilterArg(argRaw); err != nil {
			return err
		} else if !strings.HasPrefix(searchPattern, MetaBlockTypeStr_PICTURE) {
			return types.Mismatched_Format_Exception(
				MetaBlockTypeStr_PICTURE+"[:筛选条件]",
				searchPattern,
			)
		} else {
			PatternPathPairs = append(PatternPathPairs, [2]string{searchPattern, filePath})
		}
	}

	for _, pair := range PatternPathPairs {
		pattern, pArgs := pair[0], pair[1]
		for _, blockIndex := range globalFlac.FindBlocks(pattern) {
			block := globalFlac.GetBlockByIndex(blockIndex)
			if blockBody, ok := block.GetBody().(*MetaBlockT6PICT); !ok {
				return types.NewException(TMFlac_CanNotAssert_METABLOCKAsSpecificType, map[string]string{
					"index": strconv.Itoa(blockIndex),
					"type":  "MetaBlockT6PICT",
				}, nil)
			} else if err := processer(block, blockBody, pArgs); err != nil {
				return err
			}
		}
	}

	return nil
}

func TaskHandler_T6PICT_dumpPic(args interface{}) (interface{}, *types.Exception) {
	return nil, TaskHandler_T6PICT_Each(args, func(block *MetaBlock, body *MetaBlockT6PICT, path string) *types.Exception {
		if pathProcessed, _, err := task.GlobalArgFilter().FillArgs(path, map[string]map[string]interface{}{
			"flac": {"this": block},
		}); err != nil {
			return types.NewException(TMFlac_Task_CanNotParse_FileAbsolutePath, nil, err)
		} else {
			if err := utils.FileWriteBytes(pathProcessed, body.GetPicRawData()); err != nil {
				return err
			}
		}
		return nil
	})
}

func TaskHandler_T6PICT_setPicType(args interface{}) (interface{}, *types.Exception) {
	return nil, TaskHandler_T6PICT_Each(args, func(block *MetaBlock, body *MetaBlockT6PICT, picType string) *types.Exception {
		if picTypeParsed, err := strconv.ParseUint(picType, 10, 32); err != nil {
			return types.NewException(TMFlac_CanNotRead_MetaT6Type, nil, err)
		} else {
			body.SetPicType(uint32(picTypeParsed))
		}
		return nil
	})
}

