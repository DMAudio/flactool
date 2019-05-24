package flac

import (
	"dadp.flactool/types"
)

type PrependID3v2 struct {
	data []byte
}

func (pID3 *PrependID3v2) Parse(br *types.BinaryReader) *types.Exception {
	if _, err := br.ReadBytes(2); err != nil {
		return types.NewException(TMFlac_CanNotParse_ID3V2BlockSIZE, nil, err)
	}

	var blockSize uint64
	if blockSizeData, err := br.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotParse_ID3V2BlockSIZE, nil, err)
	} else {
		blockSize = types.BytesToInt64(blockSizeData)
	}

	if blockData, err := br.ReadBytes(blockSize); err != nil {
		return types.NewException(TMFlac_CanNotParse_ID3V2BlockData, nil, err)
	} else {
		pID3.data = blockData
	}

	return nil
}
