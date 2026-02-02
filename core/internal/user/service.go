package user

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	s := &Service{store: store}
	s.Init()
	return s
}

const DefaultUser = "admin"
const DefaultPassword = "admin6699"
const DefaultUserId = 1

func (s *Service) Init() {
	_, err := s.store.GetByID(DefaultUserId)
	if err == nil {
		log.Debug().Msg("initial user exists")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatal().Err(err).Msg("Failed to get user id")
	}

	err = s.newRaw(1, DefaultUser, DefaultPassword, Omnissiah)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create initial default user")
	}
}

func (s *Service) GetByUsername(username string) (User, error) {
	return s.store.GetByUsername(username)
}

func (s *Service) GetByID(id uint) (User, error) {
	return s.store.GetByID(id)
}

func (s *Service) List(q string, user *User) ([]User, error) {
	if user.Role > Magos {
		return []User{}, fmt.Errorf("insufficient privileges to list users")
	}

	return s.store.List(q)
}

func (s *Service) New(user, password string, role Role, createdBy *User) error {
	finalRole := TechPriest
	if createdBy != nil {
		if createdBy.Role > Magos {
			return fmt.Errorf("role %s does not have permission to create users", createdBy.Role)
		}

		// if requested role is higher privilege (lower number) than the creator
		if role < createdBy.Role {
			log.Warn().
				Str("requested_role", role.String()).
				Str("creator_role", createdBy.Role.String()).
				Msg("Access denied: Cannot create a user with higher privileges than yourself")
			return fmt.Errorf("insufficient privileges for role: %s", role.String())
		}

		// Allow the assignment if it's equal or lower privilege (equal or higher number)
		finalRole = role
	}

	return s.newRaw(0, user, password, finalRole)
}

// registers a new user without role checks assumes all role is verified and trusted
func (s *Service) newRaw(userid uint, user string, password string, finalRole Role) error {
	encrypted, err := EncryptPassword(password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	u := &User{
		Username:          user,
		EncryptedPassword: encrypted,
		Role:              finalRole,
	}
	u.ID = userid

	err = s.store.New(u)
	if err != nil {
		return fmt.Errorf("failed to add to DB: %w", err)
	}
	return err
}

func (s *Service) Edit(user *User, editorUser *User) (err error) {
	if editorUser == nil {
		return fmt.Errorf("editor user cannot be nil")
	}

	if user.ID == editorUser.ID && user.Role != editorUser.Role {
		return fmt.Errorf("cannot change self role")
	}

	if user.Role < editorUser.Role {
		return fmt.Errorf("insufficient privileges to edit role: %s", user.Role.String())
	}

	if user.EncryptedPassword != "" {
		user.EncryptedPassword, err = EncryptPassword(user.EncryptedPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt password: %w", err)
		}

	}

	return s.store.Edit(user)
}

func (s *Service) Delete(id uint, deleteBy *User) error {
	if deleteBy.ID == id {
		return fmt.Errorf("cannot delete self")
	}

	deleted, err := s.store.GetByID(id)
	if err != nil {
		return err
	}

	if deleted.Role < deleteBy.Role || deleteBy.Role > Magos {
		return fmt.Errorf("cannot delete user with higher permission")
	}

	return s.store.Delete(id)
}
