package user

import (
	"gorm.io/gorm"
)

type Store interface {
	GetByID(id uint) (User, error)
	GetByUsername(username string) (User, error)
	GetByEmail(username string) (User, error)

	New(user *User) error
	Edit(user *User) error
	Delete(id uint) error
	List(q string) ([]User, error)
}

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=Role -output=enum_user_role.go
type Role int

func (r Role) IsOmnissiah() bool {
	return r == Omnissiah
}

func (r Role) IsLessThenTechPriest() bool {
	return r < TechPriest
}

const (
	// Omnissiah Superuser: the first user account created by system the highest privilege can never be deleted
	Omnissiah Role = iota
	// Magos Admin: admin level priv
	Magos
	// TechPriest normal base role
	TechPriest
)

type User struct {
	gorm.Model
	Username          string `gorm:"uniqueIndex"`
	Email             string `gorm:"uniqueIndex"`
	EncryptedPassword string

	//Perms Permissions `gorm:"embedded;embeddedPrefix:perms"`
	Role Role
}
