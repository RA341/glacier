package auth

import (
	"gorm.io/gorm"
)

type StoreGorm struct {
	db          *gorm.DB
	maxSessions int
}

func NewStoreGorm(db *gorm.DB, maxSessions int) *StoreGorm {
	return &StoreGorm{
		db:          db,
		maxSessions: maxSessions,
	}
}

func (s *StoreGorm) New(token *Session) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(token).Error; err != nil {
			return err
		}

		var idsToKeep []string
		err := tx.Model(&Session{}).
			Select("id").
			Where("user_id = ? AND session_type = ?", token.UserId, token.SessionType).
			// use update instead of created to ensure a refreshed session is not deleted
			Order("updated_at DESC").
			Limit(s.maxSessions).
			Pluck("id", &idsToKeep).
			Error
		if err != nil {
			return err
		}

		// delete everything else for this user and this specific session type
		return tx.
			Where("user_id = ?", token.UserId).
			Where("session_type = ?", token.SessionType).
			Where("id NOT IN ?", idsToKeep).
			Unscoped().Delete(&Session{}).
			Error
	})
}

func (s *StoreGorm) Edit(token *Session) error {
	return s.db.Save(token).Error
}

func (s *StoreGorm) Delete(token *Session) error {
	return s.db.Unscoped().Delete(token).Error
}

func (s *StoreGorm) List() ([]Session, error) {
	var sessions []Session
	err := s.db.Find(&sessions).Error
	return sessions, err
}

func (s *StoreGorm) ListByUser(userId uint) ([]Session, error) {
	var sessions []Session
	err := s.db.Where("user_id = ?", userId).Find(&sessions).Error
	return sessions, err
}

func (s *StoreGorm) GetBySessionId(sessionId uint) (Session, error) {
	var session Session
	err := s.db.First(&session, sessionId).Error
	return session, err
}

func (s *StoreGorm) GetBySessionToken(token string) (Session, error) {
	var session Session
	err := s.db.Preload("User").Where("hashed_session_token = ?", token).First(&session).Error
	return session, err
}
