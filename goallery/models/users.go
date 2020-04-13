package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource could not be found in the database.
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID is returned when an invalid ID (like 0) is provided to a method like Delete.
	ErrInvalidID = errors.New("models: provided ID is invalid")
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(connectionInfo string) (*UserRepo, error) {

	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserRepo{db: db}, nil
}

// Close method closes the database connection.
func (ur *UserRepo) Close() error {
	return ur.db.Close()
}

// Add method inserts a new user into the repository.
func (ur *UserRepo) Add(user *User) error {
	return ur.db.Create(user).Error
}

// GetByID looks up a user with the provided ID.
// If the user is found, the error will be nil, otherwise an ErrNotFound will be returned.
// In case of any other issue, details are included in the returned error.
func (ur *UserRepo) GetByID(id uint) (*User, error) {

	var user User
	db := ur.db.Where("id = ?", id)
	if err := first(db, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail looks up a user with the given email address and returns that user, plus a nil error.
// If not found, returned user is nil, and error is ErrNotFound.
// If any other error, it will be returned and also returned user is nil.
func (ur *UserRepo) GetByEmail(email string) (*User, error) {

	var user User
	db := ur.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// Update will updates the existing record of the provided user.
func (ur *UserRepo) Update(user *User) error {
	return ur.db.Save(user).Error
}

// Delete will delete the user record with the provided ID.
// It may return ErrInvalidID if provided ID is 0, just to prevent an accidentally deletion of all users.
func (ur *UserRepo) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ur.db.Delete(&user).Error
}

// AutoMigrate attempts to automatically migrate the users table.
func (ur *UserRepo) AutoMigrate() error {

	return ur.db.AutoMigrate(&User{}).Error
}

// DestructiveReset drops the user table and recreates it.
// Needed for development purposes only.
func (ur *UserRepo) DestructiveReset() error {

	if err := ur.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ur.AutoMigrate()
}

// first will query using the provided gorm.DB to get the first item returned
// and place it into dst. If nothing found, ErrNotFound will be returned.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
