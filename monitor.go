package main

import (
	"log"
	"strconv"
	"strings"
	"time"
)

func startProcessIfNotRunning() {
	if strings.TrimSpace(processTitle) != "" {
		logInfo("Title = %s", processTitle)
		processID = findByTitle(processTitle)
	}
	if processID == 0 {
		var err error
		processID, err = startCommand(processCommand)
		if err != nil {
			log.Fatal("❌  Error starting process: ", err)
		}
		logInfo("ℹ️ Process started with PID %s", strconv.Itoa(int(processID)))
	} else {
		logInfo("ℹ️ Found already running process with PID %d. Not starting a new process", processID)
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
							logInfo("❌  Could neither kill the originally spawned pid %d nor pid %d found by title", processID, pId)
							log.Printf("  PID %d error: %s\n", processID, err1)
							log.Printf("  PID %d error: %s\n", pId, err2)
						}
					} else {
						logInfo("❌  Could not kill the originally spawned process and did not find existing one by title")
						log.Printf("  PID: %d, Error: %s\n", processID, err1)
					}
				} else {
					logInfo("❌  Could not kill the originally spawned process %d", processID)
				}
			}

			// Wait 10 seconds
			time.Sleep(10 * time.Second)

			// Start process
			processID = 0
			startProcessIfNotRunning()

			lastRequestTime = time.Now()
		} else if time.Since(lastRequestTime) > restartThreshold-1*time.Minute {
			logInfo("⚠️ WARNING: No requests were received for more than %s. Restarting process in one minute if still no request happened...", (restartThreshold - 1*time.Minute).String())
		}
	}
}
