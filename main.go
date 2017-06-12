package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/singlehopllc/lmsensorbeat/beater"
)

func main() {
	err := beat.Run("lmsensorbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
