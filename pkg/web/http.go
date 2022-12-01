// Package web exposes a web api.
package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/OpenSlides/openslides-search-service/pkg/search"
)

type controller struct {
	qs *search.QueryServer
}

func (c *controller) search(w http.ResponseWriter, r *http.Request) {

	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "'q' parameter missing", http.StatusBadRequest)
		return
	}

	answers, err := c.qs.Query(query)
	if err != nil {
		log.Printf("error: %v\n", err)
		http.Error(w, "Some went wrong. Check the server logs.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(answers); err != nil {
		log.Printf("error: %v\n", err)
	}
}

// Run starts the web server and routes the incoming requests to the controller.
func Run(ctx context.Context, addr string, qs *search.QueryServer) error {

	c := controller{qs: qs}

	mux := http.NewServeMux()

	mux.HandleFunc("/search", c.search)

	s := &http.Server{
		Addr:        addr,
		Handler:     mux,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}

	done := make(chan error)
	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			done <- fmt.Errorf("server error: %v", err)
			return
		}
		log.Println("web server done")
		done <- nil
	}()
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %v", err)
	}
	return <-done
}
