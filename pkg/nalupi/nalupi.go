package nalupi

import (
	"math/big"
)

// TenPower returns 10^power (required to offset float)
func TenPower(power int64) *big.Int {
	return big.NewInt(0).Exp(big.NewInt(10), big.NewInt(power), nil)
}

// L0 is the zero-th term for L
func L0() *big.Int {
	return big.NewInt(13591409)
}

// Lplusone takes L[k] and returns L[k+1] = L[k] + 545140134,
// function returns a new *big.Int instance instead of changing
// the existing reference so that the function is idempotent
func Lplusone(val *big.Int) *big.Int {
	Ld := big.NewInt(545140134)
	return big.NewInt(0).Add(val, Ld)
}

// X0 is the zero-th term for X
func X0() *big.Int {
	return big.NewInt(1)
}

// Xplusone computes X[k+1]
func Xplusone(val *big.Int) *big.Int {
	// 262537412640768000 = length of 18
	Xd := big.NewInt(-262537412640768000) // int64 can hold this
	return big.NewInt(0).Mul(val, Xd)
}

// M0 is the zero-th term for M
func M0() *big.Int {
	return big.NewInt(1)
}

// Mplusone computes M[k+1]
func Mplusone(Mk, Kk, k *big.Int) *big.Int {
	numerator := big.NewInt(0)
	numerator.Exp(Kk, big.NewInt(3), nil)                           // Kk^3
	numerator.Sub(numerator, big.NewInt(0).Mul(big.NewInt(16), Kk)) // Kk^3 - 16Kk
	denominator := big.NewInt(0)
	denominator.Add(k, big.NewInt(1))                // k + 1
	denominator.Exp(denominator, big.NewInt(3), nil) // (k + 1)^3
	// we are guaranteed this will be integer
	res := big.NewInt(0)
	res.Mul(numerator, Mk)
	res.Quo(res, denominator)
	return res
}

// K0 is the zero-th term for K
func K0() *big.Int {
	return big.NewInt(6)
}

// Kplusone computes K[k+1] given K[k]
func Kplusone(val *big.Int) *big.Int {
	Kd := big.NewInt(12)
	return big.NewInt(0).Add(val, Kd)
}

// C returns C value multiplied 10^factor
func C(factor int64) *big.Int {
	res := TenPower(2 * factor) // square this to make precision
	t1 := big.NewInt(426880 * 426880)
	t2 := big.NewInt(10005)
	res.Mul(res, t1)
	res.Mul(res, t2)
	res = res.Sqrt(res)
	return res
}

// OneOver returns 10^factor/term
func OneOver(factor int64, term *big.Int) *big.Int {
	pow := TenPower(factor)
	return big.NewInt(0).Div(pow, term)
}

// CalculatePIWithPrecision computes PI to precision as defined
func CalculatePIWithPrecision(precision int64) *big.Int {
	runningSum := big.NewInt(0)
	Lk := L0()
	Xk := X0()
	Mk := M0()
	Kk := K0()
	k := big.NewInt(0)
	multiplier := precision
	C0 := C(multiplier)
	for {
		termVal := TenPower(multiplier)
		termVal.Mul(termVal, Lk)
		termVal.Mul(termVal, Mk)
		termVal.Quo(termVal, Xk)
		// termVal = Mk * Lk / Xk
		// if the term is 0, stop, since further computation is not doing anything
		if termVal.Cmp(big.NewInt(0)) == 0 {
			break
		}
		runningSum.Add(runningSum, termVal)
		Lk = Lplusone(Lk)
		Xk = Xplusone(Xk)
		Mk = Mplusone(Mk, Kk, k)
		Kk = Kplusone(Kk)
		k.Add(k, big.NewInt(1))
	}
	res := TenPower(precision)
	res = res.Mul(res, C0)
	res = res.Quo(res, runningSum)
	return res
}

// CalculatePI calculates PI using Chudnovsky algorithm
func CalculatePI() {
	// 1 / pi:
	// numerator: 12*(-1)^k * (6*k)! ) * (545140134*k + 13591409)
	// denominator: (3*k)! * (k!)^3 * (640320)^(3k+3/2)
	// Or better:
	// pi = C * SUM(M[k] * L[k] / X[k])^-1
	// Where:
	// C = 426880 * SQRT(10005)
	// M[k] = (6k)! / ((3k)! * (k!)^3)
	// L[k] = 545140134*k + 13591409
	// X[k] = (-262537412640768000)^k
	// Where L, X, M can be better computed:
	// L[k + 1] = L[k] + 545140134 where L[0] = 13591409
	// X[k + 1] = X[k]*(-262537412640768000) where X[0] = 1
	// M[k + 1] = M[k]*(K[k]^3 - 16*K[k])/((k + 1)^3) where M[0] = 1, K[k + 1] = K[k] + 12, where K[0] = 6

	// NOTE THAT M, L, X ARE ALL *BIG* INTEGERS, SO WE DO NOT NEED TO COMPUTE THOSE TO CERTAIN PRECISION
	// HOWEVER, WE ARE SQRT(10005) WHICH IS A IRRATIONAL NUMBER, WHICH WE WILL NEED TO COMPUTE TO A CERTAIN PRECISION
	// TO ACCOMMODATE OUR PI PRECISION
	// SAME TO THE COMPUTATION TO SUM(terms)^-1

	// NOTE: The n-th term is a super small value
	panic("Not implemented, please use CalculatePIWithPrecision function instead")
}
