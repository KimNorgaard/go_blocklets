package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	gb "github.com/KimNorgaard/go_blocklets"
)

const PowerSupplyPath string = "/sys/class/power_supply/"

func IsBattery(supply string) bool {
	supplyTypeFile := filepath.Join(PowerSupplyPath, supply, "type")

	if str, err := gb.StringFromFile(supplyTypeFile); err == nil && strings.Contains(str, "Battery") {
		return true
	}
	return false
}

func GetBatteries() (batteries []string) {
	paths, err := filepath.Glob(PowerSupplyPath + "*")
	if err != nil {
		return
	}

	for _, path := range paths {
		supply := filepath.Base(path)
		if IsBattery(supply) {
			batteries = append(batteries, supply)
		}
	}

	return
}

func GetBatteryStatus(battery string) string {
	status, _ := gb.StringFromFile(filepath.Join(PowerSupplyPath, battery, "status"))
	return status
}

func GetBatteryChargeNow(battery string) int64 {
	chargeNow, _ := gb.StringFromFile(filepath.Join(PowerSupplyPath, battery, "energy_now"))
	i, _ := strconv.ParseInt(chargeNow, 10, 64)
	return i
}

func GetBatteryChargeFull(battery string) int64 {
	chargeFull, _ := gb.StringFromFile(filepath.Join(PowerSupplyPath, battery, "energy_full"))
	i, _ := strconv.ParseInt(chargeFull, 10, 64)
	return i
}

func GetBatteryPercent(battery string) float64 {
	return float64(100) * float64(GetBatteryChargeNow(battery)) / float64(GetBatteryChargeFull(battery))
}

func GetBatteriesChargeNow(batteries []string) int64 {
	var nowSum int64 = 0
	for _, battery := range batteries {
		nowSum = nowSum + GetBatteryChargeNow(battery)
	}
	return nowSum
}

func GetBatteriesChargeFull(batteries []string) int64 {
	var fullSum int64 = 0
	for _, battery := range batteries {
		fullSum = fullSum + GetBatteryChargeFull(battery)
	}
	return fullSum
}

func GetBatteriesPercent(batteries []string) float64 {
	return (float64(100) * float64(GetBatteriesChargeNow(batteries)) / float64(GetBatteriesChargeFull(batteries)))
}

func GetBatteriesStatus(batteries []string) string {
	var maxStatus string
	var maxStatusValue int

	statusMap := map[string]int{
		"Unknown":      0,
		"Full":         1,
		"Not Charging": 2,
		"Charging":     3,
		"Discharging":  4,
	}

	for _, battery := range batteries {
		status := GetBatteryStatus(battery)

		statusValue := statusMap[status]
		if statusValue > maxStatusValue {
			maxStatusValue = statusValue
		}
	}

	for k, v := range statusMap {
		if v == maxStatusValue {
			maxStatus = k
		}
	}

	return maxStatus
}

func main() {
	var output string = "na"
	var fullText string = "unknown"
	var shortText string = "unknown"
	var colorText string

	batteries := GetBatteries()

	pctLeft := GetBatteriesPercent(batteries)
	batStatus := GetBatteriesStatus(batteries)

	switch batStatus {
	case "Charging":
		output = "c"
	case "Not Charging":
		output = "n"
	case "Full":
		output = "f"
	case "Discharging":
		output = "d"
	}

	output = fmt.Sprintf("%s %3.0f%%", output, pctLeft)
	if pctLeft < 20 {
		colorText = gb.GreenToRed(int(pctLeft))
	}

	fullText = output
	shortText = output

	fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n", fullText, shortText, colorText)
	os.Exit(0)
}
