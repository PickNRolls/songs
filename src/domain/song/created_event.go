package song

type CreatedEvent struct {
}

func (event *CreatedEvent) Type() EventType {
	return CREATED
}

func newCreatedEvent() *CreatedEvent {
	return &CreatedEvent{}
}
