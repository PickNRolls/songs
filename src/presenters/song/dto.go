package presentersong

import "time"

type SongDTO struct {
	Id         int       `json:"id"`
	Group      string    `json:"group"`
	Name       string    `json:"name"`
	Verses     []string  `json:"verses"`
	Link       string    `json:"link"`
	ReleasedAt time.Time `json:"releasedAt"`
	AddedAt    time.Time `json:"addedAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
