package users

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resource could not be found in the database.
	ErrNotFound = errors.New("models: resource not found")
	// ErrInvalidID is returned when an invalid ID (like 0) is provided to a method like Delete.
	ErrInvalidID = errors.New("models: provided ID is invalid")
	// ErrInvalidPwd is returned when invalid password is provided at user login.
	ErrInvalidPwd = errors.New("models: provided password is invalid")

	userPwdPepper = "some-secret-random-string"
)

// UserStore is used for interacting with a user store.
// For all `Get_` methods, it either returns the user that is found and a nil error,
// or an error that is either defined by the `models` package (such as `ErrNotFound`)
// or another, more low level error.
// This is a contract (interface) used by outside package components.
type UserStore interface {
	GetByID(int uint) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByRemember(token string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Close is used for closing the connection(s) to the store (database).
	Close()

	// AutoMigrate is a helper method used for database migration
	AutoMigrate() error
	// DestructiveReset is a helper method used only for dev purposes.
	DestructiveReset() error
}

// This is an implementation of `UserStore` interface.
type userStoreGorm struct {
	db *gorm.DB
}

// Internal constructor of a userStoreGorm instance.
func newUserStoreGorm(connectionInfo string) (*userStoreGorm, error) {

	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userStoreGorm{db: db}, nil
}

// Close method closes the database connection.
func (ur *userStoreGorm) Close() {
	_ = ur.db.Close()
}

// Create method inserts a new user into the repository.
func (ur *userStoreGorm) Create(user *User) error {

	return ur.db.Create(user).Error
}

// GetByID looks up a user with the provided ID.
// If the user is found, the error will be nil, otherwise an ErrNotFound will be returned.
// In case of any other issue, details are included in the returned error.
func (ur *userStoreGorm) GetByID(id uint) (*User, error) {

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
func (ur *userStoreGorm) GetByEmail(email string) (*User, error) {

	var user User
	db := ur.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// GetByRememberHash looks up a user with the given remember token.
// If not found, returned user is nil, and error is ErrNotFound.
// If any other error, it will be returned and also returned user is nil.
func (ur *userStoreGorm) GetByRemember(rememberHash string) (*User, error) {

	var user User
	fmt.Printf(">>> userStoreGorm > GetByRemember > Looking for user with rememberHash:'%v'\n", rememberHash)
	err := first(ur.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update will updates the existing record of the provided user.
func (ur *userStoreGorm) Update(user *User) error {

	return ur.db.Save(user).Error
}

// Delete will delete the user record with the provided ID.
// It may return ErrInvalidID if provided ID is 0, just to prevent an accidentally deletion of all users.
func (ur *userStoreGorm) Delete(id uint) error {

	user := User{Model: gorm.Model{ID: id}}
	return ur.db.Delete(&user).Error
}

// AutoMigrate attempts to automatically migrate the users table.
func (ur *userStoreGorm) AutoMigrate() error {

	return ur.db.AutoMigrate(&User{}).Error
}

// DestructiveReset drops the user table and recreates it.
// Needed for development purposes only.
func (ur *userStoreGorm) DestructiveReset() error {

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
