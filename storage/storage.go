package storage

import (
	"article-storage-bot/lib/e"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
)

var ErrNoSavedPages = errors.New("no one saved page")

type Storage interface {
	Save(p *Page) error
	PickRandom(uName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	const msg = "Can't calculate hash"
	h := md5.New()
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap(msg, err)
	}
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap(msg, err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
