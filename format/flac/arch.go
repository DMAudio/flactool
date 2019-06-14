package flac

import (
	"bytes"
	base "gitlab.com/MGEs/Com.Base"
	"gitlab.com/MGEs/Com.Base/types"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Flac struct {
	//ID3v2      *PrependID3v2
	blocks []*MetaBlock
	frames []byte
}

var flacSignature = []byte("fLaC")

var id3Signature = []byte("ID3")

func (fObj *Flac) Parse(br *types.BinaryReader) *types.Exception {
	var err *types.Exception
	var marker = make([]byte, 4)
	if marker, err = br.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_FileSignature, nil, err)
	}

	//Parse Prepend ID3V2
	if bytes.Equal(marker[:3], id3Signature) {
		blockPreId3 := PrependID3v2{}
		if err := blockPreId3.Parse(br); err != nil {
			return types.NewException(TMFlac_CanNotParse_PrependID3V2Block, nil, err)
		} // else {
		//	fObj.ID3v2 = &blockPreId3
		//}

		if marker, err = br.ReadBytes(4); err != nil {
			return types.NewException(TMFlac_CanNotRead_FileSignature, nil, err)
		}
	}

	//Verify "fLaC" signature
	if !bytes.Equal(marker, flacSignature) {
		return types.NewException(TMFlac_Incorrect_FileSignature, map[string]string{
			"expected": string(flacSignature),
			"got":      string(marker),
		}, nil)
	}

	//Parse MetaBlocks
	for {
		blockMeta := &MetaBlock{}
		if isLast, err := blockMeta.Parse(br); err != nil {
			return types.NewException(TMFlac_CanNotParse_MetaDataBlock, map[string]string{
				"n": strconv.Itoa(len(fObj.blocks)),
			}, err)
		} else {
			fObj.blocks = append(fObj.blocks, blockMeta)
			if isLast {
				break
			}
		}
	}

	fObj.ResetMetaBlockProperties()

	//Read Frames
	if FrameBytes, err := br.ReadAllFollowedBytes(); err != nil {
		types.NewException(TMFlac_CanNotRead_Frames, nil, nil)
	} else {
		fObj.frames = FrameBytes
		types.Throw(types.NewException(TMFlac_Read_Frames, map[string]string{
			"length": strconv.Itoa(len(FrameBytes)),
		}, nil), types.RsInfo)
	}

	return nil
}

func (fObj *Flac) ParseFromFile(path string) *types.Exception {
	file, err := os.Open(path)
	if err != nil {
		return types.NewException(TMFlac_CanNotOpen_File, nil, err)
	}

	defer func() {
		_ = file.Close()
	}()

	return fObj.Parse(types.NewBinaryReader(file))
}

func (fObj *Flac) Encode() (*types.Buffer, *types.Exception) {
	buffer := types.NewBuffer()

	if _, err := buffer.Write(flacSignature); err != nil {
		return nil, types.NewException(TMFlac_CanNotWrite_FileSignature, nil, err)
	}

	for bIndex, block := range fObj.blocks {
		if dataBuffer, err := block.Encode(bIndex == len(fObj.blocks)-1); err != nil {
			//Encode block
			return nil, types.NewException(TMFlac_CanNotEncode_MetaDataBlock, map[string]string{
				"n": strconv.Itoa(bIndex),
			}, err)
		} else if dataDumped, err := dataBuffer.Dump(); err != nil {
			//Dump bytes
			return nil, types.NewException(TMFlac_CanNotDump_MetaDataBlock, map[string]string{
				"n": strconv.Itoa(bIndex),
			}, err)
		} else if wroteLen, err := buffer.Write(dataDumped); err != nil || wroteLen != len(dataDumped) {
			//Write to buffer
			return nil, types.NewException(TMFlac_CanNotWrite_MetaDataBlock, map[string]string{
				"n": strconv.Itoa(bIndex),
			}, err)
		}
	}

	if _, err := buffer.Write(fObj.frames); err != nil {
		return nil, types.NewException(TMFlac_CanNotWrite_Frames, nil, err)
	}

	return buffer, nil
}

