package cursor

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Cursor struct {
	Id   int
	Time time.Time
}

func New(id int, time time.Time) *Cursor {
	return &Cursor{Id: id, Time: time}
}

const separator = "/"

func fromStringErr(message string) error {
	return errors.New("cursor.FromString: " + message)
}

func FromString(s string) (*Cursor, error) {
	if s == "" {
		return nil, fromStringErr("expected string not to be zero")
	}

	split := strings.Split(s, separator)
	if len(split) != 2 {
		return nil, fromStringErr("invalid cursor string")
	}

	idstr := split[0]
	timestr := split[1]

	id, err := strconv.Atoi(idstr)
	if err != nil {
		return nil, fromStringErr("invalid cursor string")
	}

	time, err := time.Parse(time.RFC3339, timestr)
	if err != nil {
		return nil, fromStringErr("invalid cursor string")
	}

	return &Cursor{
		Id:   id,
		Time: time,
	}, nil
}

func (c *Cursor) String() string {
	return strconv.Itoa(c.Id) + separator + c.Time.Format(time.RFC3339)
}
