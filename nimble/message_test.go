package nimble

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessage_NewMessage_CreatesMessageWithPayloadAndEmptyHeaders(t *testing.T) {
	// given
	payload := "test"
	// when
	msg := NewMessage(payload)
	// then
	assert.NotNil(t, msg)
	assert.NotEmpty(t, msg.Id)
	assert.Equal(t, payload, msg.Payload)
	assert.NotNil(t, msg.Headers)
	assert.Empty(t, msg.Headers)
}

func TestMessage_NewMessageWH_CreatesMessageWithPayloadAndProvidedHeaders(t *testing.T) {
	// given
	payload := 123
	headers := map[string]any{"foo": "bar"}
	// when
	msg := NewMessageWH(payload, headers)
	// then
	assert.NotNil(t, msg)
	assert.NotEmpty(t, msg.Id)
	assert.Equal(t, payload, msg.Payload)
	assert.Equal(t, headers, msg.Headers)
}

func TestMessage_NewMessageWH_HandlesNilHeaders(t *testing.T) {
	// given
	payload := []int{1, 2, 3}
	// when
	msg := NewMessageWH(payload, nil)
	// then
	assert.NotNil(t, msg)
	assert.NotEmpty(t, msg.Id)
	assert.Equal(t, payload, msg.Payload)
	assert.Nil(t, msg.Headers)
}

func TestMessage_NewMessage_DifferentIdsForDifferentMessages(t *testing.T) {
	// given
	payload1 := "a"
	payload2 := "b"
	// when
	msg1 := NewMessage(payload1)
	msg2 := NewMessage(payload2)
	// then
	assert.NotEqual(t, msg1.Id, msg2.Id)
}

func TestMessage_NewMessageWH_DifferentIdsForDifferentMessages(t *testing.T) {
	// given
	payload1 := "a"
	payload2 := "b"
	headers := map[string]any{"x": 1}
	// when
	msg1 := NewMessageWH(payload1, headers)
	msg2 := NewMessageWH(payload2, headers)
	// then
	assert.NotEqual(t, msg1.Id, msg2.Id)
}
