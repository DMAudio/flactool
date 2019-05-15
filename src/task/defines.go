package task

import (
	"p20190417/types"
)

var TMConfig_UnableToRead_TaskFile = types.NewMask(
	"UNABLE_TO_READ_TASK_FILE",
	"无法读取任务配置",
)

var TMTask_CanNotRegister = types.NewMask(
	"CAN_NOT_REGISTER_TASK",
	"无法注册执行者: {{reason}}",
)

var TMTask_CanNotExecute = types.NewMask(
	"CAN_NOT_EXECUTE_TASK",
	"无法执行任务: {{reason}}",
)

var TMConfig_UnableToParse_TaskCollection = types.NewMask(
	"UNABLE_TO_PARSE_TASK_COLLECTION",
	"无法解析第 {{collection}} 个任务集",
)

var TMConfig_UnableToParse_TaskItem = types.NewMask(
	"UNABLE_TO_PARSE_TASK_ITEM",
	"无法解析第 {{task}} 个任务",
)

var TMConfig_Unsupported_TaskOperation = types.NewMask(
	"UNSUPPORTED_TASK_OPERATION",
	"操作名 {{operation}} 不被支持",
)

var TMConfig_UnableToParse_TaskFile = types.NewMask(
	"UNABLE_TO_PARSE_TASK_FILE",
	"无法解析任务文件：{{reason}}",
)

var TMTask_Collection_UnhandledThrowable = types.NewMask(
	"UNHANDLED_THROWABLE_OCCURRED",
	"执行 第 {{collection}} 个任务集时产生了若干个未被处理的错误(执行者: {{handler}}):",
)

var TMTask_UnhandledThrowable = types.NewMask(
	"UNHANDLED_THROWABLE_OCCURRED",
	"执行 任务配置文件 时产生了若干个未被处理的错误:",
)
