package song

type EventType string

const (
	CREATED EventType = "CREATED"
	UPDATED EventType = "UPDATED"
)

type Event interface {
	Type() EventType
}
