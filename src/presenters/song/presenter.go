package presentersong

import (
	"context"
	"errors"
	"songs/src/lib/cursor"
	"songs/src/lib/logger"
	"songs/src/presenters/dto"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Presenter struct {
	pool   *pgxpool.Pool
	logger logger.Logger
}

func New(pool *pgxpool.Pool) *Presenter {
	return &Presenter{
		pool:   pool,
		logger: &logger.StdoutLogger{},
	}
}

func (p *Presenter) mapDests(dests []dest) []*SongDTO {
	data := make([]*SongDTO, len(dests))

	for i, dest := range dests {
		data[i] = &SongDTO{
			Id:         dest.Id,
			Group:      dest.Group,
			Name:       dest.Name,
			Verses:     dest.Verses,
			Link:       dest.Link,
			ReleasedAt: dest.ReleasedAt,
			AddedAt:    dest.AddedAt,
			UpdatedAt:  dest.UpdatedAt,
		}
	}

	return data
}

func (p *Presenter) FindByFilter(parent context.Context, query QueryFindByFilter) (*dto.CursorPagingDTO[*SongDTO], error) {
	limit := query.Limit
	if limit == 0 {
		limit = 20
	}

	sql := `
  select
    id,
    name,
    group_name,
    verses,
    link,
    released_at,
    added_at,
    updated_at
  from songs`

	args := []any{}
	wheres := []string{}

	// TODO: query builder
	if query.Group != "" {
		wheres = append(wheres, `group_name = $`+strconv.Itoa(len(args)+1))
		args = append(args, query.Group)
	}

	if query.Name != "" {
		wheres = append(wheres, `name = $`+strconv.Itoa(len(args)+1))
		args = append(args, query.Name)
	}

	if query.Link != "" {
		wheres = append(wheres, `link = $`+strconv.Itoa(len(args)+1))
		args = append(args, query.Link)
	}

	if query.ReleasedAtGt != nil {
		wheres = append(wheres, `released_at > $`+strconv.Itoa(len(args)+1))
		args = append(args, *query.ReleasedAtGt)
	}

	if query.ReleasedAtLt != nil {
		wheres = append(wheres, `released_at < $`+strconv.Itoa(len(args)+1))
		args = append(args, *query.ReleasedAtLt)
	}

	if query.Next != "" {
		cursor, err := cursor.FromString(query.Next)
		if err != nil {
			return nil, err
		}

		a1 := strconv.Itoa(len(args) + 1)
		a2 := strconv.Itoa(len(args) + 2)
		wheres = append(wheres, `(id, added_at) < ($`+a1+`, $`+a2+`)`)
		args = append(args, cursor.Id, cursor.Time)
	}

	if len(wheres) == 0 {
		return nil, errors.New("presentersong.FindByFilter: expected at least 1 filter, got 0")
	}

	where := strings.Join(wheres, " and ")
	sql += " where " + where

	sql += `
  order by (id, added_at) desc
  limit $` + strconv.Itoa(len(args)+1)
	args = append(args, limit)

	p.logger.Debug("presentersong.FindByFilter SQL:")
	p.logger.Debug(sql)

	rows, err := p.pool.Query(parent, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dests := []dest{}
	for rows.Next() {
		var dest dest
		err = rows.Scan(
			&dest.Id,
			&dest.Name,
			&dest.Group,
			&dest.Verses,
			&dest.Link,
			&dest.ReleasedAt,
			&dest.AddedAt,
			&dest.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		dests = append(dests, dest)
	}

	if rows.Err() != nil {
		return nil, err
	}

	out := &dto.CursorPagingDTO[*SongDTO]{}
	if len(dests) == 0 {
		return out, nil
	}

	out.Data = p.mapDests(dests)
	last := dests[len(dests)-1]
	out.Next = cursor.New(last.Id, last.AddedAt).String()

	return out, nil
}

func (p *Presenter) FindVerses(parent context.Context, query QueryFindVerses) (*dto.OffsetPagingDTO[string], error) {
	size := query.Size
	if size == 0 {
		size = 1
	}

	sql := `
  select
    verses[$1:$2]
  from songs
  where id = $3
  `

	args := []any{query.Number + 1, query.Number + size, query.Id}
	var verses []string
	err := p.pool.QueryRow(parent, sql, args...).Scan(&verses)
	if err != nil {
		return nil, err
	}

	out := &dto.OffsetPagingDTO[string]{}
	if len(verses) == 0 {
		return out, nil
	}

	out.Data = verses

	return out, nil
}
