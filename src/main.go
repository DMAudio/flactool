package main

import (
	"flag"
	"p20190417/config"
	"p20190417/format/flac"
	"p20190417/task"
	"p20190417/types"
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
