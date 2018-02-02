package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"os/exec"

	gb "github.com/KimNorgaard/go_blocklets"
)

func volume() (int, error) {

	var volRegex = regexp.MustCompile(`\[(\d{1,3})\%\]`)

	vol := 0

	cmd := exec.Command("amixer", "sget", "Master")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(out))

	for scanner.Scan() {

		text := scanner.Text()

		if strings.Contains(text, "Front Left: Playback") {

			if strings.Contains(text, "[off]") {
				return -1, nil
			}

			matches := volRegex.FindStringSubmatch(text)
			if len(matches) != 2 {
				return 0, fmt.Errorf("expected two matches but found %d", len(matches))
			}

			vol, err = strconv.Atoi(matches[1])
			if err != nil {
				return 0, err
			}

			break
		}
	}

	return vol, nil
}

func getBatteryStatus() (string, int, error) {
	var status string
	var pctSum, pctLeft, batsFound int
	var pctAvg float64
	var stRegex = regexp.MustCompile(`^Battery \d: (\w+), (\d+)\%`)
	var matches []string

	cmd := exec.Command("acpi")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", 0, err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(out))

	for scanner.Scan() {
		text := scanner.Text()

		matches = stRegex.FindStringSubmatch(text)
		if len(matches) != 3 {
			return status, 0, fmt.Errorf("expected three matches but found %d", len(matches))
		}

		batsFound += 1

		pctLeft, err = strconv.Atoi(matches[2])
		if err != nil {
			return status, 0, err
		}
		pctSum += pctLeft
	}

	pctAvg = (float64(pctSum) / (float64(batsFound) * 100)) * 100

	return matches[1], int(pctAvg), nil
}

func getACStatus() (int, error) {
	statusFile := "/sys/class/power_supply/AC/online"

	if _, err := os.Stat(statusFile); os.IsNotExist(err) {
		return -1, nil
	}

	statusFileContent, err := ioutil.ReadFile(statusFile)
	if err != nil {
		return -1, err
	}

	lines := strings.Split(string(statusFileContent), "\n")
	if len(lines) == 0 {
		return -1, nil
	}

	status, err := strconv.Atoi(lines[0])
	if err != nil {
		return -1, err
	}

	return status, nil
}

func main() {
	var output string
	var fullText string = "unknown"
	var shortText string = "unknown"
	var colorText string
	//  - plug
	//  - empty bat
	//  - almost full bat
	//  - almost empty bat
	//  - full bat
	//  - charging / ac bat
	var icon string

	batStatus, pctLeft, err := getBatteryStatus()

	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks-battery] Failed to get battery status: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(0)
	}

	acStatus, err := getACStatus()

	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks-battery] Failed to get AC status: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(0)
	}

	// -1 means acpi was not found
	if pctLeft == -1 {
		output = "------"
		colorText = "#FF0000"
	} else {
		if acStatus == 1 {
			icon = ""
		} else if batStatus == "Full" {
			icon = ""
		} else if batStatus == "Discharging" || batStatus == "Charging" {
			if pctLeft > 50 {
				icon = ""
			} else {
				icon = ""
			}
		} else {
			icon = ""
		}

		output = fmt.Sprintf("%s %3d%%", icon, pctLeft)
		colorText = gb.GreenToRed(pctLeft)
	}

	fullText = output
	shortText = output

	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, colorText)
	os.Exit(0)
}
