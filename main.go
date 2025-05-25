package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Entry struct {
	Project  string
	Start    time.Time
	End      time.Time
	Duration string
}

var dataDir = filepath.Join(os.Getenv("HOME"), ".chronos")

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: timetrack <start|stop|status|log>")
		return
	}

	switch os.Args[1] {
	case "start":
		startTracking(os.Args[2:])
	case "stop":
		stopTracking()
	}
}

func startTracking(args []string) {
	if len(args) < 1 {
		fmt.Println("Please specify a project name.")
		return
	}

	os.MkdirAll(dataDir, 0755)

	project := strings.Join(args, " ")
	now := time.Now()

	entry := Entry{
		Project: project,
		Start:   now,
	}
}

func stopTracking() {

}
