package config

import (
	"p20190417/types"
	"sync"
)

var globalConfig *types.SSMap
var globalConfigLock sync.Mutex

func GlobalConfig() *types.SSMap {
	if globalConfig == nil {
		globalConfigLock.Lock()
		defer globalConfigLock.Unlock()
		if globalConfig == nil {
			globalConfig = &types.SSMap{}
		}
	}
	return globalConfig
}

