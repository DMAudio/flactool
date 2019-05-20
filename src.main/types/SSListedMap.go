package types

import (
	"strconv"
)

type SSListedMapAction int8

const (
	SSLM_Exception SSListedMapAction = 0
	SSLM_Prepend   SSListedMapAction = 1
	SSLM_Append    SSListedMapAction = 2
)

type SSListedMap struct {
	valMap  map[string]string
	keyList []string
	//rwLock  sync.RWMutex
}

var TMSSLM_Nil_Pointer = NewMask(
	"NIL_POINTER",
	"空指针",
)

var TMSSLM_KeyList_Broken = NewMask(
	"KEYLIST_BROKEN",
	"键顺序表损坏",
)

var TMSSLM_ValMap_Broken = NewMask(
	"VALMAP_BROKEN",
	"键值表损坏",
)

var TMSSLM_Key_Existed = NewMask(
	"KEY_EXISTED",
	"键 {{key}} 已存在",
)

var TMSSLM_Key_NotFound = NewMask(
	"KEY_NOTFOUND",
	"键 {{key}} 不存在",
)

var TMSSLM_Index_OutOfRange = NewMask(
	"INDEX_OUTOFRANGE",
	"下标越界：{{index}}",
)

func NewSSListedMap() *SSListedMap {
	s := &SSListedMap{}
	if err := s.makePrepared(); err != nil {
		panic(err.GetMessage(true))
	}
	return s
}

func (s *SSListedMap) doesPrepared() bool {
	return s != nil && s.keyList != nil && s.valMap != nil
}

func (s *SSListedMap) makePrepared() *Exception {
	if s == nil {
		return NewException(TMSSLM_Nil_Pointer, nil, nil)
	}
	if !s.doesPrepared() {
		s.valMap = map[string]string{}
		s.keyList = []string{}
	}
	return nil
}

func (s *SSListedMap) Length() int {
	if !s.doesPrepared() {
		return 0
	}
	return len(s.valMap)
}

func (s *SSListedMap) Locate(key string) (int, *Exception) {
	if !s.doesPrepared() {
		return -1, nil
	}
	if _, hasKey := s.valMap[key]; !hasKey {
		return -1, nil
	}

	for pos, elKey := range s.keyList {
		if key == elKey {
			return pos, nil
		}
	}

	return -2, NewException(TMSSLM_KeyList_Broken, nil, nil)
}

func (s *SSListedMap) Prepend(elKey string, elValue string) *Exception {
	_, e := s.InsertBeforeIndex(0, elKey, elValue)
	return e
}

func (s *SSListedMap) Append(elKey string, elValue string) (int, *Exception) {
	return s.InsertAfterIndex(len(s.valMap)-1, elKey, elValue)
}

func (s *SSListedMap) InsertBeforeKey(key string, elKey string, elValue string) (int, *Exception) {
	if keyPos, err := s.Locate(key); err != nil {
		return -1, err
	} else if keyPos == -1 {
		return -1, NewException(TMSSLM_Key_NotFound, nil, nil)
	} else {
		return s.InsertBeforeIndex(keyPos, elKey, elValue)
	}
}

func (s *SSListedMap) InsertAfterKey(key string, elKey string, elValue string) (int, *Exception) {
	if keyPos, err := s.Locate(key); err != nil {
		return -1, err
	} else if keyPos == -1 {
		return -1, NewException(TMSSLM_Key_NotFound, nil, nil)
	} else {
		return s.InsertAfterIndex(keyPos, elKey, elValue)
	}
}

func (s *SSListedMap) InsertBeforeIndex(index int, elKey string, elValue string) (int, *Exception) {
	if err := s.makePrepared(); err != nil {
		return -1, err
	}

	if _, existed := s.valMap[elKey]; existed {
		return -1, NewException(TMSSLM_Key_Existed, map[string]string{
			"key": elKey,
		}, nil)
	}

	if index == 0 {
		s.valMap[elKey] = elValue
		s.keyList = append([]string{elKey}, s.keyList...)
		return 0, nil
	}

	if index < len(s.valMap) {
		s.valMap[elKey] = elValue
		s.keyList = SListInsertBefore(s.keyList, index, elKey)
		return index, nil
	}

	return -1, NewException(TMSSLM_Index_OutOfRange, map[string]string{
		"index": strconv.Itoa(index),
	}, nil)

}

