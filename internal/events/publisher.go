package events

type Publisher interface {
	Publisher(eventType string, payload interface{}, metadata map[string]string) error
	Close() error
}
