package utils

import "os"

const (
	__start_arg_test     = "--runCase"
	__start_arg_updateDB = "--updateDb"
)

func hasParameter(para string) bool {
	for _, arg := range os.Args {
		if arg == para {
			return true
		}
	}
	return false
}

func ShouldRunTest() bool {
	return hasParameter(__start_arg_test)
}

func ShouldUpdateDB() bool {
	return hasParameter(__start_arg_updateDB)
}
