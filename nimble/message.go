package nimble

import "github.com/google/uuid"

type Message struct {
	Id      string
	Headers map[string]any
	Payload any
}

func NewMessage(payload any) *Message {
	return &Message{
		Id:      uuid.NewString(),
		Headers: make(map[string]any),
		Payload: payload,
	}
}

func NewMessageWH(payload any, headers map[string]any) *Message {
	return &Message{
		Id:      uuid.NewString(),
		Headers: headers,
		Payload: payload,
	}
}
