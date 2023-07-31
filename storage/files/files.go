package files

import (
	"article-storage-bot/lib/e"
	"article-storage-bot/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const perm = 0755

type Storage struct {
	basePath string
}

func (s Storage) Save(p *storage.Page) (err error) {
	fp := filepath.Join(s.basePath, p.UserName)
	if err := os.MkdirAll(fp, perm); err != nil {
		return e.Wrap("can't save page", err)
	}
	fName, err := filename(p)
	if err != nil {
		return err
	}
	fp = filepath.Join(fp, fName)
	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }() // ignore

	if err := gob.NewEncoder(file).Encode(p); err != nil {
		return e.Wrap("can't encode page", err)
	}
	return nil
}

func (s Storage) PickRandom(uName string) (*storage.Page, error) {
	fp := filepath.Join(s.basePath, uName)
	files, err := os.ReadDir(fp)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.NewSource(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]
	return s.decodePage(filepath.Join(fp, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	const msg = "can't remove page"
	fName, err := filename(p)
	if err != nil {
		return e.Wrap(msg, err)
	}
	fp := filepath.Join(s.basePath, p.UserName, fName)
	if err := os.Remove(fp); err != nil {
		return e.Wrap(fmt.Sprintf("%s %s", msg, fp), err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	const msg = "can't check exists page"
	fName, err := filename(p)
	if err != nil {
		return false, e.Wrap(msg, err)
	}

	fp := filepath.Join(s.basePath, p.UserName, fName)
	switch _, err := os.Stat(fp); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return true, e.Wrap(fmt.Sprintf("%s %s", msg, fp), err)
	}
	return true, nil
}

func NewStorage(bp string) Storage {
	return Storage{basePath: bp}
}

func filename(p *storage.Page) (string, error) {
	return p.Hash()
}

func (s Storage) decodePage(fp string) (*storage.Page, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't decode page %s", fp), err)
	}
	defer func() { _ = f.Close() }()
	p := new(storage.Page)
	if err := gob.NewDecoder(f).Decode(p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	return p, nil
}
