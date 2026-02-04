package secrets

import "github.com/zalando/go-keyring"

type StoreKeyring struct {
	name string
}

func NewKeyringStore(name string) Store {
	return &StoreKeyring{
		name: name,
	}
}

func (s *StoreKeyring) Add(key, value string) error {
	return keyring.Set(s.name, key, value)
}

func (s *StoreKeyring) Get(key string) (string, error) {
	return keyring.Get(s.name, key)
}
