package flac

import (
	"gitlab.com/MGEs/Com.Base/types"
)

var TMFlac_CanNotOpen_File = types.NewMask(
	"CANNOT_OPEN_FILE",
	"无法打开文件",
)

var TMFlac_CanNotSaveTo_File = types.NewMask(
	"CANNOT_SaveTo_FILE",
	"无法保存到文件",
)

var TMFlac_UndefinedObject = types.NewMask(
	"UNDEFINED_OBJECT",
	"对象未定义",
)

var TMFlac_UninitializedObject = types.NewMask(
	"UNINITIALIZED_OBJECT",
	"对象未初始化",
)

var TMFlac_CanNotAssert_FlacObject = types.NewMask(
	"CAN_NOT_ASSERT_FlacObject",
	"断言对象为*Flac失败",
)

var TMFlac_CanNotRead_FileSignature = types.NewMask(
	"CAN_NOT_READ_FILE_SIGNATURE",
	"无法读取文件标头",
)

var TMFlac_CanNotWrite_FileSignature = types.NewMask(
	"CAN_NOT_WRITE_FILE_SIGNATURE",
	"无法写入文件标头",
)

var TMFlac_Incorrect_FileSignature = types.NewMask(
	"CAN_INCORRECT_FILE_SIGNATURE",
	"文件标头错误，预期值：{{expected}}，实际值：{{got}}",
)

var TMFlac_CanNotParse_PrependID3V2Block = types.NewMask(
	"CAN_NOT_PARSE_PREPEND_ID3V2_BLOCK",
	"无法解析预置的ID3V2数据块",
)

var TMFlac_CanNotParse_ID3V2BlockSIZE = types.NewMask(
	"CAN_NOT_PARSE_ID3V2BLOCK_SIZE",
	"无法解析ID3V2数据块的大小",
)

var TMFlac_CanNotParse_ID3V2BlockData = types.NewMask(
	"CAN_NOT_PARSE_ID3V2BLOCK_SIZE",
	"无法解析ID3V2数据块的实际数据",
)

var TMFlac_CanNotParse_MetaDataBlock = types.NewMask(
	"CAN_NOT_PARSE_METADATA_BLOCK",
	"无法解析第{{n}}个元数据块",
)

var TMFlac_CanNotEncode_MetaDataBlock = types.NewMask(
	"CAN_NOT_PARSE_METADATA_BLOCK",
	"无法编码第{{n}}个元数据块",
)

var TMFlac_CanNotDump_MetaDataBlock = types.NewMask(
	"CAN_NOT_DUMP_METADATA_BLOCK",
	"无法导出第{{n}}个元数据块的数据",
)

var TMFlac_CanNotWrite_MetaDataBlock = types.NewMask(
	"CAN_NOT_WRITE_METADATA_BLOCK",
	"无法导出将{{n}}个元数据块的数据写入文件",
)

var TMFlac_Parsed_MetaBlock = types.NewMask(
	"PARSED_METABLOCK",
	"成功解析了一个元数据块（类型：{{type}}，长度：4+{{length}} Bytes）",
)

var TMFlac_CanNotParse_MetaBlockHead = types.NewMask(
	"CAN_NOT_PARSE_METABLOCK_HEAD",
	"无法解析Meta数据块的头部",
)

var TMFlac_CanNotParse_MetaBlockSIZE = types.NewMask(
	"CAN_NOT_PARSE_METABLOCK_SIZE",
	"无法解析Meta数据块的大小",
)

var TMFlac_CanNotRead_MetaBlockData = types.NewMask(
	"CAN_NOT_PARSE_METABLOCK_DATA",
	"无法读取Meta数据块的实际数据",
)

var TMFlac_CanNotParse_MetaBlockData = types.NewMask(
	"CAN_NOT_PARSE_METABLOCK_DATA",
	"无法解析Meta数据块的实际数据",
)

var TMFlac_CanNotEncode_MetaBlockData = types.NewMask(
	"CAN_NOT_ENCODE_METABLOCK_DATA",
	"无法导出Meta数据块",
)

var TMFlac_CanNotDump_MetaBlockData = types.NewMask(
	"CAN_NOT_DUMP_METABLOCK_DATA",
	"无法导出Meta数据块的实际数据",
)

var TMFlac_CanNotEncode_MetaBlockHead = types.NewMask(
	"CAN_NOT_ENCODE_METABLOCK_Head",
	"无法编码Meta数据块的头部",
)

var TMFlac_CanNotWrite_MetaBlockHead = types.NewMask(
	"CAN_NOT_WRITE_METABLOCK_Head",
	"无法写入Meta数据块的头部",
)

