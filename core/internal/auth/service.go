package auth

import (
	"errors"
	"time"

	"github.com/ra341/glacier/internal/user"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service struct {
	store   Store
	userSrv *user.Service
	conf    ConfigLoader
}

var (
	ErrTokenExpired       = errors.New("token expired")
	ErrRegistrationClosed = errors.New("registration is closed, contact your admin to create a new account")
	ErrDuplicateUser      = errors.New("username already exists, choose a different username")
	ErrInvalidUserPass    = errors.New("invalid username/password")
)

func New(store Store, userSrv *user.Service, conf ConfigLoader) *Service {
	s := &Service{
		store:   store,
		userSrv: userSrv,
		conf:    conf,
	}
	return s
}

func (s *Service) Register(username, password string, role user.Role, creatingUser *user.User) (err error) {
	if creatingUser == nil && !s.conf().OpenRegistration {
		return ErrRegistrationClosed
	}
	err = s.userSrv.New(username, password, role, creatingUser)
	if err != nil && errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicateUser
	}
	return err
}

func (s *Service) Login(username, password string, sessionType SessionType) (session Session, sessionToken string, refreshToken string, err error) {
	u, err := s.userSrv.GetByUsername(username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by username")
		return Session{}, "", "", ErrInvalidUserPass
	}

	err = user.CheckEncryptedString(password, u.EncryptedPassword)
	if err != nil {
		log.Error().Err(err).Msg("could not decrypt password")
		return Session{}, "", "", ErrInvalidUserPass
	}

	var sess Session
	sess.SessionType = sessionType
	sess.UserId = u.ID
	sess.User = u
	sessionToken, refreshToken = s.GenerateTok(&sess)

	err = s.store.New(&sess)
	if err != nil {
		return Session{}, "", "", err
	}

	return sess, sessionToken, refreshToken, nil
}

func (s *Service) GenerateTok(sess *Session) (plainSession string, plainRefresh string) {
	conf := s.conf()

	plainRefresh = user.GenerateRandomToken(20)
	hashedRefreshToken := user.HashString(plainRefresh)
	refreshExpiry := time.Now().Add(conf.GetRefreshExp())

	sess.HashedRefreshToken = hashedRefreshToken
	sess.RefreshTokenExpiry = refreshExpiry

	plainSession = user.GenerateRandomToken(20)
	hashedSessionToken := user.HashString(plainSession)
	sessionExpiry := time.Now().Add(conf.GetSessionExp())

	sess.HashedSessionToken = hashedSessionToken
	sess.SessionTokenExpiry = sessionExpiry

	return plainSession, plainRefresh
}

func (s *Service) VerifySession(sessionToken string) (session Session, err error) {
	hashedTok := user.HashString(sessionToken)
	token, err := s.store.GetBySessionToken(hashedTok)
	if err != nil {
		return Session{}, err
	}

	err = checkExpiry(token.SessionTokenExpiry)
	if err != nil {
		return token, err
	}

	return token, nil
}

func (s *Service) RefreshSession(refreshToken string) (session Session, sessionTok string, refreshTok string, err error) {
	hashedTok := user.HashString(refreshToken)
	sess, err := s.store.GetByRefreshToken(hashedTok)
	if err != nil {
		return Session{}, "", "", err
	}

	err = checkExpiry(sess.RefreshTokenExpiry)
	if err != nil {
		return Session{}, "", "", err
	}

	sessionTok, refreshTok = s.GenerateTok(&session)
	err = s.store.Edit(&session)
	if err != nil {
		return Session{}, "", "", err
	}

	return session, sessionTok, refreshTok, nil
}

func checkExpiry(expiry time.Time) error {
	now := time.Now()
	isExpired := expiry.Before(now)

	if isExpired {
		return ErrTokenExpired
	}
	return nil
}
