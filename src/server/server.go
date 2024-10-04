package server

import (
	"net/http"
	"songs/src/impls"
	"songs/src/lib/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	dbPool      *pgxpool.Pool
	mux         *http.ServeMux
	server      *http.Server
	logger      logger.Logger
	middlewares []func(http.Handler) http.Handler
}

func New(opts ...Opt) *Server {
	mux := http.NewServeMux()
	s := &Server{
		mux: mux,
		server: &http.Server{
			Handler: mux,
		},
		logger: &logger.StdoutLogger{},
	}

	for _, opt := range opts {
		opt(s)
	}

	s.registerMiddlewares()
	s.registerHandlers()

	return s
}

func (s *Server) registerMiddlewares() {
	s.registerLogger()
}

func (s *Server) registerHandlers() {
	s.registerSongsHandlers(&impls.SongDetailsFinder{})
}

func (s *Server) handle(pattern string, fn func(http.ResponseWriter, *http.Request)) {
	var handler http.Handler = http.HandlerFunc(fn)
	for _, m := range s.middlewares {
		handler = m(handler)
	}

	s.mux.Handle(pattern, handler)
}

func (s *Server) Run() error {
	s.logger.Info("Server is listening on port " + s.server.Addr)
	return s.server.ListenAndServe()
}
