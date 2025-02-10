package main

import (
	"fmt"
	"time"
)

func logInfo(msg string, arg ...any) {
	if len(arg) > 0 {
		fmt.Printf("%s - "+msg+"\n", time.Now().Format(time.DateTime), arg)
	} else {
		fmt.Printf("%s - "+msg+"\n", time.Now().Format(time.DateTime))
	}
}

func logDebug(msg string) {
	if debugLogs {
		logInfo(msg)
	}
}
