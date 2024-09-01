package utils

import (
	"github.com/dustin/go-humanize"
)

func DisplayPrice(val float64) string {
	return humanize.FormatFloat("#,###.##", val)
}
