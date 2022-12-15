// SPDX-FileCopyrightText: 2022 Since 2011 Authors of OpenSlides, see https://github.com/OpenSlides/OpenSlides/blob/master/AUTHORS
//
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// storeString returns a function to store a string.
func storeString(dst *string) func(string) error {
	return func(v string) error {
		*dst = v
		return nil
	}
}

// storeInt returns a function to store an int.
func storeInt(dst *int) func(string) error {
	return func(v string) error {
		x, err := strconv.Atoi(v)
		if err == nil {
			*dst = x
		}
		return err
	}
}

// storeDuration returns a function to store a duration.
func storeDuration(dst *time.Duration) func(string) error {
	return func(v string) error {
		// If it can be parsed as an integer take that as seconds.
		secs, err := strconv.Atoi(v)
		if err == nil {
			*dst = time.Second * time.Duration(secs)
			return nil
		}
		x, err := time.ParseDuration(v)
		if err == nil {
			*dst = x
		}
		return err
	}
}

// storeEnv maps the name of an env var to a store function.
type storeEnv struct {
	name  string
	store func(string) error
}

// store iterates over the given env/stores calls the store function
// of every env var that is found.
func storeFromEnv(se []storeEnv) error {
	for _, s := range se {
		if v, ok := os.LookupEnv(s.name); ok {
			if err := s.store(v); err != nil {
				return fmt.Errorf("parsing env var %q failed: %w", s.name, err)
			}
		}
	}
	return nil
}
