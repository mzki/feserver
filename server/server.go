package server

import (
	"context"
	"log"
	"net/http"

	"github.com/mzki/feserver/src"
)

// TODO: cache for Response from source.

// represents server which can get F.E. question from the external server,
// and can return json response containing F.E. question.
type Server struct {
	server *http.Server
	getter *src.Getter
	conf   Config
}

// it returns new constructed server with config.
// nil config is ok and use DefaultConfig insteadly.
func New(conf *Config) *Server {
	if conf == nil {
		conf = &DefaultConfig
	}
	return &Server{
		server: &http.Server{},
		getter: src.NewGetter(src.LeastIntervalTime),
		conf:   *conf,
	}
}

func (s *Server) timeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.conf.waitTime)
}

func (server *Server) getRandomQuestionJSON(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) getRandom(ctx context.Context, r *http.Request) *JSONResponse {
	qr, err := parseGetRandomQuery(r.URL.Query())
	if err != nil {
		return &JSONResponse{Error: err.Error()}
	}
	res, err := s.getter.GetRandom(ctx, &qr)
	if err != nil {
		return &JSONResponse{Error: err.Error()}
	}
	return &JSONResponse{Response: res}
}

func (server *Server) getQuestionJSON(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) getQuestion(ctx context.Context, r *http.Request) *JSONResponse {
	q, err := parseGetQuestionQuery(r.URL.Query())
	if err != nil {
		return &JSONResponse{Error: err.Error()}
	}
	res, err := s.getter.Get(ctx, q)
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

const (
	// represents URL for getting question randomly selected.
	URLGetRandom = "/r-question.json"
	// represents URL for getting question with specified query.
	URLGetQuestion = "/question.json"
)

// it starts server process.
// it blocks until process occurs any error and
// return the error.
func (s *Server) ListenAndServe() error {
	if err := s.conf.validates(); err != nil {
		return err
	}

	handler := http.NewServeMux()
	handler.HandleFunc(URLGetRandom, s.getRandomQuestionJSON)
	handler.HandleFunc(URLGetQuestion, s.getQuestionJSON)

	server := s.server
	server.Handler = handler

	serverURL := s.conf.URL + ":" + s.conf.Port
	server.Addr = serverURL
	log.Println("launch on " + serverURL + URLGetRandom)
	log.Println("launch on " + serverURL + URLGetQuestion)

	return s.server.ListenAndServe()
}

var defaultServer = New(nil)

// It starts server process using default server with
// user config.
// A nil config is OK and use DefaultConfig insteadly.
// It blocks until the process occurs any error and
// return the error.
func ListenAndServe(conf *Config) error {
	if conf == nil {
		conf = &DefaultConfig
	}
	defaultServer.conf = *conf
	return defaultServer.ListenAndServe()
}
