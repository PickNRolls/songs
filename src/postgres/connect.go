package postgres

import (
	"context"
	"songs/src/config"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() *pgxpool.Pool {
	url := "postgres://" + config.DB.Username + ":" +
		config.DB.Password + "@" + config.DB.Host + ":" +
		strconv.Itoa(config.DB.Port) + "/" + config.DB.Name

	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		panic(err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return db
}
