package users

import "golang.org/x/crypto/bcrypt"

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
		UserStore: &userValidator{UserStore: usg},
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
// passing it on to the next UserStore in the interface chain.
type userValidator struct {
	UserStore
}
