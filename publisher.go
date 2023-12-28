package gobroker

type Publisher interface {
	Publish(message []byte, topic string, properties map[string]any) error
}
