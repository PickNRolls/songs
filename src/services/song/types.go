package servicesong

import (
	"songs/src/domain/song"
	"time"
)

type Details interface {
	ReleaseDate() time.Time
	Text() string
	Link() string
}

type Id = int

type DetailsFinder interface {
	Find(group string, name string) (Details, error)
}

type Repository interface {
	FindById(id Id) (*song.Song, error)
	Save(song *song.Song) (*song.Song, error)
	DeleteById(id Id) error
}
