package flac

import (
	"bytes"
	"dadp.flactool/types"
	"dadp.flactool/utils"
	"encoding/binary"
	_ "golang.org/x/image/webp"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime"
	"strconv"
	"strings"
)

type MetaBlockT6PICT struct {
	picType        uint32
	picMime        string
	picDesc        string
	picWidth       uint32
	picHeight      uint32
	picColorDepth  uint32
	picColorAmount uint32
	picRawData     []byte
	bodyTag        *MetaBlockTags
}

func (mb *MetaBlockT6PICT) Parse(r *types.BinaryReader) *types.Exception {
	//Picture type according to the ID3v2 APIC frame
	if TypeData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6Type, nil, err)
	} else {
		mb.picType = binary.BigEndian.Uint32(TypeData)
	}

	//MIME type string
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

	//Description (in UTF-8)
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

	//Width
	if WidthData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6Width, nil, err)
	} else {
		mb.picWidth = binary.BigEndian.Uint32(WidthData)
	}

	//Height
	if HeightData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6Height, nil, err)
	} else {
		mb.picHeight = binary.BigEndian.Uint32(HeightData)
	}

	//Color depth
	if ColorDepthData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6ColorDepth, nil, err)
	} else {
		mb.picColorDepth = binary.BigEndian.Uint32(ColorDepthData)
	}

	//Number of colors used (for indexed-color pictures (e.g. GIF))
	if ColorsData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6Colors, nil, err)
	} else {
		mb.picColorAmount = binary.BigEndian.Uint32(ColorsData)
	}

	//Binary picture data
	if PicLengthData, err := r.ReadBytes(4); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6DataLength, nil, err)
	} else {
		PicLength := uint64(binary.BigEndian.Uint32(PicLengthData))
		if PicData, err := r.ReadBytes(PicLength); err != nil {
			return types.NewException(TMFlac_CanNotRead_MetaT6Data, nil, err)
		} else {
			mb.picRawData = PicData
		}
	}

	return nil
}

func (mb *MetaBlockT6PICT) ParsePictureFile(path string) *types.Exception {
	var err error
	var exception *types.Exception

	var picBytes []byte
	if picBytes, exception = utils.FileReadBytes(path); exception != nil {
		return exception
	}
	var imgObj image.Image
	var imgFormat string
	if imgObj, imgFormat, err = image.Decode(bytes.NewReader(picBytes)); err != nil {
		return types.NewException(TMFlac_CanNotRead_MetaT6Data, nil, exception)
	}

	mb.picMime = mime.TypeByExtension("." + imgFormat)
	mb.picWidth = uint32(imgObj.Bounds().Size().X)
	mb.picHeight = uint32(imgObj.Bounds().Size().Y)
	mb.picRawData = picBytes

	var imgColorModel color.Model
	imgColorModel = imgObj.ColorModel()

	switch imgColorModel.(type) {
	case color.Palette:
		mb.picColorAmount = uint32(len(imgColorModel.(color.Palette)))
	default:
		mb.picColorAmount = 0
	}

	switch imgObj.At(0, 0).(type) {
	case color.Alpha:
		mb.picColorDepth = 8
	case color.Alpha16:
		mb.picColorDepth = 16
	case color.CMYK:
		mb.picColorDepth = 4 * 8
	case color.Gray:
		mb.picColorDepth = 8
	case color.Gray16:
		mb.picColorDepth = 16
	case color.NRGBA:
		mb.picColorDepth = 4 * 8
	case color.NRGBA64:
		mb.picColorDepth = 4 * 16
	case color.NYCbCrA:
		mb.picColorDepth = 4 * 8
	case color.RGBA:
		mb.picColorDepth = 4 * 8
	case color.RGBA64:
		mb.picColorDepth = 4 * 16
	case color.YCbCr:
		mb.picColorDepth = 3 * 8
	default:
		mb.picColorDepth = 0

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
	binary.BigEndian.PutUint32(DataLength, uint32(len(mb.picRawData)))
	if length, err := buffer.Write(DataLength); err != nil || length != len(DataLength) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT6DataLength, nil, err)
	}

	if length, err := buffer.Write(mb.picRawData); err != nil || length != len(mb.picRawData) {
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

func (mb *MetaBlockT6PICT) SetPicDesc(picDesc string) {
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
	return mb.picRawData
}

func (mb *MetaBlockT6PICT) GetType() MetaBlockType {
	return MetaBlockType_PICTURE
}

func (mb *MetaBlockT6PICT) GetTags() *MetaBlockTags {
	if mb.bodyTag == nil {
		mb.bodyTag = NewMetaBlockTags()
	}

	mb.bodyTag.Set("type", strconv.FormatUint(uint64(mb.picType), 10), nil)
	mb.bodyTag.Set("mime", mb.picMime, nil)
	mb.bodyTag.Set("desc", mb.picDesc, nil)
	mb.bodyTag.Set("width", strconv.FormatUint(uint64(mb.picWidth), 10), nil)
	mb.bodyTag.Set("height", strconv.FormatUint(uint64(mb.picHeight), 10), nil)
	mb.bodyTag.Set("colorDepth", strconv.FormatUint(uint64(mb.picColorDepth), 10), nil)
	mb.bodyTag.Set("colorAmount", strconv.FormatUint(uint64(mb.picColorAmount), 10), nil)
	mb.bodyTag.Set("fileSize", strconv.Itoa(len(mb.picRawData)), nil)

	if ext, err := mime.ExtensionsByType(mb.picMime); err == nil && len(ext) > 0 {
		mb.bodyTag.Set("fileExt", strings.TrimPrefix(ext[0], "."), nil)
	}

	return mb.bodyTag
}
