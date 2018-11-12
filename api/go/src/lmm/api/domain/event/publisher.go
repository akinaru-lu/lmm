package event

import "context"

type Publisher interface {
	Publish(context.Context, Event) error
}
