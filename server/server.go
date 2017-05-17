package server

import (
	"log"
	"net/http"
	"path"
)

// TODO: cache for Response from source.

// represents server which can get F.E. question from the external server,
// and can return json response containing F.E. question.
type Server struct {
	server     *http.Server
	subServers map[string]*subServer
	conf       Config
}

// it returns new constructed server with config.
// nil config is ok and use DefaultConfig insteadly.
func New(conf *Config) *Server {
	if conf == nil {
		conf = &DefaultConfig
	}

	ss := make(map[string]*subServer, len(conf.Sources))
	for _, s := range conf.Sources {
		ss[s.SubAddr] = newSubServer(s)
	}

	return &Server{
		server:     &http.Server{},
		subServers: ss,
		conf:       *conf,
	}
}

const (
	// represents API for getting question randomly selected.
	APIGetRandom = "/r-question.json"
	// represents API for getting question with specified query.
	APIGetQuestion = "/question.json"
)

// it starts server process.
// it blocks until process occurs any error and
// return the error.
func (s *Server) ListenAndServe() error {
	if err := s.conf.validates(); err != nil {
		return err
	}

	serverURL := s.conf.URL + ":" + s.conf.Port

	handler := http.NewServeMux()
	for addr, sub := range s.subServers {
		for _, api := range []struct {
			path    string
			handler func(http.ResponseWriter, *http.Request)
		}{
			{addr + APIGetRandom, sub.getRandomQuestionJSON},
			{addr + APIGetQuestion, sub.getQuestionJSON},
		} {
			handler.HandleFunc(api.path, api.handler)
			log.Println("listen on " + path.Join(serverURL, api.path))
		}
	}

	s.server.Addr = serverURL
	s.server.Handler = handler
	return s.server.ListenAndServe()
}

// It starts server process using default server with
// user config.
// A nil config is OK and use DefaultConfig insteadly.
// It blocks until the process occurs any error and
// return the error.
func ListenAndServe(conf *Config) error {
	if conf == nil {
		conf = &DefaultConfig
	}
	defaultServer := New(conf)
	return defaultServer.ListenAndServe()
}
