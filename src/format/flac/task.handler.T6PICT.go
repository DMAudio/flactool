package flac

import (
	"p20190417/types"
)

func TaskHandler_T6PICT_dumpPic(args interface{}) (interface{}, *types.Exception) {
	if _, err := types.InterfaceToStringSlice(args); err != nil {
		return nil, err
	} else {
		//TypePathMap := map[string]string{}

		//for _, argRaw := range argList {
		//	argRaw = strings.TrimSpace(argRaw)
		//	argSplit := strings.SplitN(argRaw, "->", 2)
		//	if len(argSplit) == 2 {
		//		if _, err := strconv.Atoi(argSplit[0]); err != nil {
		//			return nil, types.NewException(TMFlac_CanNotRead_MetaT6Type, nil, err)
		//		}
		//		// 图片类型 -> 图片路径
		//		TypePathMap[argSplit[0]] = argSplit[1]
		//	} else {
		//		TypePathMap["*"] = argRaw
		//	}
		//}

		var globalFlac *Flac
		if globalFlac = GlobalFlac(); globalFlac == nil || !globalFlac.Initialized() {
			return nil, types.NewException(TMFlac_UninitializedObject, nil, nil)
		}


		//for blockIndex, blockRaw := range globalFlac.MetaBlocks {
		//	if blockRaw.GetType() != MetaBlockType_PICTURE {
		//		continue
		//	}
		//	if _, ok := blockRaw.GetBody().(*MetaBlockT6PICT); !ok {
		//		return nil, types.NewException(TMFlac_CanNotAssert_METABLOCKAsSpecificType, map[string]string{
		//			"index": strconv.Itoa(blockIndex),
		//			"type":  "MetaBlockT6PICT",
		//		}, nil)
		//	} else {
		//
		//	}
		//}
	}

	return nil, nil
}

func TaskHandler_T6PICT_setPic(args interface{}) (interface{}, *types.Exception) {

	return nil, nil
}
