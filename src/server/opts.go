package server

import (
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Opt func(s *Server)

func WithPort(port int) Opt {
	return func(s *Server) {
		s.server.Addr = ":" + strconv.Itoa(port)
	}
}

func WithDBPool(dbPool *pgxpool.Pool) Opt {
	return func(s *Server) {
		s.dbPool = dbPool
	}
}
