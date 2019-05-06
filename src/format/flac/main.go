package flac

import (
	"p20190417/task"
	"p20190417/types"
)

var TMFlac_CanNotRegister_TaskHandler = types.NewMask(
	"CAN_NOT_REGISTER_TASK_HANDLER",
	"无法注册任务执行者：{{handler}}",
)

func Init() *types.Exception {
	registerFailedErr := types.NewException(TMFlac_CanNotRegister_TaskHandler, map[string]string{}, nil)

	if err := task.GlobalHandler().Register(TaskHandler_T4VORB_Key, TaskHandler_T4VORB); err != nil {
		return registerFailedErr.SetCause(err).SetParam("handler", TaskHandler_T4VORB_Key)
	}

	return nil
}
