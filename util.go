package go_blocklets

import (
	"fmt"
	"math"
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
