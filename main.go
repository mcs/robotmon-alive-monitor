package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	port             int
	processCommand   string
	processTitle     string
	idleMinutes      int
	restartThreshold time.Duration
	lastRequestTime  time.Time
	processID        Pid
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

	ipAddress := getOutboundIP().String()
	logInfo("Use this URL within Robotmon: *** http://" + ipAddress + ":" + strconv.Itoa(port) + " ***")

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

	logInfo("ℹ️ Server started on port %s", strconv.Itoa(port))
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println("❌  HTTP serve error:", err)
		os.Exit(1)
	}
}

func getOutboundIP() net.IP {
	// Destination here must be valid, else it may give a loopback address
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}

func isCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
