package webhooks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriber(t *testing.T) {
	ed := eventData{}
	ims := NewInMemorySubscriber(func(data eventData) {
		ed = data
	})
	
	ims.StartHandler()
	
	err := Subscribe(ims)
	assert.Nil(t, err)

	err = DispatchEvent(EventUserCreated, 100)
	assert.Nil(t, err)
	
	assert.Equal(t, ed.topic, EventUserCreated)
	assert.Equal(t, ed.object, 100)
}