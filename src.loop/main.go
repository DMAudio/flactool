package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var startTime = strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)

func FileWriteString(fPath string, str string) error {
	var absPath string
	var err error
	if absPath, err = filepath.Abs(strings.TrimSpace(fPath)); err != nil {
		return fmt.Errorf("无法解析文件绝对路径")
	}

	if _, err := os.Stat(path.Dir(fPath)); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(fPath), 0777); err != nil && !os.IsExist(err) {
			return fmt.Errorf("无法创建日志目录：%s\n%v", path.Dir(fPath), err)
		}
	}

	if file, err := os.Create(absPath); err != nil {
		return fmt.Errorf("无法创建文件")
	} else {
		defer func() { _ = file.Close() }()

		if _, err := file.WriteString(str); err != nil {
			return fmt.Errorf("无法写入文件")
		}
	}
	return nil
}

func fixEnv(cmd *exec.Cmd, extra map[string]string) {
	env := os.Environ()
	cmdEnv := make([]string, 0)

	if extra == nil {
		extra = map[string]string{}
	}

	for _, e := range env {
		i := strings.Index(e, "=")
		if i > 0 {
			if _, existInExtra := extra[e[:i]]; !existInExtra {
				cmdEnv = append(cmdEnv, e)
			}
		}
	}

	for envKey, envValue := range extra {
		cmdEnv = append(cmdEnv, fmt.Sprintf("%s=%s", envKey, envValue))
	}

	cmd.Env = cmdEnv

}

func main() {
	exe := flag.String("exe", "", "可执行文件路径")
	task := flag.String("task", "", "配置文件路径")
	logDir := flag.String("logDir", "", "日志目录")
	inputDir := flag.String("inputDir", "", "源文件目录")
	outputDir := flag.String("outputDir", "", "输出目的目录")
	outputTemplate := flag.String("outputTemplate", "", "输出文件名格式, 原文件名使用%filename%代替")
	flag.Parse()

	if logDir == nil {
		logDir = new(string)
	}
	if *logDir == "" {
		*logDir = "./logs/" + startTime
	} else {
		*logDir = path.Join(*logDir, startTime)
	}
	if absPath, err := filepath.Abs(strings.TrimSpace(*logDir)); err != nil {
		panic("无法解析日志文件夹的绝对路径")
	} else {
		*logDir = absPath
	}

	if exe == nil || *exe == "" {
		panic("请指定可执行文件路径")
	} else if exeParsed, err := filepath.Abs(filepath.FromSlash(strings.TrimSpace(*exe))); err != nil {
		panic("无法解析可执行文件的绝对路径")
	} else {
		*exe = exeParsed
	}

	if task == nil || *task == "" {
		panic("请指定任务文件路径")
	} else if taskParsed, err := filepath.Abs(filepath.FromSlash(strings.TrimSpace(*task))); err != nil {
		panic("无法解析任务文件的绝对路径")
	} else {
		*task = taskParsed
	}

	if inputDir == nil || *inputDir == "" {
		panic("请指定源文件目录")
	} else if inputDirParsed, err := filepath.Abs(filepath.FromSlash(strings.TrimSpace(*inputDir))); err != nil {
		panic("无法解析源文件目录的绝对路径")
	} else {
		*inputDir = inputDirParsed
	}

	if outputDir == nil || *outputDir == "" {
		panic("请指定输出目的目录")
	} else if outputDirParsed, err := filepath.Abs(filepath.FromSlash(strings.TrimSpace(*outputDir))); err != nil {
		panic("无法解析输出目的目录的绝对路径")
	} else {
		*outputDir = outputDirParsed
	}

	if files, err := ioutil.ReadDir(*inputDir); err != nil {
		panic(fmt.Errorf("无法遍历源文件目录:%v", err))
	} else {
		taskAmount := 0
		taskFailed := make([]string, 0)
		for _, f := range files {
			inputFile := filepath.ToSlash(path.Join(*inputDir, f.Name()))
			if path.Ext(inputFile) != ".flac" {
				continue
			}

			var outputFileName string
			if outputTemplate == nil || strings.TrimSpace(*outputTemplate) == "" {
				outputFileName = f.Name()
			} else {
				outputFileName = strings.TrimSpace(*outputTemplate)
			}
			outputFile := filepath.ToSlash(path.Join(*outputDir, outputFileName))

			command := []string{
				filepath.ToSlash(*exe),
				"-input", inputFile,
				"-task", filepath.ToSlash(*task),
				"-output", outputFile,
				"-trace", "48",
			}
			cmd := exec.Command(command[0], command[1:]...)
			fixEnv(cmd, map[string]string{
				"inputDir":             path.Dir(inputFile),
				"inputBase":            path.Base(inputFile),
				"inputExt":             strings.TrimPrefix(path.Ext(inputFile), "."),
				"inputBaseTrimmedExt":  strings.TrimSuffix(path.Base(inputFile), path.Ext(inputFile)),
				"outputDir":            path.Dir(outputFile),
				"outputBase":           path.Base(outputFile),
				"outputExt":            strings.TrimPrefix(path.Ext(outputFile), "."),
				"outputBaseTrimmedExt": strings.TrimSuffix(path.Base(outputFile), path.Ext(outputFile)),
			})

			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err = cmd.Run()

			outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

			fmt.Println("\n==================================")
			fmt.Println(strings.Join(command, " "))
			taskAmount += 1

			if err != nil {
				message := fmt.Sprintf("无法执行任务: %v", err)
				fmt.Println(message)
				errStr += message
			} else {
				fmt.Println(outStr)
			}

			errStr = strings.TrimSpace(errStr)
			if errStr != "" {
				taskFailed = append(taskFailed, path.Base(inputFile))

				logBase := path.Base(inputFile) + "_err.log"
				logBase = strings.ReplaceAll(logBase, "\u003A", "\uFF1A")
				logBase = strings.ReplaceAll(logBase, "\u002F", "\uFF0F")
				logBase = strings.ReplaceAll(logBase, "\u005C", "\uFF3C")
				logBase = strings.ReplaceAll(logBase, "\u003F", "\uFF1F")
				logBase = strings.ReplaceAll(logBase, "\u0022", "\u0027\u0027")
				logBase = strings.ReplaceAll(logBase, "\u002A", "\uFF0A")
				logBase = strings.ReplaceAll(logBase, "\u003C", "\uFF1C")
				logBase = strings.ReplaceAll(logBase, "\u003E", "\uFF1E")
				logBase = strings.ReplaceAll(logBase, "\u007C", "\uFF5C")

				logPath := path.Join(*logDir, logBase)
				if err := FileWriteString(logPath, errStr); err != nil {
					panic(fmt.Errorf("无法创建日志文件：%s\n%v", logPath, err))
				}
			}
		}

		fmt.Printf("共处理了 %d 个文件，%d 个失败\n", taskAmount, len(taskFailed))
		if len(taskFailed) > 0 {
			fmt.Printf("以下文件处理失败：(日志文件路径：%s)\n", *logDir)
			for _, fileBase := range taskFailed {
				fmt.Printf("\t%s\n", fileBase)
			}
		}

	}
}