var TMFlac_CanNotEncode_MetaBlockBodySize = types.NewMask(
	"CAN_NOT_ENCODE_METABLOCK_BODYSIZE",
	"无法编码Meta数据块的实际数据大小",
)

var TMFlac_CanNotWrite_MetaBlockBodySize = types.NewMask(
	"CAN_NOT_WRITE_METABLOCK_BODYSIZE",
	"无法写入Meta数据块的实际数据大小",
)

var TMFlac_CanNotREAD_MetaT1Data = types.NewMask(
	"CAN_NOT_READ_META_T1_Refer_Data",
	"无法读取PADDING数据块的内容",
)

var TMFlac_CanNotParseMetaT1Data = types.NewMask(
	"CAN_NOT_Parse_META_T1_Refer_Data",
	"无法编码PADDING数据块的内容",
)

var TMFlac_MetaT4_NotFound = types.NewMask(
	"CAN_NOT_FIND_META_T4",
	"找不到VORBIS_COMMENT数据块",
)

var TMFlac_CanNotParse_MetaT4ReferSize = types.NewMask(
	"CAN_NOT_PARSE_META_T4_Refer_Size",
	"无法解析VORBIS_COMMENT的refer信息长度",
)

var TMFlac_CanNotWrite_MetaT4ReferSize = types.NewMask(
	"CAN_NOT_Write_META_T4_Refer_Size",
	"无法写入VORBIS_COMMENT的refer信息长度",
)

var TMFlac_CanNotParse_MetaT4ReferData = types.NewMask(
	"CAN_NOT_PARSE_META_T4_Refer_Data",
	"无法解析VORBIS_COMMENT的refer信息内容",
)

var TMFlac_CanNotWrite_MetaT4ReferData = types.NewMask(
	"CAN_NOT_WRITE_META_T4_Refer_Data",
	"无法写入VORBIS_COMMENT的refer信息内容",
)

var TMFlac_CanNotParse_MetaT4CommentAmount = types.NewMask(
	"CAN_NOT_PARSE_META_T4_Comment_Amount",
	"无法解析VORBIS_COMMENT的Comment列表长度",
)

var TMFlac_CanNotWrite_MetaT4CommentAmount = types.NewMask(
	"CAN_NOT_WRITE_META_T4_Comment_Amount",
	"无法写入VORBIS_COMMENT的Comment列表长度",
)

var TMFlac_CanNotParse_MetaT4CommentItemLength = types.NewMask(
	"CAN_NOT_PARSE_META_T4_Comment_Length",
	"无法解析VORBIS_COMMENT的Comment条目长度",
)

var TMFlac_CanNotWrite_MetaT4CommentItemLength = types.NewMask(
	"CAN_NOT_WRITE_META_T4_Comment_Length",
	"无法写入VORBIS_COMMENT的Comment条目长度",
)

var TMFlac_CanNotParse_MetaT4CommentData = types.NewMask(
	"CAN_NOT_PARSE_META_T4_Comment_Data",
	"无法解析VORBIS_COMMENT的Comment条目内容",
)

var TMFlac_CanNotDump_MetaT4CommentList = types.NewMask(
	"CAN_NOT_DUMP_META_T4_Comment_List",
	"无法导出VORBIS_COMMENT的Comment列表",
)

var TMFlac_CanNotAssert_METABLOCKAsSpecificType = types.NewMask(
	"CAN_NOT_ASSERT_METABLOCK_AS_SpecificType",
	"断言第 {{index}} 个元数据块为 {{type}} 类型失败",
)

var TMFlac_CanNotRead_MetaT6Type = types.NewMask(
	"CAN_NOT_READ_META_T6_Type",
	"无法读取PICTURE的类型",
)

var TMFlac_CanNotWrite_MetaT6Type = types.NewMask(
	"CAN_NOT_WRITE_META_T6_Type",
	"无法写入PICTURE的类型",
)

var TMFlac_CanNotRead_MetaT6MIMELength = types.NewMask(
	"CAN_NOT_READ_META_T6_MIME_Length",
	"无法读取PICTURE的MIME类型长度",
)

var TMFlac_CanNotWrite_MetaT6MIMELength = types.NewMask(
	"CAN_NOT_WRITE_META_T6_MIME_Length",
	"无法写入PICTURE的MIME类型长度",
)

var TMFlac_CanNotRead_MetaT6MIME = types.NewMask(
	"CAN_NOT_READ_META_T6_MIME",
	"无法读取PICTURE的MIME类型",
)

