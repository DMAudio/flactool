package flac

import (
	"p20190417/types"
	"p20190417/util"
)

type MetaBlockT1PADD struct {
	data []byte
}

func (mb *MetaBlockT1PADD) Parse(r *util.BinaryReader) *types.Exception {
	if data, err := r.ReadAllFollowedBytes(); err != nil {
		return types.NewException(TMFlac_CanNotREAD_MetaT1Data, nil, err)
	} else {
		mb.data = data
	}
	return nil
}

func (mb *MetaBlockT1PADD) Encode() (*types.Buffer, *types.Exception) {
	buffer := types.NewBuffer()
	if _, err := buffer.Write(mb.data); err != nil {
		return nil, types.NewException(TMFlac_CanNotParseMetaT1Data, nil, err)
	}

	return buffer, nil
}
