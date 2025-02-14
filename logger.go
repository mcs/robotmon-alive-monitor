package main

import (
	"fmt"
	"time"
)

func logInfo(msg string, arg ...any) {
	args := []any{time.Now().Format(time.DateTime)}
	if len(arg) > 0 {
		args = append(args, arg)
	}
	fmt.Printf("%s - "+msg+"\n", args...)
}

func logDebug(msg string) {
	if debugLogs {
		logInfo(msg)
	}
}
