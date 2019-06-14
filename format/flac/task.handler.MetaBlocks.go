package flac

import (
	"gitlab.com/MGEs/Com.Base/types"
	"strings"
)

func TaskHandler_MAIN_SortBlocks(flac *Flac, args interface{}) (interface{}, *types.Exception) {
	var searchPatterns []string
	if sortItemsParsed, err := types.InterfaceToStringSlice(args, types.TypeString_Error); err != nil {
		return nil, err
	} else {
		searchPatterns = sortItemsParsed
	}

	if len(searchPatterns) == 0 || searchPatterns[0] != MetaBlockTypeStr_STREAMINFO {
		searchPatterns = append([]string{MetaBlockTypeStr_STREAMINFO}, searchPatterns...)
	}

	blocks := flac.GetBlocks()
	blockUnMatched := make([]int, len(blocks), len(blocks))
	blockMatched := make([]int, 0, len(blockUnMatched))
	preserveUnMatchedBlocks := -1
	for i := 0; i < cap(blockUnMatched); i++ {
		blockUnMatched[i] = i
	}
	for _, pattern := range searchPatterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "..." {
			preserveUnMatchedBlocks = len(blockMatched) - 1
		}
		if blockIndexes, err := flac.FindBlocks(pattern); err != nil {
			return nil, err
		} else if blockIndexes != nil && len(blockIndexes) != 0 {
			for _, blockIndex := range blockIndexes {
				if types.IListFindElement(blockMatched, blockIndex) > -1 {
					continue
				}
				blockMatched = append(blockMatched, blockIndex)
				blockUnMatched = types.IListDeleteByElement(blockUnMatched, blockIndex)
			}
		}
	}

	if preserveUnMatchedBlocks > -1 {
		blockMatched = types.IListInsertAfter(blockMatched, preserveUnMatchedBlocks, blockUnMatched...)
	}

	newBlockList := make([]*MetaBlock, len(blockMatched))
	for i, blockIndex := range blockMatched {
		newBlockList[i] = blocks[blockIndex]
	}

	flac.SetBlocks(newBlockList)

	return nil, nil
}

func TaskHandler_MAIN_DeleteBlocks(flac *Flac, args interface{}) (interface{}, *types.Exception) {
	var searchPatterns []string
	if sortItemsParsed, err := types.InterfaceToStringSlice(args, types.TypeString_Error); err != nil {
		return nil, err
	} else {
		searchPatterns = sortItemsParsed
	}

	blocks := flac.GetBlocks()
	blockUnMatched := make([]int, len(blocks), len(blocks))
	for i := 0; i < cap(blockUnMatched); i++ {
		blockUnMatched[i] = i
	}
	for _, pattern := range searchPatterns {
		pattern = strings.TrimSpace(pattern)
		if blockIndexes, err := flac.FindBlocks(pattern); err != nil {
			return nil, err
		} else if blockIndexes != nil && len(blockIndexes) != 0 {
			for _, blockIndex := range blockIndexes {
				blockUnMatched = types.IListDeleteByElement(blockUnMatched, blockIndex)
			}
		}
	}

	newBlockList := make([]*MetaBlock, len(blockUnMatched))
	for i, blockIndex := range blockUnMatched {
		newBlockList[i] = blocks[blockIndex]
	}

	flac.SetBlocks(newBlockList)

	return nil, nil
}
