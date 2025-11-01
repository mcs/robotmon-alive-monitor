package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Pid int

func findByTitle(processTitle string) Pid {
	pId := 0
	if strings.TrimSpace(processTitle) != "" {
		filterQuery := fmt.Sprintf("WINDOWTITLE eq %s", processTitle)
		stdout, err := exec.Command("cmd", "/C", "tasklist", "/FI", filterQuery, "/FO", "CSV", "/NH").CombinedOutput()
		if err != nil {
			log.Fatal("❌  Unable to start tasklist command", err)
		} else {
			csvReader := csv.NewReader(strings.NewReader(string(stdout[:])))
			records, err := csvReader.ReadAll()
			if err != nil {
				log.Fatal("❌  Unable to parse file as CSV", err)
			}
			if len(records) > 1 {
				log.Fatal("❌  Found more than one process with window title", processTitle)
			}
			if len(records) == 1 {
				line := records[0]
				if len(line) >= 2 {
					processIdStr := line[1]
					var err error
					pId, err = strconv.Atoi(processIdStr)
					if err != nil {
						log.Fatal("❌  Unable to fetch pId from tasklist", err)
					}
					if pId <= 0 {
						logInfo("⚠️ Got pid %d which is a non-positive value", pId)
					}
				}
			}
		}
	}
	return Pid(pId)
}

func kill(pId Pid) error {
	pid := int(pId)
	// kill process using tskill or taskkill
	var killCmd *exec.Cmd
	if isCommandAvailable("tskill") {
		killCmd = exec.Command("tskill", strconv.Itoa(pid))
	} else {
		killCmd = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
	}

	stdout, err := killCmd.CombinedOutput()
	if err != nil {
		logInfo("❌  Error killing process %d: %s, rc = %s", pid, string(stdout), err.Error())
	}
	return err
}

func startCommand(processCommand string) (Pid, error) {
	logInfo("Received command parameter: %s", processCommand)
	pid := 0
	commandParams := strings.Fields(processCommand)
	cmdIndex := 0
	candidate := ""
	for i := 1; cmdIndex == 0 && i <= len(commandParams); i++ {
		candidate = strings.Join(commandParams[:i], " ")
		logDebug("Found candidate: " + candidate)
		_, err := os.Stat(candidate)
		if err == nil {
			// found a real file
			cmdIndex = i
		} else {
			logDebug("Could not find file: " + candidate)
		}
	}
	if cmdIndex > 0 {
		cmd := exec.Command(candidate, commandParams[cmdIndex:]...)
		err := cmd.Start()
		if err == nil {
			pid = cmd.Process.Pid
		}
		return Pid(pid), err

	}
	return Pid(pid), fmt.Errorf("could not resolve executable path in: %q", processCommand)

}
