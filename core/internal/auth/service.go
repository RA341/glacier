package auth

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/internal/user"
	"golang.org/x/oauth2"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service struct {
	store   Store
	userSrv *user.Service
	conf    ConfigLoader

	oidcProvider *oidc.Provider
	oauthConfig  *oauth2.Config
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

	config := conf()
	if config.OIDCEnable {
		ctx := getOidcContext(nil)

		provider, err := oidc.NewProvider(ctx, config.OIDCIssuerURL)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to query OIDC provider")
		}

		oauth2Config := &oauth2.Config{
			ClientID:     config.OIDCClientID,
			ClientSecret: config.OIDCClientSecret,
			RedirectURL:  config.OIDCRedirectURL,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		}

		s.oidcProvider = provider
		s.oauthConfig = oauth2Config
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

	return s.createSession(&u, sessionType)
}

func (s *Service) createSession(u *user.User, sessionType SessionType) (session Session, sessionToken string, refreshToken string, err error) {
	var sess Session
	sess.SessionType = sessionType
	sess.UserId = u.ID
	sessionToken, refreshToken = s.GenerateTok(&sess)

	err = s.store.New(&sess)
	if err != nil {
		return Session{}, "", "", err
	}

	return sess, sessionToken, refreshToken, nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// OIDC stuff

// useful to set disable self-signed cert warnings for dev
// pass nil client to use default context.background
func getOidcContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if info.IsDev() {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		customClient := &http.Client{Transport: tr}
		ctx = oidc.ClientContext(ctx, customClient)
	}
	return ctx
}

func (s *Service) GetOIDCLoginURL(state string) (string, error) {
	if !s.conf().OIDCEnable {
		return "", fmt.Errorf("OIDC is disabled")
	}

	return s.oauthConfig.AuthCodeURL(state), nil
}

func (s *Service) LoginOIDC(
	ctx context.Context,
	code string,
	sessionType SessionType,
) (session Session, sessionToken string, refreshToken string, err error) {
	ctx = getOidcContext(ctx)
	oauth2Token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return Session{}, "", "", fmt.Errorf("failed to exchange token: %w", err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return Session{}, "", "", fmt.Errorf("no id_token field in oauth2 token")
	}

	verifier := s.oidcProvider.Verifier(&oidc.Config{ClientID: s.conf().OIDCClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return Session{}, "", "", fmt.Errorf("failed to verify ID Token: %w", err)
	}

	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
	}

	err = idToken.Claims(&claims)
	if err != nil {
		return Session{}, "", "", fmt.Errorf("failed to parse claims: %w", err)
	}

	u, err := s.userSrv.GetByEmail(claims.Email)
	if err != nil {
		return Session{}, "", "", fmt.Errorf("failed to get user by email: %w", err)
	}

	return s.createSession(&u, sessionType)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// utils

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
