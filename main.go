package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	port             int
	processCommand   string
	processTitle     string
	idleMinutes      int
	restartThreshold time.Duration
	lastRequestTime  time.Time
	processID        int
	debugLogs        bool
)

func init() {
	flag.StringVar(&processCommand, "process", "", "Name of the process to manage")
	flag.IntVar(&port, "port", 0, "Port number to listen on")
	flag.StringVar(&processTitle, "title", "", "Title of the target window")
	flag.IntVar(&idleMinutes, "idletime", 10, "Time in `minutes` which the game needs to be stuck or idle before restart is triggered")
	flag.BoolVar(&debugLogs, "debug", false, "Log all HTTP requests on console?")
	flag.Parse()

	if processCommand == "" || port == 0 {
		fmt.Println("Both -process and -port parameters are required")
		flag.PrintDefaults()
		os.Exit(1)
	}
	restartThreshold = time.Duration(idleMinutes) * time.Minute
	logDebug("idletime = " + restartThreshold.String())
	logDebug("Startup successful")
}

/*
Example call:
main.exe -debug=true -port=8888 -process="C:\LDPlayer\LDPlayer4.0\dnplayer.exe index=1" -title "TsumTsum Alice"
*/
func main() {
	startProcessIfNotRunning()

	http.HandleFunc("/", handler)
	go monitorAndRestartProcess()

	logInfo("Server started on port %s", strconv.Itoa(port))
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println("HTTP serve error:", err)
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		fmt.Println("HTTP handler error:", err)
		os.Exit(1)
	}
	logDebug("received GET")
	lastRequestTime = time.Now()
}

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

func startProcessIfNotRunning() {
	logInfo("Title = %s", processTitle)
	if processTitle != "" {
		filterQuery := fmt.Sprintf("WINDOWTITLE eq %s", processTitle)
		stdout, err := exec.Command("cmd", "/C", "tasklist", "/FI", filterQuery, "/FO", "CSV", "/NH").CombinedOutput()
		if err != nil {
			log.Fatal("Unable to start tasklist command", err)
		} else {
			csvReader := csv.NewReader(strings.NewReader(string(stdout[:])))
			records, err := csvReader.ReadAll()
			if err != nil {
				log.Fatal("Unable to parse file as CSV", err)
			}
			if len(records) > 1 {
				log.Fatal("Found more than one process with window title", processTitle)
			}
			if len(records) == 1 {
				line := records[0]
				if len(line) >= 2 {
					processIdStr := line[1]
					processIdInt, err := strconv.Atoi(processIdStr)
					if err != nil {
						log.Fatal("Unable to fetch processId from tasklist", err)
					}
					if processIdInt > 0 {
						logInfo("Found already running process with id %d. Not starting a new process", processIdInt)
						processID = processIdInt
					}
				}
			}
		}
	}
	if processID == 0 {
		command := strings.Fields(processCommand)
		logInfo("Start command: %s", processCommand)
		cmd := exec.Command(command[0], command[1:]...)
		err := cmd.Start()
		if err != nil {
			log.Fatal("Error starting process:", err)
		}
		processID = cmd.Process.Pid
		logInfo("Process started with PID %s", strconv.Itoa(processID))
	}
	lastRequestTime = time.Now()
}

func monitorAndRestartProcess() {
	for {
		time.Sleep(1 * time.Minute)
		logDebug("lastRequestTime = " + lastRequestTime.String())

		if time.Since(lastRequestTime) > restartThreshold {
			logInfo("No requests were received for more than %s. Restarting process...", restartThreshold.String())

			// Kill process
			killCmd := exec.Command("tskill", strconv.Itoa(processID))
			stdout, err := killCmd.CombinedOutput()
			processID = 0
			if err != nil {
				fmt.Println("Error killing process:", stdout, ", rc = ", err)
			}

			// Wait 10 seconds
			time.Sleep(10 * time.Second)

			// Start process
			startProcessIfNotRunning()

			lastRequestTime = time.Now()
		} else if time.Since(lastRequestTime) > restartThreshold-1*time.Minute {
			logInfo("WARNING: No requests were received for more than %s. Restarting process in one minute if still no request happened...", (restartThreshold - 1*time.Minute).String())
		}
	}
}
