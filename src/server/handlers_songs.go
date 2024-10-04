package server

import (
	"context"
	"encoding/json"
	"net/http"
	repositorysong "songs/src/postgres/repositories/song"
	presentersong "songs/src/presenters/song"
	"songs/src/server/exec"
	servicesong "songs/src/services/song"
	"strconv"
	"time"
)

func (s *Server) registerSongsHandlers(detailsFinder servicesong.DetailsFinder) {
	s.logger.Debug("Registering /songs handlers")

	service := servicesong.New(detailsFinder, repositorysong.New(s.dbPool))
	presenter := presentersong.New(s.dbPool)

	s.handle("POST /songs", exec.Command(service.Add))
	s.handle("DELETE /songs/{id}", exec.Command(service.Delete, exec.WithUrlIdBind, exec.AvoidBodyBind))
	s.handle("PATCH /songs/{id}", exec.Command(service.Update, exec.WithUrlIdBind))

	s.handle("GET /songs", func(w http.ResponseWriter, r *http.Request) {
		var status int
		var err error

		defer func() {
			if status != 0 {
				// TODO: per endpoint handler
				s.logger.Error("GET /songs, " + err.Error())

				// TODO: response presenter middleware
				w.Header().Add("content-type", "application/json")
				w.WriteHeader(status)
				json.NewEncoder(w).Encode(struct {
					Err string `json:"error"`
				}{
					Err: err.Error(),
				})
			}
		}()

		values := r.URL.Query()

		gt, err := time.Parse(time.RFC3339, values.Get("releasedAtGt"))
		if values.Has("releasedAtGt") && err != nil {
			status = http.StatusBadRequest
			return
		}
		releasedAtGt := &gt
		if err != nil {
			releasedAtGt = nil
		}

		lt, err := time.Parse(time.RFC3339, values.Get("releasedAtLt"))
		if values.Has("releasedAtLt") && err != nil {
			status = http.StatusBadRequest
			return
		}
		releasedAtLt := &lt
		if err != nil {
			releasedAtLt = nil
		}

		limit, err := strconv.Atoi(values.Get("limit"))
		if values.Get("limit") != "" && err != nil {
			status = http.StatusBadRequest
			return
		}

		paging, err := presenter.FindByFilter(context.TODO(), presentersong.QueryFindByFilter{
			Group:        values.Get("group"),
			Name:         values.Get("name"),
			Link:         values.Get("link"),
			ReleasedAtGt: releasedAtGt,
			ReleasedAtLt: releasedAtLt,
			Limit:        limit,
			Next:         values.Get("next"),
		})
		if err != nil {
			status = http.StatusInternalServerError
			return
		}

		// TODO: response presenter middleware
		w.Header().Add("content-type", "application/json")
		err = json.NewEncoder(w).Encode(paging)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
	})

	s.handle("GET /songs/{id}/verses", func(w http.ResponseWriter, r *http.Request) {
		var status int
		var err error

		defer func() {
			if status != 0 {
				// TODO: per endpoint handler
				s.logger.Error("GET /songs/{id}/verses, " + err.Error())

				// TODO: response presenter middleware
				w.Header().Add("content-type", "application/json")
				w.WriteHeader(status)
				json.NewEncoder(w).Encode(struct {
					Err string `json:"error"`
				}{
					Err: err.Error(),
				})
			}
		}()

		values := r.URL.Query()
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			status = http.StatusBadRequest
			return
		}

		number, err := strconv.Atoi(values.Get("number"))
		if values.Get("number") != "" && err != nil {
			status = http.StatusBadRequest
			return
		}

		size, err := strconv.Atoi(values.Get("size"))
		if values.Get("size") != "" && err != nil {
			status = http.StatusBadRequest
			return
		}

		paging, err := presenter.FindVerses(context.TODO(), presentersong.QueryFindVerses{
			Id:     id,
			Number: number,
			Size:   size,
		})
		if err != nil {
			status = http.StatusInternalServerError
			return
		}

		// TODO: response presenter middleware
		w.Header().Add("content-type", "application/json")
		err = json.NewEncoder(w).Encode(paging)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
	})
}
