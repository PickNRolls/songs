package repositorysong

import (
	"context"
	"errors"
	"songs/src/domain/song"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Id = int

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) Save(s *song.Song) (*song.Song, error) {
	for _, event := range s.Events {
		switch event.Type() {
		case song.CREATED:
			var id Id

			err := r.pool.QueryRow(context.TODO(), `
      insert into songs (
        name,
        group_name,
        verses,
        link,
        released_at,
        added_at,
        updated_at
      ) values ($1, $2, $3, $4, $5, $6, $7) returning id
      `, s.Name, s.Group, s.Verses, s.Link, s.ReleasedAt, s.AddedAt, s.UpdatedAt).Scan(&id)
			if err != nil {
				return nil, err
			}

			return song.New(id, s.Group, s.Name, s.Verses, s.Link, s.ReleasedAt, s.AddedAt, s.UpdatedAt, false), nil

		case song.UPDATED:
			_, err := r.pool.Exec(context.TODO(), `
      update songs set
      
      name = $1,
      group_name = $2,
      verses = $3,
      link = $4,
      released_at = $5,
      added_at = $6,
      updated_at = $7
      
      where id = $8
      `, s.Name, s.Group, s.Verses, s.Link, s.ReleasedAt, s.AddedAt, s.UpdatedAt, s.Id)
			if err != nil {
				return nil, err
			}

			return s, nil
		}
	}

	return nil, errors.New("Unexpected return")
}

func (r *Repository) FindById(id Id) (*song.Song, error) {
	var dest song.Song

	err := r.pool.QueryRow(context.TODO(), `
  select (
    id,
    name,
    group_name,
    verses,
    link,
    released_at,
    added_at,
    updated_at
  ) from songs where id = $1 
  `, id).Scan(&dest)
	if err != nil {
		return nil, err
	}

	return &dest, nil
}

func (r *Repository) DeleteById(id Id) error {
	_, err := r.pool.Exec(context.TODO(), `delete from songs where id = $1`, id)
	return err
}
