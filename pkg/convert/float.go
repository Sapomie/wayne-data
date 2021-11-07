package convert

import (
	"fmt"
	"strconv"
)

type FloatTo float64

func (f FloatTo) Decimal(decimal int) float64 {
	format := "%." + fmt.Sprint(decimal) + "f"
	str := fmt.Sprintf(format, f)
	f2, _ := strconv.ParseFloat(str, 64)
	return f2
}
