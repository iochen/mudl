package size_test

import (
	"fmt"
	"testing"

	"github.com/iochen/mudl/utils/size"
)

func TestSize_String(t *testing.T) {
	i := size.Size(123456789)
	fmt.Println(i.String())
	size.Measure = 1000
	fmt.Println(i.String())

	size.Precision = 4
	size.Measure = 1 << 10
	fmt.Println(i.String())
	size.Measure = 1000
	fmt.Println(i.String())
}
