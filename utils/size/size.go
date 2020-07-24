package size

import (
	"fmt"
	"strconv"
)

// Count in byte(8bits)
type Size int64

var Measure = Size(1 << 10)

var Precision = 1

var HaveSpace = true

// 1TB = 1024 GB = 1024*1024
func (s *Size) String() string {
	if *s <0 {
		return "unknown"
	}
	p := strconv.Itoa(Precision)
	sp := ""
	f := "%ciB"
	if Measure == 1000 {
		f = "%cB"
	}
	if HaveSpace {
		sp = " "
	}
	if *s < Measure {
		return fmt.Sprintf("%d"+sp+f, *s)
	}
	div, exp := Measure, 0
	for n := *s / Measure; n >= Measure; n /= Measure {
		div *= Measure
		exp++
	}
	return fmt.Sprintf("%."+p+"f"+sp+f, float64(*s)/float64(div), "KMGTPE"[exp])
}
