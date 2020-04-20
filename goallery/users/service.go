package users

import (
	"devisions.org/goallery/utils/hash"
	"devisions.org/goallery/utils/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// UserService contains the features of working with User model.
type UserService interface {

	// embedded interface
	UserStore

	// Authenticate validates the provided email and password.
	// If correct, the user is returned.
	// Otherwise, either ErrNotFound, ErrInvalidPassword, or any
	// other error if something goes wrong.
	Authenticate(email, password string) (*User, error)
}

// Implementation of UserService.
type userService struct {
	UserStore
}

// NewUserService creates a new instance of a `UserStore` implementation.
func NewUserService(connectionInfo string) (UserService, error) {
	usg, err := newUserStoreGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	return &userService{
		UserStore: &userValidator{
			UserStore: usg,
			hmac:      hash.NewHMAC(hmacSecretKey),
		},
	}, nil
}

// Authenticate is used for authenticating the provided user credentials.
func (us *userService) Authenticate(email, password string) (*User, error) {

	foundUser, err := us.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+userPwdPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPwd
	default:
		return nil, err
	}
}

// ---------------------------------
//  Validation Layer
// ---------------------------------

// This is a layer that validates and normalizes data before
// passing it on to the next layer in the interface chain, the UserStore implementation.
// It is (and remains) hidden within UserService.
type userValidator struct {
	UserStore
	hmac hash.HMAC
}

const hmacSecretKey = "secret-hmac-key"

// GetByRemember will hash the remember token and then call
// the UserStore's method with the same name.
func (uv *userValidator) GetByRemember(token string) (*User, error) {

	user := User{Remember: token} // working with a user just for rememberHash
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	fmt.Printf(">>> userValidator > GetByRemember > Normalized token '%v' as hash '%v' before searching.\n",
		token, user.RememberHash)
	return uv.UserStore.GetByRemember(user.RememberHash)
}

// Create method inserts a new user into the store.
// This intermediate method of userValidator does any data validation and normalization.
func (uv *userValidator) Create(user *User) error {

	if err := runUserValFns(user,
		uv.bcryptPassword,
		uv.setRememberIfUnset,
		uv.hmacRemember); err != nil {
		return err
	}
	return uv.UserStore.Create(user)
}

// Update will updates the existing record of the provided user.
// This intermediate method of userValidator does any data validation and normalization.
func (uv *userValidator) Update(user *User) error {

	if err := runUserValFns(user,
		uv.bcryptPassword,
		uv.hmacRemember); err != nil {
		return err
	}
	return uv.UserStore.Update(user)
}

// Delete will delete the user record with the provided ID.
// It may return ErrInvalidID if provided ID is 0, just to prevent an accidentally deletion of all users.
// This intermediate method of userValidator does any data validation and normalization.
func (uv *userValidator) Delete(id uint) error {

	var user User
	user.ID = id
	if err := runUserValFns(&user, uv.genUserValFnIdGreaterThan(0)); err != nil {
		return err
	}
	return uv.UserStore.Delete(id)
}

// bcryptPassword will hash a user's password with an app-wide pepper
// and brypt's salt.
func (uv *userValidator) bcryptPassword(user *User) error {

	if user.Password == "" {
		// Skipping if password hasn't been changed.
		return nil
	}
	pwdBytes := []byte(user.Password + userPwdPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwdBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

type userValFn func(*User) error

func runUserValFns(user *User, fns ...userValFn) error {

	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {

	if user.Remember != "" {
		user.RememberHash = uv.hmac.Hash(user.Remember)
	}
	return nil // needed just to satisfy the signature of userValFn
}

// Generate a userValFn closure that checks if the user's ID is greater than provided n.
func (uv *userValidator) genUserValFnIdGreaterThan(n uint) userValFn {

	return userValFn(
		func(user *User) error {
			if user.ID <= n {
				return ErrInvalidID
			}
			return nil
		})
}
