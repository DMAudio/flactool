package flac

import (
	"dadp.flactool/types"
	"strings"
)

func TaskHandler_MAIN_SortBlocks(args interface{}) (interface{}, *types.Exception) {
	var err *types.Exception
	var globalFlac *Flac
	if globalFlac, err = GlobalFlacInit(); err != nil {
		return "", err
	}

	var searchPatterns []string
	if sortItemsParsed, err := types.InterfaceToStringSlice(args, types.TypeString_Error); err != nil {
		return nil, err
	} else {
		searchPatterns = sortItemsParsed
	}

	if len(searchPatterns) == 0 || searchPatterns[0] != MetaBlockTypeStr_STREAMINFO {
		searchPatterns = append([]string{MetaBlockTypeStr_STREAMINFO}, searchPatterns...)
	}

	blocks := globalFlac.GetBlocks()
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
		if blockIndexes, err := globalFlac.FindBlocks(pattern); err != nil {
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

	globalFlac.SetBlocks(newBlockList)

	return nil, nil
}

func TaskHandler_MAIN_DeleteBlocks(args interface{}) (interface{}, *types.Exception) {
	var err *types.Exception
	var globalFlac *Flac
	if globalFlac, err = GlobalFlacInit(); err != nil {
		return "", err
	}

	var searchPatterns []string
	if sortItemsParsed, err := types.InterfaceToStringSlice(args, types.TypeString_Error); err != nil {
		return nil, err
	} else {
		searchPatterns = sortItemsParsed
	}

	blocks := globalFlac.GetBlocks()
	blockUnMatched := make([]int, len(blocks), len(blocks))
	for i := 0; i < cap(blockUnMatched); i++ {
		blockUnMatched[i] = i
	}
	for _, pattern := range searchPatterns {
		pattern = strings.TrimSpace(pattern)
		if blockIndexes, err := globalFlac.FindBlocks(pattern); err != nil {
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

	globalFlac.SetBlocks(newBlockList)

	return nil, nil
}
