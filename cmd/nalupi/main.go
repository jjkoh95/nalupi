package main

import (
	"flag"
	"fmt"

	"github.com/jjkoh95/nalupi/pkg/nalupi"
)

var (
	decimal = flag.Int64("decimal", 100, "Precision of PI returned")
)

func main() {
	flag.Parse()
	fmt.Println(nalupi.CalculatePIWithPrecision(*decimal))
}
