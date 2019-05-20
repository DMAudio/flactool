package types

import (
	"fmt"
	"sync"
)

type SSMap struct {
	data map[string]string
	lock sync.RWMutex
}

func (s SSMap) Has(key string) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, isExist := s.data[key]

	return isExist
}

func (s SSMap) Get(key string) string {
	if !s.Has(key) {
		return ""
	}

	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data[key]
}

func (s *SSMap) Set(key string, value string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.data == nil {
		s.data = make(map[string]string)
	}

	s.data[key] = value
}

func (s SSMap) Del(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, key)
}

func (s SSMap) Each(walker func(key string, value string)) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	for key, value := range s.data {
		walker(key, value)
	}
}

func (s *SSMap) From(input map[string]string) *SSMap {
	s.mergeSSMap(input)
	return s
}

func (s SSMap) Dump() map[string]string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	output := map[string]string{}
	for key, value := range s.data {
		output[key] = value
	}

	return output
}

var (
	ERR_Unexpected_Input_Type = fmt.Errorf("Unexpected_Input_Type")
)

func (s *SSMap) Merge(input interface{}) error {
	switch input.(type) {
	case SSMap:
		inputObj := input.(SSMap)
		s.mergeSelfPtr(&inputObj)
	case *SSMap:
		s.mergeSelfPtr(input.(*SSMap))
	case map[string]string:
		s.mergeSSMap(input.(map[string]string))
	case *map[string]string:
		inputMap := input.(*map[string]string)
		s.mergeSSMap(*inputMap)
	default:
		return ERR_Unexpected_Input_Type
	}
	return nil
}

func (s *SSMap) mergeSelfPtr(input *SSMap) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.data == nil {
		s.data = make(map[string]string)
	}

	input.Each(func(key string, value string) {
		s.data[key] = value
	})
}

func (s *SSMap) mergeSSMap(input map[string]string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.data == nil {
		s.data = make(map[string]string)
	}

	for key, value := range input {
		s.data[key] = value
	}
}
