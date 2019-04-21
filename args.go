package main

import (
	"os"
	"strings"
)

func checkArg(arg string, handler func(bool, string)) bool {
	args := os.Args
	for _, v := range args {
		if strings.ToLower(v) == strings.ToLower(arg) {
			argTmp := strings.Split(v, arg)
			if len(argTmp) > 1 {
				handler(true, argTmp[1])
				return true
			}
			handler(true, "")
			return true
		}
	}

	handler(false, "Not Found")
	return false
}
