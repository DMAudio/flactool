package flac

import (
	"encoding/binary"
	"p20190417/types"
)

type MetaBlockT6PICT struct {
	picType   uint32
	picMime   string
	picDesc   string
	picWidth  uint32
	picHeight uint32
	picCDepth uint32
	picColors uint32
	picData   []byte
}

func (mb *MetaBlockT6PICT) Parse(r *types.BinaryReader) *types.Exception {
	//图片类型
	if TypeData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6Type, nil, err)
	} else {
		mb.picType = binary.BigEndian.Uint32(TypeData)
	}

	//MIME信息
	if MIMETypeLengthData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6MIMELength, nil, err)
	} else {
		MimeTypeLength := uint64(binary.BigEndian.Uint32(MIMETypeLengthData))
		if MIMETypeData, err := r.ReadBytes(MimeTypeLength); err != nil {
			return types.NewException(TMFlac_CanNotRead_MetaT6MIME, nil, err)
		} else {
			mb.picMime = string(MIMETypeData)
		}
	}

	//图片介绍
	if DescLengthData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6DescriptionLength, nil, err)
	} else {
		DescLength := uint64(binary.BigEndian.Uint32(DescLengthData))
		if descData, err := r.ReadBytes(DescLength); err != nil {
			return types.NewException(TMFlac_CanNotRead_MetaT6Description, nil, err)
		} else {
			mb.picDesc = string(descData)
		}
	}

	//宽度
	if WidthData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6Width, nil, err)
	} else {
		mb.picWidth = binary.BigEndian.Uint32(WidthData)
	}

	//长度
	if HeightData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6Height, nil, err)
	} else {
		mb.picHeight = binary.BigEndian.Uint32(HeightData)
	}

	//颜色深度
	if ColorDepthData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6ColorDepth, nil, err)
	} else {
		mb.picCDepth = binary.BigEndian.Uint32(ColorDepthData)
	}

	//色数(用于 GIF 等图片格式)
	if ColorsData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6Colors, nil, err)
	} else {
		mb.picColors = binary.BigEndian.Uint32(ColorsData)
	}

	//图片原始数据
	if PicLengthData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6DataLength, nil, err)
	} else {
		PicLength := uint64(binary.BigEndian.Uint32(PicLengthData))
		if PicData, err := r.ReadBytes(PicLength); err != nil {
			return types.NewException(TMFlac_CanNotRead_MetaT6Data, nil, err)
		} else {
			mb.picData = PicData
		}
	}

	return nil
}

func (mb *MetaBlockT6PICT) Encode() (*types.Buffer, *types.Exception) {
	buffer := types.NewBuffer()

	TypeData := make([]byte, 4)
	binary.BigEndian.PutUint32(TypeData, uint32(mb.picType))
	if length, err := buffer.Write(TypeData); err != nil || length != len(TypeData) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6Type, nil, err)
	}

	MIMESize := make([]byte, 4)
	binary.BigEndian.PutUint32(MIMESize, uint32(len(mb.picMime)))
	if length, err := buffer.Write(MIMESize); err != nil || length != len(MIMESize) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6MIMELength, nil, err)
	}

	if length, err := buffer.WriteString(mb.picMime); err != nil || length != len(mb.picMime) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6MIME, nil, err)
	}

	DescSize := make([]byte, 4)
	binary.BigEndian.PutUint32(DescSize, uint32(len(mb.picDesc)))
	if length, err := buffer.Write(DescSize); err != nil || length != len(DescSize) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6DescriptionLength, nil, err)
	}

	if length, err := buffer.WriteString(mb.picDesc); err != nil || length != len(mb.picDesc) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6Description, nil, err)
	}

	PicWidth := make([]byte, 4)
	binary.BigEndian.PutUint32(PicWidth, uint32(mb.picWidth))
	if length, err := buffer.Write(PicWidth); err != nil || length != len(PicWidth) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6Width, nil, err)
	}

	PicHeight := make([]byte, 4)
	binary.BigEndian.PutUint32(PicHeight, uint32(mb.picHeight))
	if length, err := buffer.Write(PicHeight); err != nil || length != len(PicHeight) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6Height, nil, err)
	}

	ColorDepth := make([]byte, 4)
	binary.BigEndian.PutUint32(ColorDepth, uint32(mb.picCDepth))
	if length, err := buffer.Write(ColorDepth); err != nil || length != len(ColorDepth) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6ColorDepth, nil, err)
	}

	Colors := make([]byte, 4)
	binary.BigEndian.PutUint32(Colors, uint32(mb.picColors))
	if length, err := buffer.Write(Colors); err != nil || length != len(Colors) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6Colors, nil, err)
	}

	DataLength := make([]byte, 4)
	binary.BigEndian.PutUint32(DataLength, uint32(len(mb.picData)))
	if length, err := buffer.Write(DataLength); err != nil || length != len(DataLength) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6DataLength, nil, err)
	}

	if length, err := buffer.Write(mb.picData); err != nil || length != len(mb.picData) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6Data, nil, err)
	}

	return buffer, nil
}

func (mb *MetaBlockT6PICT) GetPicType() int {
	return int(mb.picType)
}

func (mb *MetaBlockT6PICT) GetTags() *MetaBlockTags {
	m := NewMetaBlockTags()

	return m
}
