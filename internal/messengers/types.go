package messengers

type Recipient struct {
	UUID  string
	Email string
}

type Message struct {
	Title      string
	Body       string
	Recipients []Recipient
}

type Messenger interface {
	Send(message Message) error
	Close()
}
