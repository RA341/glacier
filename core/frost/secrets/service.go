package secrets

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

const SecretKeySession = "secrets-session"
const SecretKeyRefresh = "secrets-refresh"

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) AddSession(session string, refresh string) error {
	err := s.store.Add(SecretKeySession, session)
	if err != nil {
		return fmt.Errorf("could not store session key: %w", err)
	}

	err = s.store.Add(SecretKeyRefresh, refresh)
	if err != nil {
		return fmt.Errorf("could not store refresh key: %w", err)
	}

	return nil
}

func (s *Service) GetSession() (session string, refresh string) {
	session = s.Get(SecretKeySession)
	refresh = s.Get(SecretKeyRefresh)

	return session, refresh
}

func (s *Service) Get(key string) string {
	get, err := s.store.Get(key)
	if err != nil {
		log.Warn().Err(err).Str("key", key).Msg("could not get secret")
	}

	return get
}
