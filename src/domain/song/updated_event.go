package song

type UpdatedEvent struct {
}

func (event *UpdatedEvent) Type() EventType {
	return UPDATED
}

func newUpdatedEvent() *UpdatedEvent {
	return &UpdatedEvent{}
}
