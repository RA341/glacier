package auth

import "time"

const (
	Day   = time.Hour * 24
	Month = Day * 30
)

type ConfigLoader func() *Config

type Config struct {
	Disable          bool `yaml:"disable" env:"AUTH_DISABLE" default:"false" help:"disable authentication"`
	OpenRegistration bool `yaml:"openRegistration" env:"AUTH_OPEN_REGISTRATION" default:"false" help:"open registration for anyone to signup"`

	MaxConcurrentSessions int `yaml:"maxConcurrentSessions" env:"AUTH_MAX_SESSIONS" default:"8" help:"maximum number of logged in sessions per user"`
	SessionExpiryInDays   int `yaml:"sessionExpiryInDays" env:"AUTH_SESSION_EXPIRY" default:"1" help:"time limit for a session token"`
	RefreshExpiryInMonths int `yaml:"refreshExpiryInMonths" env:"AUTH_REFRESH_EXPIRY" default:"12" help:"time limit for a refresh token"`

	OidcEnable       bool   `yaml:"OIDCEnable" env:"AUTH_OIDC_ENABLE" default:"false" help:"enable OIDC support"`
	OidcIssuerUrl    string `yaml:"OIDCIssuerURL" env:"AUTH_OIDC_ISSUER" default:"TODO" help:"url for your OIDC issuer"`
	OidcClientId     string `yaml:"OIDCClientID" env:"AUTH_OIDC_CLIENT_ID" default:"TODO" help:"client id for OIDC" hide:"true"`
	OidcClientSecret string `yaml:"OIDCClientSecret" env:"AUTH_OIDC_CLIENT_SECRET" default:"TODO" help:"client secret for OIDC" hide:"true"`
	OidcRedirectUrl  string `yaml:"OIDCRedirectURL" env:"AUTH_OIDC_REDIRECT_URL" default:"TODO" help:"redirect url for OIDC"`
}

func (c *Config) GetSessionExp() time.Duration {
	return Day * time.Duration(c.SessionExpiryInDays)
}

func (c *Config) GetRefreshExp() time.Duration {
	return Month * time.Duration(c.SessionExpiryInDays)
}
