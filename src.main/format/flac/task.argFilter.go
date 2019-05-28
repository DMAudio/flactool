package flac

import (
	"dadp.flactool/types"
	"strings"
)

func SplitFillerArg(argRaw string) (string, string, *types.Exception) {
	argSplit := strings.SplitN(argRaw, "->", 2)
	if len(argSplit) != 2 {
		return "", "", types.Exception_Mismatched_Format("FilterPattern -> OperationArg", argRaw)
	} else if filterPattern := strings.TrimSpace(argSplit[0]); filterPattern == "" {
		return "", "", types.Exception_Mismatched_Format("FilterPattern", "(Empty)")
	} else if operationArg := strings.TrimSpace(argSplit[1]); operationArg == "" {
		return "", "", types.Exception_Mismatched_Format("OperationArg", "(Empty)")
	} else {
		return filterPattern, operationArg, nil
	}
}

func ArgFiller(input string, extraArgs map[string]interface{}) (string, *types.Exception) {
	var err *types.Exception
	var blockFilter, blockTagPath string
	if blockFilter, blockTagPath, err = SplitFillerArg(input); err != nil {
		return "", err
	}

	var block *MetaBlock
	if blockFilter == "this" {
		if thisBlock, exist := extraArgs["this"]; !exist {
			return "", types.NewException(TMFlac_Arg_CanNotParseThisBlock, nil, nil)
		} else if thisBlockParsed, ok := thisBlock.(*MetaBlock); !ok {
			return "", types.NewException(TMFlac_Arg_CanNotParseThisBlock, nil, nil)
		} else {
			block = thisBlockParsed
		}
	} else {
		var globalFlac *Flac
		var blocks []*MetaBlock

		if globalFlac, err = SIMap_GetFlac(extraArgs); err != nil {
			return "", err
		} else if blockIndexSlice, err := globalFlac.FindBlocks(blockFilter); err != nil {
			return "", err
		} else {
			blocks = globalFlac.GetBlocksByIndexSlice(blockIndexSlice)
		}

		if len(blocks) == 0 {
			return "", types.NewException(TMFlac_Arg_CanNotFind_Block, map[string]string{
				"pattern": blockFilter,
			}, nil)
		} else if len(blocks) > 1 {
			return "", types.NewException(TMFlac_Arg_VaguePattern, map[string]string{
				"pattern": blockFilter,
			}, nil)
		}

		block = blocks[0]
	}

	return block.GetTags().Get(blockTagPath), nil

}

func ToArgFillerParameter(flac *Flac) map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		"flac": {"flac": flac},
	}
}
