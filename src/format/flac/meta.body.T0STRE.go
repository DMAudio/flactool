package flac

import (
	"p20190417/types"
	"p20190417/util"
)

//不可用

type MetaBlockT0STRE struct {
	data []byte
}

func (mb *MetaBlockT0STRE) Parse(r *util.BinaryReader) *types.Exception {
	//TODO
	return nil
}

func (mb *MetaBlockT0STRE) Encode() (*types.Buffer, *types.Exception) {
	buffer := types.NewBuffer()
	//TODO
	return buffer, nil
}
