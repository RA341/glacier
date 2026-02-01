package user

import "gorm.io/gorm"

type StoreGorm struct {
	db *gorm.DB
}

func NewStoreGorm(db *gorm.DB) *StoreGorm {
	return &StoreGorm{db: db}
}

func (s StoreGorm) GetByUsername(username string) (User, error) {
	var u User
	err := s.db.Where("username = ?", username).First(&u).Error
	return u, err
}

func (s StoreGorm) GetByID(id uint) (User, error) {
	var u User
	err := s.db.First(&u, id).Error
	return u, err
}

func (s StoreGorm) New(user *User) error {
	return s.db.Create(user).Error
}

func (s StoreGorm) Edit(user *User) error {
	return s.db.Save(user).Error
}

func (s StoreGorm) Delete(id uint) error {
	return s.db.Delete(&User{}, id).Error
}

func (s StoreGorm) List() ([]User, error) {
	var users []User
	err := s.db.Find(&users).Error
	return users, err
}
