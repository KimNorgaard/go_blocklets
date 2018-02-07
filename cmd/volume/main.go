package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"os/exec"
)

func getVolume() (int, error) {
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

			if vol, err = strconv.Atoi(matches[1]); err != nil {
				return 0, err
			}

			break
		}
	}

	return vol, nil
}

func main() {
	var output string
	var fullText string = "unknown"
	var shortText string = "unknown"
	var colorText string

	vol, err := getVolume()

	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks volume] Failed to get volume: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(0)
	}

	if vol == -1 {
		output = "----"
	} else {
		output = fmt.Sprintf("%02d%%", vol)
	}

	fullText = output
	shortText = output

	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, colorText)
	os.Exit(0)
}
