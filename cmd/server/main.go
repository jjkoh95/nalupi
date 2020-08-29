package main

import (
	"github.com/jjkoh95/nalupi/pkg/rest"
)

func main() {
	srv := rest.New()
	srv.ListenAndServe()
}
