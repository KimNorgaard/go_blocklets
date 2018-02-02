package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	gobl "github.com/KimNorgaard/go_blocklets"
)

func wifi(wlIf string) (int, error) {
	// wifi quality
	var quality int = 0

	var ifDir = fmt.Sprintf("/sys/class/net/%s", wlIf)
	if _, err := os.Stat(ifDir + "/wireless"); os.IsNotExist(err) {
		return 0, err
	}

	operState, err := ioutil.ReadFile(ifDir + "/operstate")
	if err != nil {
		return 0, err
	}

	// Interface is down, exit normally
	if string(operState) == "down\n" {
		return -1, nil
	}

	wlInfoF, err := os.Open("/proc/net/wireless")
	if err != nil {
		return 0, err
	}
	defer wlInfoF.Close()

	scanner := bufio.NewScanner(wlInfoF)

	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, wlIf) {
			fields := strings.Fields(text)
			qualifyF, err := strconv.ParseFloat(fields[2], 64)
			quality = int(qualifyF)
			if err != nil {
				return 0, err
			}
			return quality, nil
			break
		}
	}

	return quality, nil
}

func main() {

	// Set display texts to defaults.
	var output string
	var fullText string = "unknown"
	var shortText string = "unknown"
	var colorText string = "#19C500"

	var wlIf = os.Getenv("BLOCK_INSTANCE")
	if len(wlIf) == 0 {
		wlIf = "wlan0"
	}

	// Retrieve current WIFI info
	quality, err := wifi(wlIf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[i3blocks wifi] Failed to get info: %s", err.Error())
		fmt.Fprintf(os.Stdout, "%s\n%s\n", fullText, shortText)
		os.Exit(0)
	}

	output = fmt.Sprintf("%3d%%", quality)

	// Interface down. Exit cleanly.
	if quality == -1 {
		output = wlIf
		colorText = "#FF0000"
	} else {
		colorText = gobl.getGreenToRed(quality)
	}

	fullText = output
	shortText = output

	// Write out gathered information to STDOUT.
	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, colorText)
	os.Exit(0)
}
