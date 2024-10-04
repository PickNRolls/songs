package config

import (
	"os"
	"strconv"
)

var env map[string]string

var Port int = 8000

var DB = dbConfig{
	Username: "postgres",
	Host:     "localhost",
	Port:     5432,
	Name:     "library",
}

func init() {
	var err error
	Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	if value, has := os.LookupEnv("DB_USERNAME"); has {
		DB.Username = value
	}

	if path, has := os.LookupEnv("POSTGRES_PASSWORD_FILE"); has {
		b, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		DB.Password = string(b[:len(b)-1])
	}

	if value, has := os.LookupEnv("DB_HOST"); has {
		DB.Host = value
	}

	if value, has := os.LookupEnv("DB_PORT"); has {
		var err error
		DB.Port, err = strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
	}

	if value, has := os.LookupEnv("POSTGRES_DB"); has {
		DB.Name = value
	}
}
