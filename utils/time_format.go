package utils

import (
	"fmt"
	"math/bits"
)

// Port of https://github.com/DesterLib/Backend/blob/main/app/utils/time_formatter.py#L1
func TimeFormat(seconds uint64) (timeStr string) {
	hours, remainder := bits.Div64(0, seconds, 3600)
	minutes, seconds := bits.Div64(0, remainder, 60)
	days, hours := bits.Div64(0, hours, 24)
	timeStr = ""
	if days > 0 {
		timeStr += fmt.Sprintf("%d days, ", days)
	}
	timeStr += fmt.Sprintf("%d hours, %d minutes, %d seconds", hours, minutes, seconds)
	return
}
