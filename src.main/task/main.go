package task

import "dadp.flactool/types"

func Init() *types.Exception {
	if err := GlobalHandler().Register("COMMON", Handler_COMMON); err != nil {
		return err
	}

	if err := GlobalArgFiller().Register("env", Filter_Env); err != nil {
		return err
	}

	if err := GlobalArgFiller().Register("fmtFName", Filter_FmtFileName); err != nil {
		return err
	}

	if err := GlobalArgFiller().Register("u", Filter_DecodeURI); err != nil {
		return err
	}

	return nil
}
