package main

import (
	"flag"
	"p20190417/config"
	"p20190417/format/flac"
	"p20190417/task"
	"p20190417/util"
)

var inputFile *string
var outputFile *string

func InitMain() {
	inputFile = flag.String("input", "", "")
	outputFile = flag.String("output", "", "")
}

func main() {
	util.InitThrowableExt()
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
			util.Throw(err, util.RsError)
		}
	} else {
		panic("未指定任务配置文件")
	}

	if *outputFile != "" {
		flac.GlobalFlac().WriteToFile(*outputFile)
	}
}
