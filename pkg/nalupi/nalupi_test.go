package nalupi_test

import (
	"testing"

	"github.com/jjkoh95/nalupi/pkg/nalupi"
	"github.com/stretchr/testify/require"
)

func TestTenPower(t *testing.T) {
	var p int64
	var got string
	var want string

	p = 0
	got = nalupi.TenPower(p).String()
	want = "1"
	require.Equal(t, got, want, "Expect to get 1")

	p = 10
	got = nalupi.TenPower(p).String()
	want = "10000000000"
	require.Equal(t, got, want, "Expect to get 10000000000")
}

func TestCalculatePITo5Digits(t *testing.T) {
	got := nalupi.CalculatePIWithPrecision(5).String()
	got = got[:5]
	want := "31415"
	require.Equal(t, want, got, "Expect to get first 5 digits of PI correctly with one iteration")
}

func TestCalculatePITo100Digits(t *testing.T) {
	got := nalupi.CalculatePIWithPrecision(100).String()
	got = got[:100]
	want := "3141592653589793238462643383279502884197169399375105820974944592307816406286208998628034825342117067"
	require.Equal(t, want, got, "Expect to get the first 100 digits of PI correctly")
}
