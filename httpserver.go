package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	hasBeenCalled bool
)

func handler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		fmt.Println("❌  HTTP handler error:", err)
		os.Exit(1)
	}
	if !hasBeenCalled {
		logInfo("✅  Connection from game successfully established")
		hasBeenCalled = true
	}
	logDebug("received GET")
	lastRequestTime = time.Now()
}
