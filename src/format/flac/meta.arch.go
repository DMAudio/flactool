package flac

import (
	"bytes"
	"fmt"
	"net/url"
	"p20190417/types"
	"strings"
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

	MetaBlockTypeStr_UNKNOWN = "UNKNOWN"
)

type MetaBlock struct {
	blockType MetaBlockType
	blockBody MetaBlockBody
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
	m.blockType = MetaBlockType(uint8(blockHead & 0x7F))

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
			m.blockBody = blockBody
		}
	}

	types.Throw(types.NewException(TMFlac_Parsed_MetaBlock, map[string]string{
		"type":   m.GetTypeStr(),
		"length": fmt.Sprintf("%08d", blockSize),
	}, nil), types.RsInfo)

	return isLastFlag, nil
}

func (m *MetaBlock) ParseBody(data []byte) (MetaBlockBody, *types.Exception) {
	var body MetaBlockBody
	r := bytes.NewReader(data)

	switch m.blockType {
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
		blockHead = uint8(1<<7 | m.blockType)
	} else {
		blockHead = uint8(m.blockType)
	}

	if blockHeadBytes, err := types.UIntToBytes(uint64(blockHead), 1); err != nil {
		return nil, types.NewException(TMFlac_CanNotEncode_MetaBlockHead, nil, err)
	} else if _, err := buffer.Write(blockHeadBytes); err != nil {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaBlockHead, nil, err)
	}

	var bodyData []byte
	if bodyDataBuffer, err := m.blockBody.Encode(); err != nil {
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

func (m *MetaBlock) GetType() MetaBlockType {
	return m.blockType
}

func (m *MetaBlock) GetTypeStr() string {
	switch m.blockType {
	case MetaBlockType_STREAMINFO:
		return MetaBlockTypeStr_STREAMINFO
	case MetaBlockType_PADDING:
		return MetaBlockTypeStr_PADDING
	case MetaBlockType_APPLICATION:
		return MetaBlockTypeStr_APPLICATION
	case MetaBlockType_SEEKTABLE:
		return MetaBlockTypeStr_SEEKTABLE
	case MetaBlockType_VORBIS_COMMENT:
		return MetaBlockTypeStr_VORBIS_COMMENT
	case MetaBlockType_CUESHEET:
		return MetaBlockTypeStr_CUESHEET
	case MetaBlockType_PICTURE:
		return MetaBlockTypeStr_PICTURE
	default:
		return MetaBlockTypeStr_UNKNOWN
	}
}

func (m *MetaBlock) GetBody() MetaBlockBody {
	return m.blockBody
}

func (m *MetaBlock) MatchesPattern(patternRaw string, extraTags map[string]*MetaBlockTags) bool {
	patternSplit := strings.SplitN(patternRaw, ":", 2)

	if blockType := patternSplit[0]; blockType != m.GetTypeStr() {
		return false
	}

	var filters map[string][]string
	if len(patternSplit) == 2 {
		if filtersTmp, err := url.ParseQuery(patternSplit[1]); err != nil {
			panic(err)
		} else {
			filters = filtersTmp
		}
	}

	if len(filters) == 0 {
		return true
	}

	var bodyTags *MetaBlockTags

	for key, values := range filters {
		if strings.HasPrefix(key, "body.") {
			if bodyTags == nil {
				bodyTags = m.GetBody().GetTags()
			}
			if !bodyTags.Match(strings.TrimPrefix(key, "body."), values) {
				return false
			}
			continue
		}

		for prefix, matcher := range extraTags {
			if !strings.HasPrefix(key, prefix+".") {
				continue
			}
			if !matcher.Match(strings.TrimPrefix(key, prefix+"."), values) {
				return false
			}
		}
	}

	return true
}
