// Package web exposes a web api.
package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/OpenSlides/openslides-autoupdate-service/pkg/auth"
	"github.com/OpenSlides/openslides-search-service/pkg/config"
	"github.com/OpenSlides/openslides-search-service/pkg/search"
)

type controller struct {
	cfg  *config.Config
	auth *auth.Auth
	qs   *search.QueryServer
}

/*
func userIDFromRequest(r *http.Request) (int, error) {
	user := r.FormValue("u")
	if user == "" {
		return 0, errors.New("'u' parameter missing")
	}
	userID, err := strconv.Atoi(user)
	if err != nil {
		return 0, errors.New("'u' is not an user id")
	}
	return userID, nil
}
*/

func (c *controller) search(w http.ResponseWriter, r *http.Request) {

	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "'q' parameter missing", http.StatusBadRequest)
		return
	}

	answers, err := c.qs.Query(query)
	if err != nil {
		log.Printf("error: %v\n", err)
		http.Error(w, "Something went wrong. Check the server logs.",
			http.StatusInternalServerError)
		return
	}

	if c.cfg.Restricter.URL != "" {

		userID := c.auth.FromContext(r.Context())
		/*
			userID, err := userIDFromRequest(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		*/

		requestBody := struct {
			UserID int      `json:"user_id"`
			FQIDs  []string `json:"fqids"`
		}{
			UserID: userID,
			FQIDs:  answers,
		}

		body, err := json.Marshal(&requestBody)
		if err != nil {
			log.Printf("error: %v\n", err)
			http.Error(w, "Something went wrong. Check the server logs.",
				http.StatusInternalServerError)
			return
		}
		resp, err := http.Post(
			c.cfg.Restricter.URL,
			"application/json",
			bytes.NewReader(body))
		if err != nil {
			log.Printf("error: restricter call failed: %v\n", err)
			http.Error(w, "Something went wrong. Check the server logs.",
				http.StatusInternalServerError)
			return
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("error: restricter call failed: %q (%d)\n",
				resp.Status, resp.StatusCode)
			http.Error(w, "Something went wrong. Check the server logs.",
				http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("error: copy restricter output failed: %v\n", err)
		}
		return
	}

	// No restricter configured.

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(answers); err != nil {
		log.Printf("error: %v\n", err)
	}
}

// Run starts the web server and routes the incoming requests to the controller.
func Run(
	ctx context.Context,
	cfg *config.Config,
	auth *auth.Auth,
	qs *search.QueryServer,
) error {

	c := controller{
		cfg:  cfg,
		auth: auth,
		qs:   qs,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/search", c.search)

	addr := fmt.Sprintf("%s:%d", cfg.Web.Host, cfg.Web.Port)
	log.Printf("listen web on %s\n", addr)

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
