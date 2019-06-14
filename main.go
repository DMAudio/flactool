package main

import (
	"flag"
	"gitlab.com/KTGWKenta/DADP.FlacTool/config"
	"gitlab.com/KTGWKenta/DADP.FlacTool/format/flac"
	"gitlab.com/KTGWKenta/DADP.FlacTool/task"
	"gitlab.com/MGEs/Com.Base/types"
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
	task.Init()
	flac.Init()
	InitMain()

	config.ParseFlags()

	fObj := &flac.Flac{}
	if *inputFile != "" {
		if err := fObj.ParseFromFile(*inputFile); err != nil {
			types.Throw(err, types.RsPanic)
		}
	} else {
		panic("未指定源文件")
	}

	if *taskFile == "" {
		panic("未指定任务配置文件")
	} else if taskList, err := task.LoadTaskList(*taskFile); err != nil {
		types.Throw(err, types.RsError)
	} else if err := taskList.ExecuteTasks(map[string]interface{}{
		"flac": fObj,
	}); err != nil {
		types.Throw(err, types.RsError)
	} else if *outputFile != "" {
		if outputFileProcessed, _, err := task.GlobalArgFiller().FillArgs(*outputFile, flac.ToArgFillerParameter(fObj)); err != nil {
			types.Throw(err, types.RsError)
		} else {
			fObj.WriteToFile(outputFileProcessed)
		}
	}
}
