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
		_, pi, err := recipi.GetLastComputedPI()
		if err != nil {
			w.Write([]byte("Unable to get PI"))
		}
		w.Write([]byte(pi.String()))
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

		fractionSum, err := recipi.LoadFractionMeta(precision + 1)
		fmt.Println(fractionSum)
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

			tempTerm := big.NewInt(0).Mul(Lk, Mk)
			recipi.SaveFractionMeta(tempTerm.String(), Xk.String())
			tempTerm.Mul(tempTerm, nalupi.TenPower(precision+1))

			fmt.Println(tempTerm, Xk)
			tempTerm.Quo(tempTerm, Xk)
			if tempTerm.Cmp(big.NewInt(0)) == 0 {
				break
			}
			// else add to fraction
			fractionSum.Add(fractionSum, tempTerm)
		}
		recipi.SaveSnapshot(k.String(), Lk.String(), Xk.String(), Kk.String(), Mk.String())
		fmt.Println(fractionSum)
		fractionSum.Quo(nalupi.C(precision+1), fractionSum)
		recipi.SaveComputedPI(strconv.FormatInt(precision+1, 10), fractionSum.String())
		w.Write([]byte("Ok"))
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
		WriteTimeout: 10 * time.Minute,
		ReadTimeout:  10 * time.Minute,
	}
}
