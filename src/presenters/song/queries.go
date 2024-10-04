package presentersong

import "time"

type QueryFindByFilter struct {
	Group        string
	Name         string
	Link         string
	ReleasedAtGt *time.Time
	ReleasedAtLt *time.Time
	Limit        int
	Next         string
}

type QueryFindVerses struct {
  Id int
  Number int
  Size int 
}

