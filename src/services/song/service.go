package servicesong

import (
	"context"
	"errors"
	"log"
	"songs/src/domain/song"
	"strings"
	"time"
)

type Service struct {
	detailsFinder DetailsFinder
	repository    Repository
}

func New(detailsFinder DetailsFinder, repository Repository) *Service {
	return &Service{
		detailsFinder: detailsFinder,
		repository:    repository,
	}
}

func (s *Service) Add(parent context.Context, command CommandAdd) (*song.Song, error) {
	if command.Group == "" || command.Name == "" {
		return nil, errors.New("Group and Name must be filled")
	}

	details, err := s.detailsFinder.Find(command.Group, command.Name)
	if err != nil {
		return nil, err
	}

	verses := strings.Split(details.Text(), "\n\n")
	song := song.New(0, command.Group, command.Name, verses, details.Link(), details.ReleaseDate(), time.Now(), time.Now(), true)

	song, err = s.repository.Save(song)
	if err != nil {
		return nil, err
	}

	return song, nil
}

func (s *Service) Update(parent context.Context, command CommandUpdate) (*song.Song, error) {
	song, err := s.repository.FindById(command.Id)
	if err != nil {
		return nil, err
	}

	song.SetGroup(command.Group)
	song.SetName(command.Name)
	song.SetLink(command.Link)
	if command.ReleasedAt != nil {
		song.SetReleasedAt(*command.ReleasedAt)
	}

	song, err = s.repository.Save(song)
	if err != nil {
		return nil, err
	}

	return song, err
}

func (s *Service) Delete(parent context.Context, command CommandDelete) (Id, error) {
	log.Println(command)
	err := s.repository.DeleteById(command.Id)
	if err != nil {
		return 0, nil
	}

	return command.Id, nil
}
