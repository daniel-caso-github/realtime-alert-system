package event

import "context"

// Publisher defines the interface for publishing events.
type Publisher interface {
	Publish(ctx context.Context, event *Event) error
	PublishToStream(ctx context.Context, stream string, event *Event) error
}

// Subscriber defines the interface for subscribing to events.
type Subscriber interface {
	Subscribe(ctx context.Context, stream string, group string, handler Handler) error
	Unsubscribe() error
}

// Handler defines the interface for handling events.
type Handler func(ctx context.Context, event *Event) error

// Bus combines Publisher and Subscriber interfaces.
type Bus interface {
	Publisher
	Subscriber
}

// Stream names.
const (
	StreamAlerts        = "alerts"
	StreamNotifications = "notifications"
	StreamDeadLetter    = "dead-letter"
)

// Consumer group names.
const (
	GroupAlertProcessors      = "alert-processors"
	GroupNotificationSenders  = "notification-senders"
	GroupDeadLetterProcessors = "dead-letter-processors"
)
