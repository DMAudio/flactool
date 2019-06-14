package task

import (
	"gitlab.com/MGEs/Com.Base/types"
)

func Init() *types.Exception {
	if err := GlobalHandler().Register("COMMON", Handler_COMMON); err != nil {
		return err
	}

	if err := GlobalArgFiller().Register("env", Filler_Env); err != nil {
		return err
	}

	if err := GlobalArgFiller().Register("fmtFName", Filler_FmtFileName); err != nil {
		return err
	}

	if err := GlobalArgFiller().Register("u", Filler_DecodeURI); err != nil {
		return err
	}

	return nil
}
