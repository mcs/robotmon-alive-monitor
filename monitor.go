package main

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func startProcessIfNotRunning() {
	logInfo("Title = %s", processTitle)
	processID = findByTitle(processTitle)
	if processID == 0 {
		processID, err := startCommand(processCommand)
		if err != nil {
			log.Fatal("❌  Error starting process:", err)
		}
		logInfo("Process started with PID %s", strconv.Itoa(int(processID)))
	}
	lastRequestTime = time.Now()
}

func monitorAndRestartProcess() {
	for {
		time.Sleep(1 * time.Minute)
		logDebug("lastRequestTime = " + lastRequestTime.String())

		if time.Since(lastRequestTime) > restartThreshold {
			logInfo("❗  No requests were received for more than %s. Restarting process...", restartThreshold.String())

			err1 := kill(processID)
			if err1 != nil {
				if strings.TrimSpace(processTitle) != "" {
					// lookup existing running process, maybe Pid got messed up due to manual user restart of the original process
					pId := findByTitle(processTitle)
					if pId > 0 && pId != processID {
						// the process has a different process id than expected
						err2 := kill(pId)
						if err2 != nil {
							logInfo("❌  Could neither kill the originally spawned process nor the one found by title")
							log.Printf("  PID 1: %d, Error 1: %s\n", processID, err1)
							log.Printf("  PID 2: %d, Error 2: %s\n", pId, err2)
						}
						// situation rescued
					} else {
						logInfo("❌  Could not kill the originally spawned process and did not find existing one by title")
						log.Printf("  PID: %d, Error: %s\n", processID, err1)
					}
				} else {
					logInfo("❌  Could not kill the originally spawned process")
				}
			} else {
				// nice one
			}

			// Wait 10 seconds
			time.Sleep(10 * time.Second)

			// Start process
			startProcessIfNotRunning()

			lastRequestTime = time.Now()
		} else if time.Since(lastRequestTime) > restartThreshold-1*time.Minute {
			logInfo("⚠️ WARNING: No requests were received for more than %s. Restarting process in one minute if still no request happened...", (restartThreshold - 1*time.Minute).String())
		}
	}
}
