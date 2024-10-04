package impls

import (
	servicesong "songs/src/services/song"
	"time"
)

type SongDetailsFinder struct{}

func (df *SongDetailsFinder) Find(group string, name string) (servicesong.Details, error) {
	return &details{}, nil
}

type details struct{}

func (d *details) ReleaseDate() time.Time {
	return time.Now()
}

func (d *details) Text() string {
	return "text\n\nmy new verse\n\nmy third verse"
}

func (d *details) Link() string {
	return "google.com"
}
