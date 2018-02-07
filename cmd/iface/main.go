package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	gb "github.com/KimNorgaard/go_blocklets"
)

func getDefaultInterface() (ifName string, err error) {
	cmd := exec.Command("ip", "route")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(out))

	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "default") {
			ifName = strings.Fields(text)[4]
			break
		}
	}

	return ifName, nil
}

func main() {
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

	if len(ifName) == 0 {
		fmt.Fprintf(os.Stderr, "[i3blocks iface] no default interface found")
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(0)
	}

	status, err := gb.GetIfaceStatus(ifName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks iface] Failed to get info: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(0)
	}

	if status == gb.IfStatusNonExistant {
		os.Exit(0)
	} else if status == gb.IfStatusUp {
		output = ifName
	}

	fullText = output
	shortText = output

	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, colorText)
	os.Exit(0)
}
