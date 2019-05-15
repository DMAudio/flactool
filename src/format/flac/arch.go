package flac

import (
	"bytes"
	"os"
	"p20190417/types"
	"strconv"
	"sync"
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

	//解析前置ID3V2
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

	//检测fLaC标头
	if !bytes.Equal(marker, flacSignature) {
		return types.NewException(TMFlac_Incorrect_FileSignature, map[string]string{
			"expected": string(flacSignature),
			"got":      string(marker),
		}, nil)
	}

	//解析Meta数据块
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

	//读取Frames
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
			return nil, types.NewException(TMFlac_CanNotEncode_MetaDataBlock, map[string]string{
				"n": strconv.Itoa(bIndex),
			}, err)
		} else if dataDumped, err := dataBuffer.Dump(); err != nil {
			return nil, types.NewException(TMFlac_CanNotDump_MetaDataBlock, map[string]string{
				"n": strconv.Itoa(bIndex),
			}, err)
		} else if wroteLen, err := buffer.Write(dataDumped); err != nil || wroteLen != len(dataDumped) {
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
		return types.NewException(TMFlac_CanNotWrite_File, nil, err)
	} else if dataDumped, err := dataBuffer.Dump(); err != nil {
		return types.NewException(TMFlac_CanNotWrite_File, nil, err)
	} else {
		fileBytes = dataDumped
	}

	if file, err := os.Create(path); err != nil {
		panic(err)
	} else {
		defer func() { _ = file.Close() }()
		if _, err := file.Write(fileBytes); err != nil {
			return types.NewException(TMFlac_CanNotWrite_File, nil, err)
		}
	}
	return nil
}

func (fObj *Flac) Initialized() bool {
	return fObj.blocks != nil && fObj.frames != nil
}

func (fObj *Flac) GetBlocks() []*MetaBlock {
	return fObj.blocks
}

func (fObj *Flac) SetBlocks(blocks []*MetaBlock) {
	fObj.blocks = blocks
}

func (fObj *Flac) FindBlock(pattern string) []int {
	result := make([]int, 0)

	for blockIndex, block := range fObj.blocks {
		blockExtraTags := NewMetaBlockTags()
		blockExtraTags.Set("index", strconv.Itoa(blockIndex), nil)
		if block.MatchesPattern(pattern, map[string]*MetaBlockTags{
			"block": blockExtraTags,
		}) {
			result = append(result, blockIndex)
		}
	}

	return result
}

var globalFlac *Flac
var globalFlacLock sync.Mutex

func GlobalFlac() *Flac {
	if globalFlac == nil {
		globalFlacLock.Lock()
		defer globalFlacLock.Unlock()
		if globalFlac == nil {
			globalFlac = &Flac{}
		}
	}
	return globalFlac
}
