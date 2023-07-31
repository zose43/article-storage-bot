package telegram

import (
	"article-storage-bot/clients/telegram"
	"article-storage-bot/events"
	"article-storage-bot/lib/e"
	"article-storage-bot/storage"
	"errors"
	"fmt"
)

var (
	UnknownEventType = errors.New("unknown event type")
	UnknownMetaType  = errors.New("unknown meta type")
)

type Manager struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func (m *Manager) Processor(event events.Event) error {
	switch event.Type {
	case events.Message:
		return m.processMessage(event)
	default:
		return UnknownEventType
	}
}

type Meta struct {
	ChatId   int
	Username string
}

func (m *Manager) Fetch(limit int) ([]events.Event, error) {
	updates, err := m.tg.Updates(m.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't fetch updates", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, v := range updates {
		res = append(res, event(v))
	}

	m.offset = updates[len(updates)-1].Id + 1
	return res, nil
}

func (m *Manager) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		msg := fmt.Sprintf("can't process message %v", event.Meta)
		return e.Wrap(msg, err)
	}
	if err = m.doCmd(event.Text, meta.Username, meta.ChatId); err != nil {
		msg := fmt.Sprintf("can't process message %s", event.Text)
		return e.Wrap(msg, err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	meta, ok := event.Meta.(Meta)
	if ok {
		return meta, nil
	}
	return Meta{}, e.Wrap("can't get meta", UnknownMetaType)
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchMessage(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatId:   upd.Message.Chat.Id,
			Username: upd.Message.User.Username,
		}
	}
	return res
}

func fetchMessage(upd telegram.Update) string {
	if upd.Message != nil {
		return upd.Message.Text
	}
	return ""
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	}
	return events.Unknown
}

func NewManager(tg *telegram.Client, storage storage.Storage) *Manager {
	return &Manager{tg: tg, storage: storage}
}
