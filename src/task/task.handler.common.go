package task

import (
	"bytes"
	"fmt"
	"os/exec"
	"p20190417/types"
	"strings"
)

func Handler_COMMON(operation string, args interface{}) (interface{}, *types.Exception) {
	switch operation {
	case "execute":
		return Handler_COMMON_Cmd(args)
	default:
		return nil, types.NewException(TMConfig_Unsupported_TaskOperation, map[string]string{
			"operation": operation,
		}, nil)
	}
}

func Handler_COMMON_Cmd(args interface{}) (interface{}, *types.Exception) {
	var argsRaw []string
	var exception *types.Exception
	var err error
	if argsRaw, exception = types.InterfaceToStringSlice(args); exception != nil {
		return nil, exception
	}

	for argIndex, argRaw := range argsRaw {
		if argParsed, _, err := GlobalArgFilter().FillArgs(argRaw, nil); err != nil {
			return nil, err
		} else {
			argsRaw[argIndex] = argParsed
		}
	}

	if len(argsRaw) == 0 {
		return nil, types.NewException(types.NewMask(
			"UnableTo_Execute_Command",
			fmt.Sprintf("命令执行失败：%s", "未指定可执行文件路径"),
		), nil, err)
	}

	cmd := exec.Command(argsRaw[0], argsRaw[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		return nil, types.NewException(types.NewMask(
			"UnableTo_Execute_Command",
			fmt.Sprintf("命令执行失败：\n$> %s\n%s", strings.Join(argsRaw, " "), errStr),
		), nil, err)
	} else {
		return outStr, nil
	}

}
