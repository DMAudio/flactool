package flac

import (
	"fmt"
	"gitlab.com/KTGWKenta/DADP.FlacTool/task"
	base "gitlab.com/MGEs/Com.Base"
	"gitlab.com/MGEs/Com.Base/types"
	"reflect"
	"strconv"
	"strings"
)

func TaskHandler_T6PICT_Each(
	flac *Flac,
	args interface{},
	argChecker func([]string) *types.Exception,
	processor func(*MetaBlock, *MetaBlockT6PICT, []string) *types.Exception,
) *types.Exception {
	var err *types.Exception
	var argSlice []interface{}
	if argSlice, err = types.InterfaceToInterfaceSlice(args); err != nil {
		return err
	}

	type BlockId_ArgVal_Pair struct {
		errTag   string
		bIdSlice []int
		argSlice []string
	}

	BlockId_ArgVal_Pairs := make([]BlockId_ArgVal_Pair, 0)

	argValParseExceptions := map[string]*types.Exception{}
	for argIndex, argVal := range argSlice {
		errTag := "Sa" + strconv.Itoa(argIndex)
		if argValSlice, err := types.InterfaceToStringSlice(argVal, types.TypeString_Error); err != nil {
			argValParseExceptions[errTag] = err
		} else if len(argValSlice) == 0 {
			argValParseExceptions[errTag] = types.Exception_Mismatched_Format(
				"[BlockFilter, ArgVal1, ...]",
				strings.Join(argValSlice, ","),
			)
		} else if blockIdSlice, err := flac.FindBlocks(argValSlice[0]); err != nil {
			argValParseExceptions[errTag] = err
		} else if err := argChecker(argValSlice[1:]); err != nil {
			argValParseExceptions[errTag] = err
		} else {
			BlockId_ArgVal_Pairs = append(BlockId_ArgVal_Pairs, BlockId_ArgVal_Pair{
				errTag:   errTag,
				bIdSlice: blockIdSlice,
				argSlice: argValSlice[1:],
			})
		}
	}

	if len(argValParseExceptions) > 0 {
		return types.NewException(task.TMTask_UnableToParse_SubArg, nil, argValParseExceptions)
	}

	argValExecuteExceptions := map[string]*types.Exception{}
	for _, pairVal := range BlockId_ArgVal_Pairs {
		for _, blockIndex := range pairVal.bIdSlice {
			block := flac.GetBlockByIndex(blockIndex)
			if blockBody, ok := block.GetBody().(*MetaBlockT6PICT); !ok {
				argValExecuteExceptions[pairVal.errTag] = types.NewException(
					TMFlac_CanNotAssert_METABLOCKAsSpecificType, map[string]string{
						"index": strconv.Itoa(blockIndex),
						"type":  "MetaBlockT6PICT",
					}, nil)
			} else if err := processor(block, blockBody, pairVal.argSlice); err != nil {
				argValExecuteExceptions[pairVal.errTag] = err
			}
		}
	}
	if len(argValExecuteExceptions) > 0 {
		return types.NewException(task.TMTask_UnableToExecute_SubTask, nil, argValExecuteExceptions)
	} else {
		return nil
	}
}

func TaskHandler_T6PICT_SingleArgChecker(args []string) *types.Exception {
	if len(args) != 1 {
		return types.Exception_Mismatched_ArgumentList(
			"[1]string{Path}",
			fmt.Sprintf("[%d]string{%s}", len(args), strings.Join(args, ", ")),
		)
	}
	return nil
}

func TaskHandler_T6PICT_dumpPic(flac *Flac, args interface{}) (interface{}, *types.Exception) {
	handler := func(block *MetaBlock, body *MetaBlockT6PICT, args []string) *types.Exception {
		if pathProcessed, _, err := task.GlobalArgFiller().FillArgs(args[0], map[string]map[string]interface{}{
			"flac": {"flac": flac, "this": block},
		}); err != nil {
			return err
		} else {
			if err := base.FileWriteBytes(pathProcessed, body.GetPicRawData()); err != nil {
				return err
			}
		}
		return nil
	}
	return nil, TaskHandler_T6PICT_Each(flac, args, TaskHandler_T6PICT_SingleArgChecker, handler)
}

func TaskHandler_T6PICT_setPic(flac *Flac, args interface{}) (interface{}, *types.Exception) {
	handler := func(block *MetaBlock, body *MetaBlockT6PICT, args []string) *types.Exception {
		if pathProcessed, _, err := task.GlobalArgFiller().FillArgs(args[0], map[string]map[string]interface{}{
			"flac": {"flac": flac, "this": block},
		}); err != nil {
			return err
		} else {
			if err := body.ParsePictureFile(pathProcessed); err != nil {
				return err
			}
		}
		return nil
	}
	return nil, TaskHandler_T6PICT_Each(flac, args, TaskHandler_T6PICT_SingleArgChecker, handler)
}

