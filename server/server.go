package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mzki/feserver/src"
)

// TODO: cache for Response from source.

const (
	DefaultURL  = ""
	DefaultPort = "8080"

	WaitTime = 2 * time.Second
)

func getRandomQuestionJSON(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), WaitTime)
	defer cancel()

	resCh := make(chan src.Response)
	errCh := make(chan error)
	go func() {
		defer close(resCh)
		defer close(errCh)
		res, err := src.GetRandom(ctx)
		if err != nil {
			errCh <- err
		} else {
			resCh <- res
		}
	}()

	select {
	case res := <-resCh:
		if err := writeJSON(w, res); err != nil {
			log.Println("Error: " + err.Error())
			http.Error(w, "JSON API Error. Check server log", http.StatusInternalServerError)
		}
	case err := <-errCh:
		log.Println("Error: " + err.Error())
		http.Error(w, "API Error. Check server log", http.StatusInternalServerError)
	case <-ctx.Done():
		if err := ctx.Err(); err == context.DeadlineExceeded {
			http.Error(w, "Timeout", http.StatusRequestTimeout)
		}
	}
}

func writeJSON(w io.Writer, res src.Response) error {
	return json.NewEncoder(w).Encode(&res)
}

func ListenAndServe() error {
	http.HandleFunc("/r-question.json", getRandomQuestionJSON)

	serverURL := DefaultURL + ":" + DefaultPort
	log.Println("launch on " + serverURL + "/r-question.json")
	return http.ListenAndServe(serverURL, nil)
}
