package main

import (
	"songs/src/config"
	"songs/src/postgres"
	"songs/src/server"
)

func main() {
	s := server.New(
		server.WithPort(config.Port),
		server.WithDBPool(postgres.Connect()),
	)

	s.Run()
}
