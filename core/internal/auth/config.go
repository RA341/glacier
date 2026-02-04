package auth

import "time"

const (
	Day   = time.Hour * 24
	Month = Day * 30
)

type ConfigLoader func() *Config

type Config struct {
	MaxConcurrentSessions int  `yaml:"maxConcurrentSessions" env:"MAX_SESSIONS" default:"8" help:"maximum number of logged in sessions per user"`
	Disable               bool `yaml:"disable" env:"AUTH_DISABLE" default:"false" help:"disable authentication"`
	OpenRegistration      bool `yaml:"openRegistration" env:"AUTH_OPEN_REGISTRATION" default:"false" help:"open registration for anyone to signup"`
	SessionExpiryInDays   int  `yaml:"sessionExpiryInDays" env:"AUTH_SESSION_EXPIRY" default:"1" help:"time validity for a session"`
	RefreshExpiryInMonths int  `yaml:"refreshExpiryInMonths" env:"AUTH_REFRESH_EXPIRY" default:"12" help:"time validity for a refresh"`
}

func (c *Config) GetSessionExp() time.Duration {
	return Day * time.Duration(c.SessionExpiryInDays)
}

func (c *Config) GetRefreshExp() time.Duration {
	return Month * time.Duration(c.SessionExpiryInDays)
}
