package rest

import (
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jjkoh95/nalupi/pkg/nalupi"
	"github.com/jjkoh95/nalupi/pkg/recipi"
)

// New generates a new http.Server instance
func New() *http.Server {
	r := mux.NewRouter()

	// health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok!"))
	})

	r.HandleFunc("/pi/current", func(w http.ResponseWriter, r *http.Request) {
		precision, pi, err := recipi.GetLastComputedPI()
		if err != nil {
			w.Write([]byte("Unable to get PI"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		integral, fractional := big.NewInt(0).QuoRem(pi, nalupi.TenPower(precision), big.NewInt(0))
		w.Write([]byte(fmt.Sprintf("%s.%s", integral.String(), fractional.String())))
	})

	r.HandleFunc("/pi/trigger", func(w http.ResponseWriter, r *http.Request) {
		k, Lk, Xk, Mk, Kk, err := recipi.GetSnapshot()
		if err != nil {
			w.Write([]byte("Unable to get snapshot"))
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		precision, _, err := recipi.GetLastComputedPI()
		if err != nil {
			w.Write([]byte("Unable to get precision"))
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		numerator, denominator, err := recipi.LoadFractionMeta()
		if err != nil {
			w.Write([]byte("Unable to read fraction meta"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// compute next iteration
		for {
			Lk = nalupi.Lplusone(Lk)
			Xk = nalupi.Xplusone(Xk)
			Mk = nalupi.Mplusone(Mk, Kk, k)
			Kk = nalupi.Kplusone(Kk)
			k.Add(k, big.NewInt(1))

			termNumerator := big.NewInt(0).Mul(Lk, Mk)
			// Mk*Lk*(262537412640768000)
			// this step is required for making common denominator
			numerator.Mul(
				numerator,
				big.NewInt(0).Abs(nalupi.Xmultiplier()),
			)
			// Mk*Lk*(262537412640768000)
			// we need to take into account what is the sign of Xk
			// to determine addition or subtraction action to do here
			if Xk.Sign() == -1 {
				numerator.Sub(
					numerator,
					termNumerator,
				)
			} else {
				numerator.Add(
					numerator,
					termNumerator,
				)
			}
			denominator = Xk

			// temporary term here to determine if we should proceed to the next iteration
			tempTerm := termNumerator
			tempTerm.Mul(tempTerm, nalupi.TenPower(precision+1))
			tempTerm.Quo(tempTerm, Xk)
			if tempTerm.Cmp(big.NewInt(0)) == 0 {
				break
			}
		}
		recipi.SaveSnapshot(k.String(), Lk.String(), Xk.String(), Mk.String(), Kk.String())
		recipi.SaveFractionMeta(numerator.String(), denominator.String())

		res := nalupi.C(precision + 1)
		res.Mul(res, big.NewInt(0).Abs(denominator))
		res.Quo(res, numerator)

		// also note that 3 is precision value = 0
		recipi.SaveComputedPI(strconv.FormatInt(precision+1, 10), res.String())
		w.Write([]byte("Ok"))
	})

	r.HandleFunc("/sun/circumference", func(w http.ResponseWriter, r *http.Request) {
		// get last computed PI and return 2*PI*R

		precision, pi, err := recipi.GetLastComputedPI()
		if err != nil {
			w.Write([]byte("Unable to get PI"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		sunRadius := nalupi.GetSunRadiusInKM()
		// circumference = 2*pi*radius
		res := big.NewInt(2)
		res.Mul(res, pi)
		res.Mul(res, sunRadius)

		// we can either do divisor, remainder
		// or index slicing
		// since big has a nice feature/function QuoRem, we will stick with this
		// also note precision-1 because on top we regard "3" as precision = 1
		// because we have 10^0 for precision = 1 when estimation of PI = 3
		integral, fractional := big.NewInt(0).QuoRem(res, nalupi.TenPower(precision), big.NewInt(0))

		w.Write([]byte(fmt.Sprintf("%s.%s", integral.String(), fractional.String())))
	})

	// setting port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // default port
	}

	// server instance
	return &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}
}
