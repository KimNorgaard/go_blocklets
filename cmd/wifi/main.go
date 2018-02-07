package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	gb "github.com/KimNorgaard/go_blocklets"
)

func getWifiQuality(ifName string) (int, error) {
	var quality float64

	f, err := os.Open("/proc/net/wireless")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, ifName) {
			fields := strings.Fields(text)
			if quality, err = strconv.ParseFloat(fields[2], 64); err != nil {
				return 0, err
			}
			break
		}
	}

	return int(quality), nil
}

func main() {
	var output string
	var fullText string = "unknown"
	var shortText string = "unknown"
	var colorText string

	ifName := os.Getenv("BLOCK_INSTANCE")
	if len(ifName) == 0 {
		ifName = "wlan0"
	}

	status, err := gb.GetIfaceStatus(ifName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks wifi] Failed to get info: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(0)
	}

	if status == gb.IfStatusNonExistant {
		os.Exit(0)
	} else if status == gb.IfStatusUp {
		quality, err := getWifiQuality(ifName)

		if err != nil {
			fmt.Fprintf(os.Stderr, "[i3blocks wifi] Failed to get info: %s", err.Error())
			fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
			os.Exit(0)
		}

		if quality < 50 {
			colorText = gb.GreenToRed(quality)
			output = fmt.Sprintf("%s <span color='%s'>%02d%%</span>", ifName, colorText, quality)
		} else {
			output = fmt.Sprintf("%s %02d%%", ifName, quality)
		}

	}

	fullText = output
	shortText = output

	fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
	os.Exit(0)
}
