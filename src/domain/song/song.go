package song

import "time"

type Song struct {
	Id         int       `json:"id" db:"id"`
	Group      string    `json:"group" db:"group_name"`
	Name       string    `json:"name" db:"name"`
	Verses     []string  `json:"verses" db:"verses"`
	Link       string    `json:"link" db:"link"`
	ReleasedAt time.Time `json:"releasedAt" db:"released_at"`
	AddedAt    time.Time `json:"addedAt" db:"added_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`

	Events []Event `json:"-"`
}

func New(id int, group string, name string, verses []string, link string, releaseDate time.Time, addedAt time.Time, updatedAt time.Time, raiseCreatedEvent bool) *Song {
	song := &Song{
		Id:         id,
		Group:      group,
		Name:       name,
		Verses:     verses,
		Link:       link,
		ReleasedAt: releaseDate,
		AddedAt:    addedAt,
		UpdatedAt:  updatedAt,

		Events: make([]Event, 0),
	}

	if raiseCreatedEvent {
		song.Events = append(song.Events, newCreatedEvent())
	}

	return song
}

func (s *Song) update(update func()) {
	s.UpdatedAt = time.Now()
	s.Events = append(s.Events, newUpdatedEvent())
	update()
}

func (s *Song) SetGroup(group string) {
	if group == "" {
		return
	}

	s.update(func() {
		s.Group = group
	})
}

func (s *Song) SetName(name string) {
	if name == "" {
		return
	}

	s.update(func() {
		s.Name = name
	})
}

func (s *Song) SetLink(link string) {
	if link == "" {
		return
	}

	s.update(func() {
		s.Link = link
	})
}

func (s *Song) SetReleasedAt(releasedAt time.Time) {
	s.update(func() {
		s.ReleasedAt = releasedAt
	})
}
