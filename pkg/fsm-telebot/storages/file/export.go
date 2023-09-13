package file

type (
	// ChatsStorage in intermediate representation for data in Storage.
	// This type using simple type for key.
	// In any way you can create you custom type and custom provider
	// with compatibility to this type.
	ChatsStorage map[ChatID]UsersStorage
	UsersStorage map[UserID]Record
	Record       struct {
		State string            `json:"state"`
		Data  map[string][]byte `json:"data"`
	}

	ChatID = int64
	UserID = int64
)

func (d *dataCache) export(p Provider) ([]byte, error) {
	if d.raw != nil {
		return d.raw, nil
	}
	bytes, err := p.Encode(d.loaded)
	if err != nil {
		return nil, err
	}
	d.raw = bytes
	return bytes, nil

}

func (r *record) exportData(p Provider) (map[string][]byte, error) {
	if len(r.data) < 1 {
		return nil, nil
	}

	m := make(map[string][]byte)
	for k, d := range r.data {
		data, err := d.export(p)
		if err != nil {
			return nil, err
		}
		m[k] = data
	}

	return m, nil
}

func (s *Storage) dump() (ChatsStorage, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	chats := make(ChatsStorage)
	for key, r := range s.data {
		chat, ok := chats[key.c]
		if !ok {
			chat = make(UsersStorage)
		}

		exportData, err := r.exportData(s.p)
		if err != nil {
			return nil, err
		}

		chat[key.u] = Record{
			State: string(r.state),
			Data:  exportData,
		}
		chats[key.c] = chat
	}
	return chats, nil
}

func (s *Storage) reset(dump ChatsStorage) {
	s.rw.Lock()
	defer s.rw.Unlock()

	for chatId, usersStorage := range dump {
		for userId, r := range usersStorage {
			data := make(map[string]dataCache)
			for key, d := range r.Data {
				data[key] = dataCache{raw: d}
			}

			s.data[newKey(chatId, userId)] = record{
				state: fsm.State(r.State),
				data:  data,
			}
		}
	}
}