func (fObj *Flac) WriteToFile(path string) *types.Exception {
	var fileBytes []byte

	if dataBuffer, err := fObj.Encode(); err != nil {
		return types.NewException(TMFlac_CanNotSaveTo_File, nil, err)
	} else if dataDumped, err := dataBuffer.Dump(); err != nil {
		return types.NewException(TMFlac_CanNotSaveTo_File, nil, err)
	} else {
		fileBytes = dataDumped
	}

	if err := base.FileWriteBytes(path, fileBytes); err != nil {
		return types.NewException(TMFlac_CanNotSaveTo_File, nil, err)
	}

	return nil
}

func (fObj *Flac) ResetMetaBlockProperties() {
	for blockIndex, block := range fObj.blocks {
		block.GetTags().Set("index", strconv.Itoa(blockIndex), nil)
		if blockIndex == len(fObj.blocks)-1 {
			block.GetTags().Set("isLast", "true", TagMatcher_Bool)
		} else {
			block.GetTags().Set("isLast", "false", TagMatcher_Bool)
		}
	}
}

func (fObj *Flac) Initialized() bool {
	return fObj.blocks != nil && fObj.frames != nil
}

func (fObj *Flac) GetBlocks() []*MetaBlock {
	return fObj.blocks
}

func (fObj *Flac) FindBlocks(pattern string) ([]int, *types.Exception) {
	patternTrimmed := strings.TrimSpace(pattern)
	if patternTrimmed == "" {
		return nil, types.NewException(
			TMFlac_FailedTo_Parse_FilterPattern,
			map[string]string{"pattern": pattern},
			types.Exception_Mismatched_Format(
				"BlockType or BlockType:TagName1=Filter1[&TagName2=Filter2...]",
				"(Empty)",
			),
		)
	}

	var blockType string
	var blockFilters url.Values = nil

	patternSplit := strings.SplitN(patternTrimmed, ":", 2)
	blockType = strings.TrimSpace(patternSplit[0])
	if len(patternSplit) == 2 {
		if blockFiltersTemp, err := url.ParseQuery(strings.TrimSpace(patternSplit[1])); err != nil {
			return nil, types.NewException(
				TMFlac_FailedTo_Parse_FilterPattern,
				map[string]string{"pattern": pattern},
				err,
			)
		} else if len(blockFiltersTemp) == 0 {
			return nil, types.NewException(
				TMFlac_FailedTo_Parse_FilterPattern,
				map[string]string{"pattern": pattern},
				types.Exception_Mismatched_Format(
					"BlockType or BlockType:TagName1=Filter1[&TagName2=Filter2...]",
					"BlockType:(Empty)",
				),
			)
		} else {
			blockFilters = blockFiltersTemp
		}
	}

	result := make([]int, 0)

	for blockIndex, block := range fObj.blocks {
		if block.Matches(blockType, blockFilters) {
			result = append(result, blockIndex)
		}
	}

	return result, nil
}

func (fObj *Flac) GetBlockByIndex(index int) *MetaBlock {
	if index >= 0 && index < len(fObj.blocks) {
		return fObj.blocks[index]
	}
	return nil
}

func (fObj *Flac) GetBlocksByIndexSlice(indexes []int) []*MetaBlock {
	blocks := make([]*MetaBlock, 0)
	if len(indexes) == 0 {
		return blocks
	}
	for _, index := range indexes {
		if index >= 0 && index < len(fObj.blocks) {
			blocks = append(blocks, fObj.blocks[index])
		}
	}
	return blocks
}

func (fObj *Flac) SetBlocks(blocks []*MetaBlock) {
	fObj.blocks = blocks
	fObj.ResetMetaBlockProperties()
}

func (fObj *Flac) AppendBlock(block *MetaBlock) {
	fObj.blocks = append(fObj.blocks, block)
	fObj.ResetMetaBlockProperties()
}
