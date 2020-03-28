package utils

import (
	"runtime"
)

func PrintFuncName() string {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(2, fpcs)
	if n == 0 {
		InfoLogger.Logln("---PrintFuncName--- : MSG: NO CALLER")
		return ""
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		InfoLogger.Logln("---PrintFuncName--- : MSG CALLER WAS NIL")
		return ""
	}

	return caller.Name()
}
