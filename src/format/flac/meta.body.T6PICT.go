package flac

import (
	"encoding/binary"
	"p20190417/types"
	"strconv"
)

type MetaBlockT6PICT struct {
	picType        uint32
	picMime        string
	picDesc        string
	picWidth       uint32
	picHeight      uint32
	picColorDepth  uint32
	picColorAmount uint32
	picFileRawData []byte
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
		mb.picColorDepth = binary.BigEndian.Uint32(ColorDepthData)
	}

	//色数(用于 GIF 等图片格式)
	if ColorsData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6Colors, nil, err)
	} else {
		mb.picColorAmount = binary.BigEndian.Uint32(ColorsData)
	}

	//图片原始数据
	if PicLengthData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6DataLength, nil, err)
	} else {
		PicLength := uint64(binary.BigEndian.Uint32(PicLengthData))
		if PicData, err := r.ReadBytes(PicLength); err != nil {
			return types.NewException(TMFlac_CanNotRead_MetaT6Data, nil, err)
		} else {
			mb.picFileRawData = PicData
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
	binary.BigEndian.PutUint32(ColorDepth, uint32(mb.picColorDepth))
	if length, err := buffer.Write(ColorDepth); err != nil || length != len(ColorDepth) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6ColorDepth, nil, err)
	}

	Colors := make([]byte, 4)
	binary.BigEndian.PutUint32(Colors, uint32(mb.picColorAmount))
	if length, err := buffer.Write(Colors); err != nil || length != len(Colors) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6Colors, nil, err)
	}

	DataLength := make([]byte, 4)
	binary.BigEndian.PutUint32(DataLength, uint32(len(mb.picFileRawData)))
	if length, err := buffer.Write(DataLength); err != nil || length != len(DataLength) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6DataLength, nil, err)
	}

	if length, err := buffer.Write(mb.picFileRawData); err != nil || length != len(mb.picFileRawData) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6Data, nil, err)
	}

	return buffer, nil
}

func (mb *MetaBlockT6PICT) GetPicType() uint32 {
	return mb.picType
}

func (mb *MetaBlockT6PICT) SetPicType(picType uint32) {
	mb.picType = picType
}

func (mb *MetaBlockT6PICT) GetPicDesc() string {
	return mb.picDesc
}

func (mb *MetaBlockT6PICT) GSetPicDesc(picDesc string) {
	mb.picDesc = picDesc
}

func (mb *MetaBlockT6PICT) GetPicMime() string {
	return mb.picMime
}

func (mb *MetaBlockT6PICT) GetPicWidth() uint32 {
	return mb.picWidth
}

func (mb *MetaBlockT6PICT) GetPicHeight() uint32 {
	return mb.picHeight
}

func (mb *MetaBlockT6PICT) GetPicColorDepth() uint32 {
	return mb.picColorDepth
}

func (mb *MetaBlockT6PICT) GetPicColorAmount() uint32 {
	return mb.picColorAmount
}

func (mb *MetaBlockT6PICT) GetPicRawData() []byte {
	return mb.picFileRawData
}

func (mb *MetaBlockT6PICT) GetTags() *MetaBlockTags {
	m := NewMetaBlockTags()

	m.Set("type", strconv.FormatUint(uint64(mb.picType), 10), nil)
	m.Set("mime", mb.picMime, nil)
	m.Set("desc", mb.picDesc, nil)
	m.Set("width", strconv.FormatUint(uint64(mb.picWidth), 10), nil)
	m.Set("height", strconv.FormatUint(uint64(mb.picHeight), 10), nil)
	m.Set("colorDepth", strconv.FormatUint(uint64(mb.picColorDepth), 10), nil)
	m.Set("colorAmount", strconv.FormatUint(uint64(mb.picColorAmount), 10), nil)
	m.Set("fileSize", strconv.Itoa(len(mb.picFileRawData)), nil)

	return m
}
