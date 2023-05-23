package telegram

import (
	"article-storage-bot/clients/telegram"
	"article-storage-bot/lib/e"
	"article-storage-bot/storage"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
)

var ErrUnknownCmd = errors.New("unknown command")

const (
	HelpCmd  = "/help"
	RndCmd   = "/rnd"
	StartCmd = "/start"
)

func (m *Manager) doCmd(text, username string, chatId int) error {
	text = strings.TrimSpace(text)

	log.SetPrefix("new: ")
	log.Printf("got command %s from %s", text, username)

	if isAddCmd(text) {
		return m.savePage(text, username, chatId)
	}

	switch text {
	case HelpCmd:
		return m.sendHelp(chatId)
	case RndCmd:
		return m.sendRandom(chatId, username)
	case StartCmd:
		return m.sendHello(chatId)
	default:
		_ = m.tg.SendMessage(chatId, UnknownCmd)
		return ErrUnknownCmd
	}
}

func (m *Manager) sendHelp(chatId int) error {
	return m.tg.SendMessage(chatId, MsgHelp)
}

func (m *Manager) sendHello(chatId int) error {
	return m.tg.SendMessage(chatId, MsgHello)
}

func (m *Manager) sendRandom(chatId int, username string) error {
	sendMsg := messageSender(chatId, m.tg)
	p, err := m.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return e.Wrap("can't send random command", err)
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		_ = sendMsg(MsgNoSavedPages)
	}
	if err = sendMsg(p.URL); err != nil {
		return e.Wrap("can't send message", err)
	}

	return m.storage.Remove(p)
}

func (m *Manager) savePage(pageUrl, username string, chatId int) error {
	p := storage.Page{
		URL:      pageUrl,
		UserName: username,
	}
	sendMsg := messageSender(chatId, m.tg)

	double, err := m.storage.IsExists(&p)
	if err != nil {
		msg := fmt.Sprintf("can't do command save %s", pageUrl)
		return e.Wrap(msg, err)
	}
	if double {
		if err := sendMsg(MsgAlreadyExist); err != nil {
			return e.Wrap("can't send message", err)
		}
	}

	if err = m.storage.Save(&p); err != nil {
		return e.Wrap("can't save page", err)
	}
	if err := sendMsg(MsgSaved); err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func messageSender(chatId int, tg *telegram.Client) func(text string) error {
	return func(text string) error {
		return tg.SendMessage(chatId, text)
	}
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	u, err := url.Parse(text)
	return err != nil && u.Host != ""
}
