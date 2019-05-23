package config

import "flag"

var globalFlagRegisters []func()

func ParseFlags() {
	for _, rFunc := range globalFlagRegisters {
		rFunc()
	}
	flag.Parse()
}

func NewRegister(rFunc func()) {
	globalFlagRegisters = append(globalFlagRegisters, rFunc)
}
