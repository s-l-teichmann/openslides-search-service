// Package main implements the daemon of the search service.
package main

import (
	"log"

	"github.com/OpenSlides/openslides-search-service/pkg/config"
)

func check(err error) {
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func main() {
	cfg, err := config.GetConfig()
	check(err)
	_ = cfg
}
