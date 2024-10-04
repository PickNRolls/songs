package exec

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type binding[C any] func(*C, *http.Request) error

func (this binding[C]) chain(b binding[C]) binding[C] {
	return func(command *C, r *http.Request) error {
		err := this(command, r)
		if err != nil {
			return err
		}
		return b(command, r)
	}
}

type Opt[C any] interface {
	Binding() binding[C]
	AvoidBodyBind() bool
}

type OptFn[C any] func() Opt[C]

func noopBind[C any](command *C, r *http.Request) error {
	return nil
}

func jsonBind[C any](command *C, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(command)
	defer r.Body.Close()

	return err
}

func urlIdBind[C any, PC PointerCanSetId[C]](command *C, r *http.Request) error {
	idstr := r.PathValue("id")
	if idstr == "" {
		return errors.New("Id must me filled")
	}

	id, err := strconv.Atoi(idstr)
	if err != nil {
		return err
	}

	PC(command).SetId(id)
	return nil
}

type urlIdBindOpt[C any, PC PointerCanSetId[C]] struct{}

func (this *urlIdBindOpt[C, PC]) Binding() binding[C] {
	return urlIdBind[C, PC]
}

func (this *urlIdBindOpt[C, PC]) AvoidBodyBind() bool { return false }

type PointerCanSetId[C any] interface {
	SetId(id int)
	*C
}

func WithUrlIdBind[C any, PC PointerCanSetId[C]]() Opt[C] {
	return &urlIdBindOpt[C, PC]{}
}

type avoidBodyBind[C any] struct {
}

func (this *avoidBodyBind[C]) Binding() binding[C] { return nil }
func (this *avoidBodyBind[C]) AvoidBodyBind() bool { return true }

func AvoidBodyBind[C any]() Opt[C] {
	return &avoidBodyBind[C]{}
}

func Command[C any, T any](
	serviceMethod func(context.Context, C) (T, error),
	opts ...OptFn[C],
) func(w http.ResponseWriter, r *http.Request) {
	var bind binding[C] = noopBind
	avoidBodyBind := false

	for _, optFn := range opts {
		opt := optFn()

		if opt.Binding() != nil {
			bind = bind.chain(opt.Binding())
		}

		if opt.AvoidBodyBind() {
			avoidBodyBind = true
		}
	}

	if !avoidBodyBind {
		bind = bind.chain(jsonBind[C])
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var command C
		err := bind(&command, r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		ctx := context.TODO()
		out, err := serviceMethod(ctx, command)

		w.Header().Add("content-type", "application/json")
		err = json.NewEncoder(w).Encode(out)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
