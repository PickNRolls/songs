package dto

type CursorPagingDTO[T any] struct {
	Data []T    `json:"data"`
	Next string `json:"next,omitempty"`
}

type OffsetPagingDTO[T any] struct {
	Data []T `json:"data"`
}
