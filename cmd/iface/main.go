package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func iface(ifName string) (string, bool, error) {
	var ifDir = fmt.Sprintf("/sys/class/net/%s", ifName)

	if _, err := os.Stat(ifDir); os.IsNotExist(err) {
		return "", false, nil
	}

	operState, err := ioutil.ReadFile(ifDir + "/operstate")
	if err != nil {
		return "", false, err
	}

	return strings.TrimSpace(string(operState)), true, nil
}

func getDefaultInterface() (string, error) {
	cmd := exec.Command("ip", "route")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(out))

	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "default") {
			return strings.Fields(text)[4], nil
		}
	}

	return "", nil
}

func main() {

	// Set display texts to defaults.
	var output string
	var fullText string = "unknown"
	var shortText string = "unknown"
	var colorText string
	var err error

	var ifName = os.Getenv("BLOCK_INSTANCE")
	if len(ifName) == 0 {
		if ifName, err = getDefaultInterface(); err != nil {
			fmt.Fprintf(os.Stderr, "[i3blocks iface] %s", err.Error())
			fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
			os.Exit(0)
		}
	}

	if ifName == "" {
		fmt.Print("----")
		fmt.Print("----")
		fmt.Print("#AAAAAA")
		os.Exit(0)
	}

	// Retrieve current iface info
	state, ifExists, err := iface(ifName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks iface] Failed to get info: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(0)
	}

	if !ifExists {
		os.Exit(0)
	}

	if state == "up" {
		colorText = "#00FF00"
	} else {
		colorText = "#FF0000"
	}

	output = ifName
	fullText = output
	shortText = output

	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, colorText)
	os.Exit(0)
}
