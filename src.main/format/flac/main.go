package flac

import (
	"dadp.flactool/task"
	"dadp.flactool/types"
)

func Init() *types.Exception {
	if err := task.GlobalHandler().Register(TaskHandler_T4VORB_Key, TaskHandler_T4VORB); err != nil {
		return err
	}

	if err := task.GlobalHandler().Register(TaskHandler_T6PICT_Key, TaskHandler_T6PICT); err != nil {
		return err
	}

	if err := task.GlobalHandler().Register(TaskHandler_BLOCKS_Key, TaskHandler_BLOCKS); err != nil {
		return err
	}

	if err := task.GlobalArgFilter().Register("flac", ArgFilter); err != nil {
		return err
	}

	return nil
}
