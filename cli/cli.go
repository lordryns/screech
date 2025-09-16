package cli

import "os"

func GetArg(index int) string {
	var args = os.Args
	if len(args) > index {
		return args[index]
	}

	return ""
}
