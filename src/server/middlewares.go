package server

import "net/http"

type middleware func(http.Handler) http.Handler

func (s *Server) registerMiddleware(m middleware) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) registerLogger() {
	s.registerMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.logger.Info(r.Method + " " + r.URL.Path)

			h.ServeHTTP(w, r)
		})
	})
}
