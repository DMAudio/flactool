package flac

import (
	"bytes"
	"fmt"
	"p20190417/types"
)

type MetaBlockType uint8

const (
	MetaBlockType_STREAMINFO     = 0
	MetaBlockType_PADDING        = 1
	MetaBlockType_APPLICATION    = 2
	MetaBlockType_SEEKTABLE      = 3
	MetaBlockType_VORBIS_COMMENT = 4
	MetaBlockType_CUESHEET       = 5
	MetaBlockType_PICTURE        = 6
)

const (
	MetaBlockTypeStr_STREAMINFO     = "STREAMINFO"
	MetaBlockTypeStr_PADDING        = "PADDING"
	MetaBlockTypeStr_APPLICATION    = "APPLICATION"
	MetaBlockTypeStr_SEEKTABLE      = "SEEKTABLE"
	MetaBlockTypeStr_VORBIS_COMMENT = "VORBIS_COMMENT"
	MetaBlockTypeStr_CUESHEET       = "CUESHEET"
	MetaBlockTypeStr_PICTURE        = "PICTURE"
)

type MetaBlock struct {
	Type MetaBlockType
	Body MetaBlockBody
}

type MetaBlockBody interface {
	Parse(r *types.BinaryReader) *types.Exception
	Encode() (*types.Buffer, *types.Exception)

	GetTags() *MetaBlockTags
}

func (m *MetaBlock) Parse(br *types.BinaryReader) (bool, *types.Exception) {
	var blockHead uint64
	var blockSize uint64

	//头部
	if blockHeadData, err := br.ReadBytes(1); err != nil {
		return true, types.NewException(TMFlac_CanNotParse_MetaBlockHead, nil, err)
	} else {
		blockHead = types.BytesToUInt64(blockHeadData)
	}

	//1bit 是否为最后一个数据块
	isLastFlag := blockHead>>7 == 1
	//7bit 数据块类型
	m.Type = MetaBlockType(uint8(blockHead & 0x7F))

	//数据长度 3byte
	if blockSizeData, err := br.ReadBytes(3); err != nil {
		return isLastFlag, types.NewException(TMFlac_CanNotParse_MetaBlockSIZE, nil, err)
	} else {
		blockSize = types.BytesToUInt64(blockSizeData)
	}

	//原始数据
	if blockData, err := br.ReadBytes(blockSize); err != nil {
		return isLastFlag, types.NewException(TMFlac_CanNotRead_MetaBlockData, nil, err)
	} else {
		if blockBody, err := m.ParseBody(blockData); err != nil {
			return isLastFlag, types.NewException(TMFlac_CanNotParse_MetaBlockData, nil, err)
		} else {
			m.Body = blockBody
		}
	}

	types.Throw(types.NewException(TMFlac_Parsed_MetaBlock, map[string]string{
		"type":   m.GetTypeStr(true),
		"length": fmt.Sprintf("%08d", blockSize),
	}, nil), types.RsInfo)

	return isLastFlag, nil
}

func (m *MetaBlock) ParseBody(data []byte) (MetaBlockBody, *types.Exception) {
	var body MetaBlockBody
	r := bytes.NewReader(data)

	switch m.Type {
	//case MetaBlockType_STREAMINFO:
	//	body = &MetaBlockT0STRE{}
	case MetaBlockType_PADDING:
		body = &MetaBlockT1PADD{}
	case MetaBlockType_VORBIS_COMMENT:
		body = &MetaBlockT4VORB{}
	case MetaBlockType_PICTURE:
		body = &MetaBlockT6PICT{}
	default:
		body = &MetaBlockT1PADD{}

	}

	err := body.Parse(types.NewBinaryReader(r))

	return body, err
}

func (m *MetaBlock) Encode(isLast bool) (*types.Buffer, *types.Exception) {
	buffer := types.NewBuffer()

	var blockHead uint8
	if isLast {
		blockHead = uint8(1<<7 | m.Type)
	} else {
		blockHead = uint8(m.Type)
	}

	if blockHeadBytes, err := types.UIntToBytes(uint64(blockHead), 1); err != nil {
		return nil, types.NewException(TMFlac_CanNotEncode_MetaBlockHead, nil, err)
	} else if _, err := buffer.Write(blockHeadBytes); err != nil {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaBlockHead, nil, err)
	}

	var bodyData []byte
	if bodyDataBuffer, err := m.Body.Encode(); err != nil {
		return nil, types.NewException(TMFlac_CanNotEncode_MetaBlockData, nil, err)
	} else if bodyDataDumped, err := bodyDataBuffer.Dump(); err != nil {
		return nil, types.NewException(TMFlac_CanNotDump_MetaBlockData, nil, err)
	} else {
		bodyData = bodyDataDumped
	}

	if blockBodySizeBytes, err := types.UIntToBytes(uint64(len(bodyData)), 3); err != nil {
		return nil, types.NewException(TMFlac_CanNotEncode_MetaBlockBodySize, nil, err)
	} else if _, err := buffer.Write(blockBodySizeBytes); err != nil {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaBlockBodySize, nil, err)
	}

	//dump, err := buffer.DumpList()
	//fmt.Printf("%08b\t%v", dump, err)

	if _, err := buffer.Write(bodyData); err != nil {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaBlockBodySize, nil, err)
	}

	return buffer, nil
}

func (m *MetaBlock) GetTypeStr(extend bool) string {
	typeStr := ""
	switch m.Type {
	case MetaBlockType_STREAMINFO:
		typeStr = MetaBlockTypeStr_STREAMINFO
	case MetaBlockType_PADDING:
		typeStr = MetaBlockTypeStr_PADDING
	case MetaBlockType_APPLICATION:
		typeStr = MetaBlockTypeStr_APPLICATION
	case MetaBlockType_SEEKTABLE:
		typeStr = MetaBlockTypeStr_SEEKTABLE
	case MetaBlockType_VORBIS_COMMENT:
		typeStr = MetaBlockTypeStr_VORBIS_COMMENT
	case MetaBlockType_CUESHEET:
		typeStr = MetaBlockTypeStr_CUESHEET
	case MetaBlockType_PICTURE:
		typeStr = MetaBlockTypeStr_PICTURE
	}

	return typeStr
}
