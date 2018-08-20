package go_blocklets

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

func GreenToRed(pct int) string {
	var r, g float64

	pctFloat := float64(pct)

	if pctFloat < 50 {
		r = 255
		g = math.Floor((pctFloat * 2) * 255 / 100)
	} else {
		r = math.Floor(255 - (pctFloat*2-100)*255/100)
		g = 255
	}

	return fmt.Sprintf("#%02X%02X%02X", int(r), int(g), 0)
}

const (
	IfStatusNonExistant = 0
	IfStatusDown        = 1
	IfStatusUp          = 2
	IfStatusError       = -1
)

func GetIfaceStatus(ifName string) (int, error) {
	fileName := fmt.Sprintf("/sys/class/net/%s/operstate", ifName)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return IfStatusNonExistant, nil
	}

	status, err := ioutil.ReadFile(fileName)
	if err != nil {
		return IfStatusError, err
	}

	statusText := strings.TrimSpace(string(status))
	if statusText == "up" {
		return IfStatusUp, nil
	} else {
		return IfStatusDown, nil
	}
}

func StringFromFile(fileName string) (string, error) {
	fh, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer fh.Close()
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(bytes.NewBuffer(b).String(), " \n"), nil
}
