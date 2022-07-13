package models

type Message struct {
	Id     int    `json:"id,omitempty"`
	Sender string `json:"sender"`
	Target string `json:"target"`
	Body   string `json:"body"`
}

func NewMessage(sender, target string, body string) *Message {
	return &Message{
		Sender: sender,
		Target: target,
		Body:   body,
	}
}
