// Package main implements the daemon of the search service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/OpenSlides/openslides-search-service/pkg/config"
	"github.com/OpenSlides/openslides-search-service/pkg/meta"
	"github.com/OpenSlides/openslides-search-service/pkg/search"
	"github.com/OpenSlides/openslides-search-service/pkg/web"
	"golang.org/x/sys/unix"
)

func check(err error) {
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func signalContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, unix.SIGTERM)
		<-sig
		cancel()
		<-sig
		os.Exit(2)
	}()
	return ctx, cancel
}

func run(cfg *config.Config) error {
	ctx, cancel := signalContext()
	defer cancel()

	models, err := meta.Fetch[meta.Collections](cfg.Models.Models)
	if err != nil {
		return fmt.Errorf("loading models failed: %w", err)
	}

	// For text indexing we can only use string fields.
	searchModels := models.Clone()
	searchModels.Retain(meta.RetainStrings(false))

	// If there are search filters configured cut search models further down.
	if cfg.Models.Search != "" {
		searchFilter, err := meta.Fetch[meta.Filters](cfg.Models.Search)
		if err != nil {
			return fmt.Errorf("loading search filters failed. %w", err)
		}
		searchModels.Retain(searchFilter.Retain(false))
	}

	db := search.NewDatabase(cfg)
	ti, err := search.NewTextIndex(cfg, db, searchModels)
	if err != nil {
		return fmt.Errorf("creating text index failed: %w", err)
	}
	defer ti.Close()

	runtime.GC()

	qs, err := search.NewQueryServer(cfg, ti)
	if err != nil {
		return err
	}
	go qs.Run(ctx)

	return web.Run(ctx, cfg, qs)
}

func main() {
	flag.Parse()
	cfg, err := config.GetConfig()
	check(err)
	check(run(cfg))
}