func (s *SSListedMap) InsertAfterIndex(index int, elKey string, elValue string) (int, *Exception) {
	if err := s.makePrepared(); err != nil {
		return -1, err
	}

	if _, existed := s.valMap[elKey]; existed {
		return -1, NewException(TMSSLM_Key_Existed, map[string]string{
			"key": elKey,
		}, nil)
	}

	if index == len(s.valMap)-1 {
		s.valMap[elKey] = elValue
		s.keyList = append(s.keyList, elKey)
		return 0, nil
	}

	if index >= 0 && index < len(s.valMap)-1 {
		s.valMap[elKey] = elValue
		s.keyList = SListInsertAfter(s.keyList, index, elKey)
		return index, nil
	}

	return -1, NewException(TMSSLM_Index_OutOfRange, map[string]string{
		"index": strconv.Itoa(index),
	}, nil)
}

func (s *SSListedMap) Get(Key string) (string, *Exception) {
	if !s.doesPrepared() {
		return "", NewException(TMSSLM_Key_NotFound, map[string]string{
			"key": Key,
		}, nil)
	}

	if value, exist := s.valMap[Key]; !exist {
		return "", NewException(TMSSLM_Key_NotFound, map[string]string{
			"key": Key,
		}, nil)
	} else {
		return value, nil
	}

}

func (s *SSListedMap) Set(Key string, Value string, actionIfNotExist SSListedMapAction) *Exception {
	if _, exist := s.valMap[Key]; !exist {
		switch actionIfNotExist {
		case SSLM_Prepend:
			return s.Prepend(Key, Value)
		case SSLM_Append:
			_, e := s.Append(Key, Value)
			return e
		default:
			return NewException(TMSSLM_Key_NotFound, map[string]string{
				"key": Key,
			}, nil)

		}

	} else {
		s.valMap[Key] = Value
		return nil
	}
}

func (s *SSListedMap) Delete(Key string) *Exception {
	if _, exist := s.valMap[Key]; exist {
		if keyIndex, err := s.Locate(Key); err != nil {
			return err
		} else {
			delete(s.valMap, Key)
			s.keyList = append(s.keyList[:keyIndex], s.keyList[keyIndex+1:]...)
		}
	}
	return nil
}

func (s *SSListedMap) DumpList() ([][2]string, *Exception) {
	if err := s.makePrepared(); err != nil {
		return nil, err
	}

	result := make([][2]string, s.Length())
	for keyIndex, key := range s.keyList {
		if value, exist := s.valMap[key]; !exist {
			return nil, NewException(TMSSLM_ValMap_Broken, nil, nil)
		} else {
			result[keyIndex] = [2]string{key, value}
		}
	}

	return result, nil
}

func (s *SSListedMap) DumpMap() (map[string]string, *Exception) {
	if err := s.makePrepared(); err != nil {
		return nil, err
	}
	result := map[string]string{}
	for _, key := range s.keyList {
		if value, exist := s.valMap[key]; !exist {
			return nil, NewException(TMSSLM_ValMap_Broken, nil, nil)
		} else {
			result[key] = value
		}
	}
	return result, nil
}

func (s *SSListedMap) Import(input [][2]string) *Exception {
	keyList := make([]string, len(input))
	valMap := map[string]string{}

	for recordIndex, record := range input {
		if _, keyExisted := valMap[record[0]]; keyExisted {
			return NewException(TMSSLM_Key_Existed, map[string]string{
				"key": record[0],
			}, nil)
		}
		keyList[recordIndex] = record[0]
		valMap[record[0]] = record[1]
	}

	s.valMap = valMap
	s.keyList = keyList

	return nil
}

func (s *SSListedMap) Sort(sortBy []string) {
	if !s.doesPrepared() {
		return
	}

	fillOtherTags := -1
	unrelatedTags := make([]string, len(s.keyList))
	copy(unrelatedTags, s.keyList)

	result := make([]string, 0)
	for _, key := range sortBy {
		if _, keyExist := s.valMap[key]; !keyExist {
			if key == "..." {
				fillOtherTags = len(result) - 1
			}
			continue
		}
		result = append(result, key)
		unrelatedTags = SListDeleteByElement(unrelatedTags, key)
	}

	if fillOtherTags > -1 {
		result = SListInsertAfter(result, fillOtherTags, unrelatedTags...)
	} else {
		newValMap := map[string]string{}
		for _, key := range result {
			newValMap[key] = s.valMap[key]
		}
		s.valMap = newValMap
	}

	s.keyList = result
}
