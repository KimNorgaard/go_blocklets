package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func getPid() (int, error) {
	pidFile := "/run/openconnect.pid"

	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return -1, nil
	}

	pidFileContent, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return -1, err
	}

	lines := strings.Split(string(pidFileContent), "\n")
	if len(lines) == 0 {
		return -1, nil
	}

	pid, err := strconv.Atoi(lines[0])
	if err != nil {
		return -1, err
	}

	return pid, nil
}

func isRunning(pid int) (bool, error) {
	procDir := fmt.Sprintf("/proc/%d", pid)
	if _, err := os.Stat(procDir); os.IsNotExist(err) {
		return false, nil
	}

	cmdline, err := ioutil.ReadFile(procDir + "/cmdline")
	if err != nil {
		return false, err
	}

	if strings.HasPrefix(string(cmdline), "openconnect\000") {
		return true, nil
	}

	return false, nil
}

func main() {

	// Set display texts to defaults.
	var output string = ""
	var fullText string = "unknown"
	var shortText string = "unknown"
	var colorText string = "#FF0000"

	pid, err := getPid()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks openconnect] %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(0)
	}

	active, err := isRunning(pid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks openconnect] %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(0)
	}

	if active {
		colorText = "#00FF00"
		output = ""
	}

	fullText = output
	shortText = output

	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, colorText)
	os.Exit(0)
}
