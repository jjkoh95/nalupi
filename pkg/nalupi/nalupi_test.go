package nalupi_test

import (
	"math/big"
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

func TestCalculatePIOneIteration(t *testing.T) {
	Lk := nalupi.L0()
	Xk := nalupi.X0()
	Mk := nalupi.M0()
	termNumerator := big.NewInt(0).Mul(Mk, Lk)
	termDenominator := Xk
	term := big.NewInt(0)
	tempTerm := big.NewInt(0).Quo(termNumerator, termDenominator)
	term = term.Add(term, tempTerm)

	// 3.1415
	// precision here is required to get the values
	// since we are taking big.Int instead of float
	var precision int64 = 100
	C := nalupi.C(precision)
	inverseTerm := nalupi.OneOver(precision, term)
	got := big.NewInt(0).Mul(C, inverseTerm).String()[:5]
	want := "31415"
	require.Equal(t, want, got, "Expect to get first 5 digits of PI correctly with one iteration")
}
