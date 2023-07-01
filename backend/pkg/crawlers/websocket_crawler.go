package crawlers

type Message struct {
	MessageDict *map[string]string
	MessageType *string
	MessageErr  *error
}

type WebsocketCrawlers interface {
	Connect() error
	Close() error
	WriteMessage(messages []string) error
	ReadMessage() <-chan Message
}
