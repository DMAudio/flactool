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
var taskFile *string

func InitMain() {
	inputFile = flag.String("input", "", "path to the source file")
	outputFile = flag.String("output", "", "path which task result shall be saved to")
	taskFile = flag.String("task", "", "path to the task list file")
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

	if *taskFile == "" {
		panic("未指定任务配置文件")
	} else if taskList, err := task.LoadTaskList(*taskFile); err != nil {
		types.Throw(err, types.RsError)
	} else if err := taskList.ExecuteTasks(); err != nil {
		types.Throw(err, types.RsError)
	} else if *outputFile != "" {
		if outputFileProcessed, _, err := task.GlobalArgFilter().FillArgs(*outputFile, nil); err != nil {
			types.Throw(err, types.RsError)
		} else {
			flac.GlobalFlac().WriteToFile(outputFileProcessed)
		}
	}
}
