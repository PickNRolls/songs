package servicesong

import "time"

type CommandAdd struct {
	Group string `json:"group"`
	Name  string `json:"song"`
}

type CommandDelete struct {
	Id int `json:"id"`
}

func (c *CommandDelete) SetId(id int) {
	c.Id = id
}

type CommandUpdate struct {
	Id         int        `json:"id"`
	Group      string     `json:"group"`
	Name       string     `json:"name"`
	ReleasedAt *time.Time `json:"releasedAt"`
	Link       string     `json:"link"`
}

func (c *CommandUpdate) SetId(id int) {
	c.Id = id
}
