package utils

import (
	"io/ioutil"
	"os"
	"p20190417/types"
	"path/filepath"
	"strings"
)

var TMFile_CanNotParse_FileAbsPath = types.NewMask(
	"CANNOT_Parse_File_AbsPath",
	"无法解析文件绝对路径：{{path}}",
)

var TMFile_CanNotCreate_File = types.NewMask(
	"CANNOT_Write_FILE",
	"无法创建文件：{{path}}",
)

var TMFile_CanNotWrite_File = types.NewMask(
	"CANNOT_Write_FILE",
	"无法写入文件：{{path}}",
)

var TMFile_CanNotOpen_File = types.NewMask(
	"CANNOT_Open_FILE",
	"无法打开文件：{{path}}",
)
var TMFile_CanNotRead_File = types.NewMask(
	"CANNOT_Read_FILE",
	"无法读取文件：{{path}}",
)

func FileWriteBytes(path string, bytes []byte) *types.Exception {
	if absPath, err := filepath.Abs(strings.TrimSpace(path)); err != nil {
		return types.NewException(TMFile_CanNotParse_FileAbsPath, nil, err)
	} else if file, err := os.Create(absPath); err != nil {
		return types.NewException(TMFile_CanNotCreate_File, map[string]string{
			"path": path,
		}, err)
	} else {
		defer func() { _ = file.Close() }()
		if _, err := file.Write(bytes); err != nil {
			return types.NewException(TMFile_CanNotWrite_File, map[string]string{
				"path": path,
			}, err)
		}
	}
	return nil
}

func FileWriteString(path string, str string) *types.Exception {
	if absPath, err := filepath.Abs(strings.TrimSpace(path)); err != nil {
		return types.NewException(TMFile_CanNotParse_FileAbsPath, nil, err)
	} else if file, err := os.Create(absPath); err != nil {
		return types.NewException(TMFile_CanNotCreate_File, map[string]string{
			"path": path,
		}, err)
	} else {
		defer func() { _ = file.Close() }()
		if _, err := file.WriteString(str); err != nil {
			return types.NewException(TMFile_CanNotWrite_File, map[string]string{
				"path": path,
			}, err)
		}
	}
	return nil
}

func FileGetReader(path string) (*os.File, *types.Exception) {
	if absPath, err := filepath.Abs(strings.TrimSpace(path)); err != nil {
		return nil, types.NewException(TMFile_CanNotParse_FileAbsPath, nil, err)
	} else if fileObj, err := os.OpenFile(absPath, os.O_RDONLY, 0644); err != nil {
		return nil, types.NewException(TMFile_CanNotOpen_File, map[string]string{
			"path": path,
		}, nil)
	} else {
		return fileObj, nil
	}
}

func FileReadBytes(path string) ([]byte, *types.Exception) {
	if fileObj, err := FileGetReader(path); err != nil {
		return nil, types.NewException(TMFile_CanNotRead_File, map[string]string{
			"path": path,
		}, err)
	} else {
		defer func() {
			_ = fileObj.Close()
		}()
		if content, err := ioutil.ReadAll(fileObj); err != nil {
			return nil, types.NewException(TMFile_CanNotRead_File, map[string]string{
				"path": path,
			}, nil)
		} else {
			return content, nil
		}
	}
}