func TaskHandler_T6PICT_setPicType(flac *Flac, args interface{}) (interface{}, *types.Exception) {
	handler := func(block *MetaBlock, body *MetaBlockT6PICT, args []string) *types.Exception {
		if picTypeProcessed, _, err := task.GlobalArgFiller().FillArgs(args[0], map[string]map[string]interface{}{
			"flac": {"flac": flac, "this": block},
		}); err != nil {
			return err
		} else if picTypeParsed, err := strconv.ParseUint(picTypeProcessed, 10, 32); err != nil {
			return types.NewException(TMFlac_CanNotRead_MetaT6Type, nil, err)
		} else {
			body.SetPicType(uint32(picTypeParsed))
		}
		return nil
	}
	return nil, TaskHandler_T6PICT_Each(flac, args, TaskHandler_T6PICT_SingleArgChecker, handler)
}

func TaskHandler_T6PICT_setDesc(flac *Flac, args interface{}) (interface{}, *types.Exception) {
	handler := func(block *MetaBlock, body *MetaBlockT6PICT, args []string) *types.Exception {
		if descProcessed, _, err := task.GlobalArgFiller().FillArgs(args[0], map[string]map[string]interface{}{
			"flac": {"flac": flac, "this": block},
		}); err != nil {
			return err
		} else {
			body.SetPicDesc(descProcessed)
		}
		return nil
	}
	return nil, TaskHandler_T6PICT_Each(flac, args, TaskHandler_T6PICT_SingleArgChecker, handler)
}

func TaskHandler_T6PICT_getDesc(flac *Flac, args interface{}) (interface{}, *types.Exception) {
	if argParsed, ok := args.(string); !ok {
		return nil, types.Exception_Mismatched_Format("Type:string", reflect.TypeOf(args).String())
	} else if pattern, _, err := task.GlobalArgFiller().FillArgs(argParsed, ToArgFillerParameter(flac)); err != nil {
		return nil, err
	} else {
		results := make([]string, 0)
		var blockIndexSlice []int
		blockIndexSlice, err = flac.FindBlocks(pattern)
		if err != nil {
			return nil, err
		}
		for _, blockIndex := range blockIndexSlice {
			block := flac.GetBlockByIndex(blockIndex)
			if blockBody, ok := block.GetBody().(*MetaBlockT6PICT); !ok {
				return nil, types.NewException(TMFlac_CanNotAssert_METABLOCKAsSpecificType, map[string]string{
					"index": strconv.Itoa(blockIndex),
					"type":  "MetaBlockT6PICT",
				}, nil)
			} else {
				results = append(results, blockBody.GetPicDesc())
			}
		}
		return strings.Join(results, "\n"), nil
	}
}

func TaskHandler_T6PICT_addPic(flac *Flac, args interface{}) (interface{}, *types.Exception) {
	if argsParsed, err := types.InterfaceToStringMap(args); err != nil {
		return nil, err
	} else {
		var picType uint32
		var picPath string
		var picDesc string
		var keyExist bool

		TMFLAC_LackOfMustInfo := types.NewMask("LackOfMustInfo", "缺少必要信息：{{info}}")

		if picTypeRaw, keyExist := argsParsed["type"]; !keyExist {
			return nil, types.NewException(TMFLAC_LackOfMustInfo, map[string]string{
				"info": "图片类型(type)",
			}, nil)
		} else if picTypeParsed, err := strconv.ParseUint(picTypeRaw, 10, 32); err != nil {
			return nil, types.NewException(TMFlac_CanNotRead_MetaT6Type, nil, err)
		} else {
			picType = uint32(picTypeParsed)
		}

		if picPath, keyExist = argsParsed["path"]; !keyExist {
			return nil, types.NewException(TMFLAC_LackOfMustInfo, map[string]string{
				"info": "图片路径(path)",
			}, nil)
		}

		blockBody := &MetaBlockT6PICT{}
		blockBody.SetPicType(picType)
		if err := blockBody.ParsePictureFile(picPath); err != nil {
			return nil, err
		}

		if picDesc, keyExist = argsParsed["desc"]; keyExist {

			blockBody.SetPicDesc(picDesc)
		}

		flac.AppendBlock(NewMetaBlock(blockBody))
	}
	return nil, nil
}
