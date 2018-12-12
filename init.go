package main

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("orion")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{time:2006-01-02 15:04:05.000} %{level:.4s} %{shortfile}:%{shortfunc} %{message}`,
)

func init() {
	logFile, err := os.OpenFile("./orion.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(-2)
	}
	backendFile := logging.NewLogBackend(logFile, "", 0)
	backendFileFormatter := logging.NewBackendFormatter(backendFile, format)
	backendFileLeveled := logging.AddModuleLevel(backendFileFormatter)
	backendFileLeveled.SetLevel(logging.DEBUG, "")

	backendStderr := logging.NewLogBackend(os.Stderr, "", 0)
	backendStderrFormatter := logging.NewBackendFormatter(backendStderr, format)
	backendStderrLeveled := logging.AddModuleLevel(backendStderrFormatter)
	backendStderrLeveled.SetLevel(logging.CRITICAL, "")

	logging.SetBackend(backendFileLeveled, backendStderrLeveled)

}
