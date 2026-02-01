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

		// Subquery finds the IDs of the newest N sessions to keep them
		// The outer query deletes everything else for this user
		err := tx.
			Where("user_id = ? AND id NOT IN (?)",
				token.UserId,
				tx.Model(&Session{}).
					Select("id").
					Where("user_id = ?", token.UserId).
					Where("session_type = ?", token.SessionType).
					Order("created_at DESC").
					Limit(s.maxSessions),
			).
			Delete(&Session{}).
			Error

		return err
	})
}

func (s *StoreGorm) Edit(token *Session) error {
	return s.db.Save(token).Error
}

func (s *StoreGorm) Delete(token *Session) error {
	return s.db.Delete(token).Error
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
