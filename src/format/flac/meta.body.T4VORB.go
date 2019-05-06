package flac

import (
	"encoding/binary"
	"fmt"
	"p20190417/types"
	"p20190417/util"
	"strings"
)

type MetaBlockT4VORB struct {
	refer    string
	comments [][2]string
}

func (mb *MetaBlockT4VORB) Parse(r *util.BinaryReader) *types.Exception {
	var referSize uint32
	if referSizeData, err := r.ReadBytes(4); err == nil {
		referSize = binary.LittleEndian.Uint32(referSizeData)
	} else {
		return types.NewException(TMFlac_CanNotParse_MetaT4ReferSize, nil, err)
	}

	if referData, err := r.ReadBytes(uint64(referSize)); err == nil {
		mb.refer = string(referData)
	} else {
		return types.NewException(TMFlac_CanNotParse_MetaT4ReferData, nil, err)
	}

	var commentListLength uint32
	if commentListLengthData, err := r.ReadBytes(4); err == nil {
		commentListLength = binary.LittleEndian.Uint32(commentListLengthData)
	} else {
		return types.NewException(TMFlac_CanNotParse_MetaT4CommentAmount, nil, err)
	}

	mb.comments = make([][2]string, commentListLength)

	for i := uint32(0); i < commentListLength; i++ {
		var commentLength uint32
		if commentLengthData, err := r.ReadBytes(4); err == nil {
			commentLength = binary.LittleEndian.Uint32(commentLengthData)
		} else {
			return types.NewException(TMFlac_CanNotParse_MetaT4CommentItemLength, nil, err)
		}

		if commentData, err := r.ReadBytes(uint64(commentLength)); err == nil {
			commentTemp := strings.SplitN(string(commentData), "=", 2)
			mb.comments[i][0] = commentTemp[0]
			mb.comments[i][1] = commentTemp[1]
		} else {
			return types.NewException(TMFlac_CanNotParse_MetaT4CommentData, nil, err)
		}
	}

	return nil
}

func (mb *MetaBlockT4VORB) Encode() (*types.Buffer, *types.Exception) {
	buffer := types.NewBuffer()

	referSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(referSize, uint32(len(mb.refer)))
	if length, err := buffer.Write(referSize); err != nil || length != len(referSize) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT4ReferSize, nil, err)
	}

	if length, err := buffer.WriteString(mb.refer); err != nil || length != len(mb.refer) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT4ReferData, nil, err)
	}

	commentAmount := make([]byte, 4)
	binary.LittleEndian.PutUint32(commentAmount, uint32(len(mb.comments)))
	if length, err := buffer.Write(commentAmount); err != nil || length != len(commentAmount) {
		return nil, types.NewException(TMFlac_CanNotWrite_MetaT4CommentAmount, nil, err)
	}

	for _, commentItem := range mb.comments {
		commentContent := fmt.Sprintf("%s=%s", commentItem[0], commentItem[1])

		commentItemLength := make([]byte, 4)
		binary.LittleEndian.PutUint32(commentItemLength, uint32(len(commentContent)))
		if length, err := buffer.Write(commentItemLength); err != nil || length != len(commentItemLength) {
			return nil, types.NewException(TMFlac_CanNotWrite_MetaT4CommentAmount, nil, err)
		}

		if length, err := buffer.WriteString(commentContent); err != nil || length != len(commentContent) {
			return nil, types.NewException(TMFlac_CanNotWrite_MetaT4CommentItemLength, nil, err)
		}
	}

	return buffer, nil
}

func (mb *MetaBlockT4VORB) GetRefer() string {
	return mb.refer
}

func (mb *MetaBlockT4VORB) SetRefer(referText string) {
	mb.refer = referText
}

func (mb *MetaBlockT4VORB) GetCommentList() [][2]string {
	comments := make([][2]string, len(mb.comments))
	copy(comments, mb.comments)
	return comments
}

func (mb *MetaBlockT4VORB) GetCommentMap() map[string]string {
	comments := map[string]string{}
	for _, commentItem := range mb.comments {
		comments[commentItem[0]] = commentItem[1]
	}

	return comments
}

func (mb *MetaBlockT4VORB) SetComment(key, Value string) {
	found := false
	for i := 0; i < len(mb.comments); i++ {
		if mb.comments[i][0] == key {
			mb.comments[i][1] = Value
			found = true
		}
	}
	if !found {
		mb.comments = append(mb.comments, [2]string{key, Value})
	}
}

func (mb *MetaBlockT4VORB) SetComments(commentList [][2]string) {
	comments := make([][2]string, len(commentList))
	copy(comments, commentList)
	mb.comments = comments
}

func (mb *MetaBlockT4VORB) DeleteComment(key string) {
	for i := 0; i < len(mb.comments); i++ {
		if mb.comments[i][0] == key {
			mb.comments = append(mb.comments[:i], mb.comments[i+1:]...)
		}
	}
}
