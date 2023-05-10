package vk

import (
	"VKBot/client/vk"
	"VKBot/events"
	"VKBot/lib/e"
	"errors"
)

type Processor struct {
	Vk  *vk.Client
	Ts  string
	Key string
}

type Meta struct {
	MessageID int
	UserID    int
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

const (
	eventTypeMessege = "message_new"
)

func New(client *vk.Client) *Processor {
	return &Processor{
		Vk: client,
	}
}

func (p *Processor) Fetch() ([]events.Event, error) {
	res, err := p.Vk.GetLongPool(p.Key, p.Ts)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(res.Updates) == 0 {
		return nil, nil
	}

	resEvent := make([]events.Event, 0, len(res.Updates))

	for _, u := range res.Updates {
		if u.EventType == eventTypeMessege {
			resEvent = append(resEvent, event(u))
		}
	}

	p.Ts = res.Ts

	return resEvent, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.MessageID, meta.UserID); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd vk.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			MessageID: upd.Object.Message.Id,
			UserID:    upd.Object.Message.UserID,
		}
	}

	return res
}

func fetchText(upd vk.Update) string {
	return upd.Object.Message.Text
}

func fetchType(upd vk.Update) events.Type {
	return events.Message
}
