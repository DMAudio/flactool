package flac

import (
	"p20190417/types"
	"strings"
)

func TaskHandler_MAIN_SortBlocks(args interface{}) (interface{}, *types.Exception) {
	globalFlac := GlobalFlac()
	if globalFlac == nil || !globalFlac.Initialized() {
		return nil, types.NewException(TMFlac_UninitializedObject, nil, nil)
	}

	var sortItems []string
	if sortItemsParsed, err := types.InterfaceToStringSlice(args); err != nil {
		return nil, err
	} else {
		sortItems = sortItemsParsed
	}

	if len(sortItems) == 0 || sortItems[0] != MetaBlockTypeStr_STREAMINFO {
		sortItems = append([]string{MetaBlockTypeStr_STREAMINFO}, sortItems...)
	}

	metaBlocks := make([]*MetaBlock, 0)
	metaBlocksIndexHandled := make([]bool, len(globalFlac.MetaBlocks))
	for _, sItem := range sortItems {
		VariadicBlocksAmount, sortItemTrimmed := types.StringRepresentsVariadicType(sItem)

		extendTypeRequired := strings.Index(sItem, ":") != -1
		for blockIndex, block := range globalFlac.MetaBlocks {
			if metaBlocksIndexHandled[blockIndex] {
				continue
			}
			if block.GetTypeStr(extendTypeRequired) == sortItemTrimmed || sortItemTrimmed == "" {
				metaBlocks = append(metaBlocks, block)
				metaBlocksIndexHandled[blockIndex] = true
				if !VariadicBlocksAmount {
					break
				}
			}
		}
	}

	globalFlac.MetaBlocks = metaBlocks

	return nil, nil
}
//
//func TaskHandler_MAIN_DeleteBlocks(args interface{}) (interface{}, *types.Exception) {
//	globalFlac := GlobalFlac()
//	if globalFlac == nil || !globalFlac.Initialized() {
//		return nil, types.NewException(TMFlac_UninitializedObject, nil, nil)
//	}
//
//	var deleteList []string
//	if deleteItemsParsed, err := util.InterfaceToStringSlice(args); err != nil {
//		return nil, err
//	} else {
//		deleteList = deleteItemsParsed
//	}
//
//	metaBlocks := make([]*MetaBlock, 0)
//
//	for blockIndex, block := range globalFlac.MetaBlocks {
//	//	for _, dItem := range deleteList {
//	//
//	//	}
//	}
//
//	return nil,nil
//}