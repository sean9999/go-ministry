package main

type memPersister map[key]Message

func (m memPersister) Set(k key, msg Message) {
	m[k] = msg
}

func (m memPersister) Has(k key) bool {
	_, has := m[k]
	return has
}

func (m memPersister) Get(k key) (Message, error) {
	if !m.Has(k) {
		return NilMessage, ErrNoRecord
	}
	return m[k], nil
}

func (m memPersister) Delete(k key) error {
	if !m.Has(k) {
		return ErrNoRecord
	}
	delete(m, k)
	return nil
}
