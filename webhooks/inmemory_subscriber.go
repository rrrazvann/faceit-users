package webhooks

type eventData struct {
	topic  string
	object any
}

type InMemorySubscriber struct {
	channel chan eventData
	handler func(data eventData)
}

func NewInMemorySubscriber(handler func(data eventData)) *InMemorySubscriber {
	return &InMemorySubscriber{
		channel: make(chan eventData),
		handler: handler,
	}
}

func (w InMemorySubscriber) StartHandler() {
	go func() {
		for data := range w.channel {
			w.handler(data)
		}
	}()
}

func (w InMemorySubscriber) Notify(topic string, object any) error {
	w.channel <- eventData{
		topic:  topic,
		object: object,
	}

	return nil
}
