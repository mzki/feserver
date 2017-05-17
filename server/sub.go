package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mzki/feserver/src"
)

// subServer listens on the sub address of the top-level server,
// and serves the questions from src.Source.
type subServer struct {
	getter   *src.Getter
	source   Source
	waitTime time.Duration
}

func newSubServer(s Source) *subServer {
	return &subServer{
		getter:   src.NewGetter(s.Source, src.LeastIntervalTime),
		source:   s,
		waitTime: time.Duration(s.WaitSecond) * time.Second,
	}
}

func (s *subServer) timeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.waitTime)
}

func (server *subServer) getRandomQuestionJSON(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := server.timeoutContext()
	defer cancel()

	resCh := make(chan *JSONResponse, 1)
	go func() {
		defer close(resCh)
		resCh <- server.getRandom(ctx, r)
	}()

	select {
	case jres := <-resCh:
		if err := jres.write(w); err != nil {
			serverError(w, err, "Writing JSON Error. Check server log.", http.StatusInternalServerError)
		}
	case <-ctx.Done():
		contextError(w, ctx)
	}
}

func (sub *subServer) getRandom(ctx context.Context, r *http.Request) *JSONResponse {
	qr, err := parseGetRandomQuery(r.URL.Query(), sub.source)
	if err != nil {
		return &JSONResponse{Error: err.Error()}
	}
	res, err := sub.getter.GetRandom(ctx, qr)
	if err != nil {
		return &JSONResponse{Error: err.Error()}
	}
	return &JSONResponse{Response: res}
}

func (server *subServer) getQuestionJSON(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := server.timeoutContext()
	defer cancel()

	resCh := make(chan *JSONResponse, 1)
	go func() {
		defer close(resCh)
		resCh <- server.getQuestion(ctx, r)
	}()

	select {
	case jres := <-resCh:
		if err := jres.write(w); err != nil {
			serverError(w, err, "Writing JSON Error. Check server log.", http.StatusInternalServerError)
		}
	case <-ctx.Done():
		contextError(w, ctx)
	}
}

func (sub *subServer) getQuestion(ctx context.Context, r *http.Request) *JSONResponse {
	q, err := parseGetQuestionQuery(r.URL.Query(), sub.source)
	if err != nil {
		return &JSONResponse{Error: err.Error()}
	}
	res, err := sub.getter.Get(ctx, q)
	if err != nil {
		return &JSONResponse{Error: err.Error()}
	}
	return &JSONResponse{Response: res}
}

func serverError(w http.ResponseWriter, err error, mes string, status int) {
	log.Println("Error: " + err.Error())
	http.Error(w, mes, status)
}

func contextError(w http.ResponseWriter, ctx context.Context) {
	err := ctx.Err()
	if err != context.DeadlineExceeded {
		serverError(w, err, "Unknown error. Check server log", http.StatusInternalServerError)
		return
	}

	jres := &JSONResponse{Error: "Request timeout. Please try again later."}
	if err := jres.write(w); err != nil {
		serverError(w, err, "Writing JSON Error. Check server log.", http.StatusInternalServerError)
	}
}
