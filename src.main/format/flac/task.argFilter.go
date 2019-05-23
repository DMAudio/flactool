package flac

import (
	"dadp.flactool/types"
	"strings"
)

func SplitFilterArg(argRaw string) (string, string, *types.Exception) {
	argSplit := strings.SplitN(argRaw, "->", 2)
	if len(argSplit) != 2 {
		return "", "", types.Mismatched_Format_Exception(
			"筛选条件 -> 操作参数",
			argRaw,
		)
	} else if searchPattern := strings.TrimSpace(argSplit[0]); searchPattern == "" {
		return "", "", types.Mismatched_Format_Exception(
			":筛选条件",
			searchPattern,
		)
	} else if operationArg := strings.TrimSpace(argSplit[1]); operationArg == "" {
		return "", "", types.Mismatched_Format_Exception(
			"操作参数",
			"(空)",
		)
	} else {
		return searchPattern, operationArg, nil
	}
}

func ArgFilter(input string, extraArgs map[string]interface{}) (string, *types.Exception) {
	var err *types.Exception
	var blockFilter, blockTagPath string
	if blockFilter, blockTagPath, err = SplitFilterArg(input); err != nil {
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
		if globalFlac, err = GlobalFlacInit(); err != nil {
			return "", err
		}

		blocks := globalFlac.GetBlocksByIndexSlice(globalFlac.FindBlocks(blockFilter))

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
