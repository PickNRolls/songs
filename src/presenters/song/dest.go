package presentersong

import "time"

type dest struct {
	Id         int       `db:"id"`
	Group      string    `db:"group_name"`
	Name       string    `db:"name"`
	Verses     []string  `db:"verses"`
	Link       string    `db:"link"`
	ReleasedAt time.Time `db:"released_at"`
	AddedAt    time.Time `db:"added_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
