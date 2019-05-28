package task

import (
	"dadp.flactool/types"
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

var TMTask_FailedTo_Parse_WhenPattern = types.NewMask(
	"FailedTo_Parse_WhenPattern",
	"无法解析开关表达式:\n{{pattern}}",
)

var TMTask_UnableToParse_SubArg = types.NewMask(
	"UNABLE_TO_PARSE_SubARG",
	"子参数解析失败",
)

var TMTask_UnableToExecute_SubTask = types.NewMask(
	"UNABLE_TO_Execute_SubTask",
	"子任务执行失败",
)

var TMFiller_FailedTo_CompileRegex = types.NewMask(
	"FailedTo_CompileRegex",
	"无法编译正则表达式",
)

var TMFiller_UnableToParse_Arg = types.NewMask(
	"UNABLE_TO_PARSE_ARG",
	"参数解析失败：{{reason}}",
)

var TMFiller_Undefined_Handler = types.NewMask(
	"Undefined_Handler",
	"未定义的处理者：{{handler}}",
)

var TMFiller_FailedToExecute_Filler = types.NewMask(
	"Failed_To_Execute_Filler",
	"Filler执行失败（执行者：{{handler}}，参数：{{parameters}}）",
)

var TMFiller_FailedToFill_Args = types.NewMask(
	"Failed_To_Fill_Args",
	"参数填充失败",
)

var TMFiller_CanNotRegister = types.NewMask(
	"CAN_NOT_REGISTER_Filler",
	"无法注册执行者: {{reason}}",
)