var TMFlac_CanNotWrite_MetaT6MIME = types.NewMask(
	"CAN_NOT_WRITE_META_T6_MIME",
	"无法写入PICTURE的MIME类型",
)

var TMFlac_CanNotRead_MetaT6DescriptionLength = types.NewMask(
	"CAN_NOT_READ_META_T6_DESCRIPTION_Length",
	"无法读取PICTURE的介绍长度",
)

var TMFlac_CanNotWrite_MetaT6DescriptionLength = types.NewMask(
	"CAN_NOT_WRITE_META_T6_DESCRIPTION_Length",
	"无法写入PICTURE的介绍长度",
)

var TMFlac_CanNotRead_MetaT6Description = types.NewMask(
	"CAN_NOT_READ_META_T6_DESCRIPTION",
	"无法读取PICTURE的介绍",
)

var TMFlac_CanNotWrite_MetaT6Description = types.NewMask(
	"CAN_NOT_WRITE_META_T6_DESCRIPTION",
	"无法写入PICTURE的介绍",
)

var TMFlac_CanNotRead_MetaT6Width = types.NewMask(
	"CAN_NOT_READ_META_T6_WIDTH",
	"无法读取PICTURE的宽度",
)

var TMFlac_CanNotWrite_MetaT6Width = types.NewMask(
	"CAN_NOT_WRITE_META_T6_WIDTH",
	"无法写入PICTURE的宽度",
)

var TMFlac_CanNotRead_MetaT6Height = types.NewMask(
	"CAN_NOT_READ_META_T6_HEIGHT",
	"无法读取PICTURE的长度",
)

var TMFlac_CanNotWrite_MetaT6Height = types.NewMask(
	"CAN_NOT_WRITE_META_T6_HEIGHT",
	"无法写入PICTURE的长度",
)

var TMFlac_CanNotRead_MetaT6ColorDepth = types.NewMask(
	"CAN_NOT_READ_META_T6_COLOR_DEPTH",
	"无法读取PICTURE的色深",
)

var TMFlac_CanNotWrite_MetaT6ColorDepth = types.NewMask(
	"CAN_NOT_WRITE_META_T6_COLOR_DEPTH",
	"无法写入PICTURE的色深",
)

var TMFlac_CanNotRead_MetaT6Colors = types.NewMask(
	"CAN_NOT_READ_META_T6_COLORS",
	"无法读取PICTURE的长度",
)

var TMFlac_CanNotWrite_MetaT6Colors = types.NewMask(
	"CAN_NOT_WRITE_META_T6_COLORS",
	"无法写入PICTURE的长度",
)

var TMFlac_CanNotRead_MetaT6DataLength = types.NewMask(
	"CAN_NOT_READ_META_T6_DATA_LENGTH",
	"无法读取PICTURE的图片原始数据长度",
)

var TMFlac_CanNotWrite_MetaT6DataLength = types.NewMask(
	"CAN_NOT_WRITE_META_T6_DATA_LENGTH",
	"无法写入PICTURE的图片原始数据长度",
)

var TMFlac_CanNotRead_MetaT6Data = types.NewMask(
	"CAN_NOT_READ_META_T6_DATA",
	"无法读取PICTURE的图片原始数据",
)

var TMFlac_CanNotWrite_MetaT6Data = types.NewMask(
	"CAN_NOT_WRITE_META_T6_DATA",
	"无法写入PICTURE的图片原始数据",
)

var TMFlac_Read_Frames = types.NewMask(
	"READ_FRAMES",
	"成功读取了音频源数据（长度：{{length}} Bytes）",
)

var TMFlac_CanNotRead_Frames = types.NewMask(
	"CAN_NOT_READ_FRAMES",
	"无法读取源数据",
)

var TMFlac_CanNotWrite_Frames = types.NewMask(
	"CAN_NOT_WRITE_FRAMES",
	"无法写入源数据",
)

var TMFlac_FailedTo_Parse_FilterPattern = types.NewMask(
	"FailedTo_Parse_FilterPattern",
	"无法解析过滤表达式:\n表达式：{{pattern}}",
)

var TMFlac_Arg_CanNotFind_Block = types.NewMask(
	"CAN_NOT_FIND_BLOCK",
	"找不到数据块: {{pattern}}",
)

var TMFlac_Arg_VaguePattern = types.NewMask(
	"Vague_Pattern",
	"通过表达式筛选的数据块不唯一: {{pattern}}",
)

var TMFlac_Arg_CanNotParseThisBlock = types.NewMask(
	"Can_Not_Parse_ThisBlock",
	"无法解析@this",
)
