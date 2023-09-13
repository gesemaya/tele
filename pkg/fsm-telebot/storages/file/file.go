// Package file implemented storage what stored in files and can restore its state from files.
//
// For universality, saving data in the format is moved to the Provider interface.
// This allows you not to think at the storage level about how the data will be stored.
package file

import (
	"errors"
	"fmt"
	"github.com/gesemaya/sniper/pkg/fsm-telebot"
	"io"
	"sync"


)

type WriterFunc func() (io.WriteCloser, error)

// Provider saves data to files (or streams). With custom format.
// Some providers (json, gob, base64) are already implemented in
// the `provider` sub-package.
type Provider interface {
	ProviderName() string
	Save(w io.Writer, data ChatsStorage) error
	Read(r io.Reader) (ChatsStorage, error)
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}

// chatKey represents  pair {c: chat id, u: user id}
type chatKey struct {
	c, u int64
}

func newKey(chat, user int64) chatKey {
	return chatKey{
		c: chat,
		u: user,
	}
}

type record struct {
	state fsm.State
	data  map[string]dataCache
}

// dataCache stores data in two variants.
// Decoded in loaded
// and raw.
type dataCache struct {
	// loaded decoded content from raw via provider
	// see dataCache.get in ./internal
	loaded any
	// raw content from file.
	raw []byte
}

// Storage is file storage. In run time data storages in RAM.
// On save (Close) and restore (Init) data will edit
// by result of Provider.
//
// For safe format operations storage serialization
// to special format - ChatsStorage.
type Storage struct {
	rw       sync.RWMutex
	data     map[chatKey]record
	p        Provider
	writerFn WriterFunc
}

func NewStorage(p Provider, writerFn WriterFunc) *Storage {
	return &Storage{p: p, writerFn: writerFn, data: make(map[chatKey]record)}
}

// Init storage set from readr.
//
// If the reader is equal to "nil",
// function will finish without an error.
// But the storage state will not change too.
func (s *Storage) Init(r io.Reader) error {
	if r == nil {
		return nil
	}
	dump, err := s.p.Read(r)
	if err != nil {
		return err
	}
	s.reset(dump)
	return nil
}

func (s *Storage) GetState(chatId, userId int64) (fsm.State, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.data[newKey(chatId, userId)].state, nil
}

func (s *Storage) SetState(chatId, userId int64, state fsm.State) error {
	s.do(chatId, userId, func(r *record) {
		r.state = state
	})
	return nil
}

func (s *Storage) ResetState(chatId, userId int64, withData bool) error {
	s.do(chatId, userId, func(r *record) {
		r.state = ""
		if withData {
			for key := range r.data {
				delete(r.data, key)
			}
		}
	})
	return nil
}

func (s *Storage) UpdateData(chatId, userId int64, key string, data any) error {
	s.do(chatId, userId, func(r *record) {
		r.updateData(key, data)
	})
	return nil
}

func (s *Storage) GetData(chatId, userId int64, key string, to any) error {
	s.rw.RLock()
	defer s.rw.RUnlock()
	d, ok := s.data[newKey(chatId, userId)].data[key]
	if !ok {
		return fsm.ErrNotFound
	}

	return d.get(to, s.p)
}

// Close saves storage data to writer from writer function.
//
// Also, the method closes writer, minimum once time.
func (s *Storage) Close() (err error) {
	w, err := s.writerFn()
	if err != nil {
		return err
	}

	defer func(w io.WriteCloser) {
		errClose := w.Close()
		err = errors.Join(err, errClose)
	}(w)

	err = s.save(w)
	return
}

// SaveTo saves storage data to writer.
// You can use this method to create dumps in runtime.
func (s *Storage) SaveTo(w io.Writer) error {
	return s.save(w)
}

func (s *Storage) save(w io.Writer) error {
	dump, err := s.dump()
	if err != nil {
		return err
	}

	return s.p.Save(&wrapWriter{w}, dump)
}

// wrapWriter protects base wrapper from type assertions.
type wrapWriter struct {
	w io.Writer
}

func (w *wrapWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

type ProviderError struct {
	ProviderType string
	Operation    string
	Err          error
}

func (e ProviderError) Unwrap() error { return e.Err }
func (e ProviderError) Error() string {
	return fmt.Sprintf("fsm-telebot/storage/file/provider: %s: %s: %v", e.ProviderType, e.Operation, e.Err)
}
