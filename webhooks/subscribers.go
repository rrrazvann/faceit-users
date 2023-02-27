package webhooks

const (
	EventUserCreated = "USER_CREATED"
	EventUserUpdated = "USER_UPDATED"
	EventUserDeleted = "USER_DELETED"
)

type subscriber interface {
	Notify(topic string, object any) error
}

var subscribers = []subscriber{}

func Subscribe(s subscriber) error {
	subscribers = append(subscribers, s)

	return nil
}

// todo: implement as async notifications + retry system
func DispatchEvent(topic string, object any) error {
	for _, subscriber := range subscribers {
		err := subscriber.Notify(topic, object)
		if err != nil {
			return err
		}
	}

	return nil
}
