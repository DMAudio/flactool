package main

import (
	"dadp.flactool/config"
	"dadp.flactool/format/flac"
	"dadp.flactool/task"
	"dadp.flactool/types"
	"flag"
)

var inputFile *string
var outputFile *string

func InitMain() {
	inputFile = flag.String("input", "", "")
	outputFile = flag.String("output", "", "")
}

func main() {
	types.InitThrowableExt()
	task.Init()
	flac.Init()
	InitMain()

	config.ParseFlags()

	if *inputFile != "" {
		flac.GlobalFlac().ParseFromFile(*inputFile)
	} else {
		panic("未指定源文件")
	}

	if *task.CollectionFile != "" {
		if err := task.ExecuteTasks(); err != nil {
			types.Throw(err, types.RsError)
		}
	} else {
		panic("未指定任务配置文件")
	}

	if *outputFile != "" {
		if outputFileProcessed, _, err := task.GlobalArgFilter().FillArgs(*outputFile, nil); err != nil {
			types.Throw(err, types.RsError)
		} else {
			flac.GlobalFlac().WriteToFile(outputFileProcessed)
		}
	}
}
