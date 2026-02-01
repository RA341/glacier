package auth

import (
	"time"

	"github.com/ra341/glacier/internal/user"
	"gorm.io/gorm"
)

type Store interface {
	New(token *Session) error
	Edit(token *Session) error
	Delete(token *Session) error

	List() ([]Session, error)
	ListByUser(userId uint) ([]Session, error)

	GetBySessionId(sessionId uint) (Session, error)
	GetBySessionToken(token string) (Session, error)
}

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=SessionType -output=enum_session_type.go
type SessionType int

const (
	// Web for direct glacier web UI access
	Web SessionType = iota

	// Frost for frost clients
	Frost
)

type Session struct {
	gorm.Model

	UserId uint
	User   user.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	HashedRefreshToken string `gorm:"uniqueIndex"`
	RefreshTokenExpiry time.Time
	HashedSessionToken string `gorm:"uniqueIndex"`
	SessionTokenExpiry time.Time

	SessionType SessionType
}
