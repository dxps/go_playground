package web

import "encoding/gob"

func init() {
	gob.Register(CreatePostForm{})
	gob.Register(CreateThreadForm{})
	gob.Register(CreateCommentForm{})
	gob.Register(UserRegistrationForm{})
	gob.Register(FormErrors{})
}

type FormErrors map[string]string

type CreatePostForm struct {
	Title   string
	Content string
	Errors  FormErrors
}

func (f *CreatePostForm) Validate() bool {

	f.Errors = FormErrors{}
	if f.Title == "" {
		f.Errors["Title"] = "Please provide a title"
	}
	if f.Content == "" {
		f.Errors["Content"] = "Please enter a content"
	}
	return len(f.Errors) == 0
}

type CreateThreadForm struct {
	Title       string
	Description string
	Errors      FormErrors
}

func (f *CreateThreadForm) Validate() bool {

	f.Errors = FormErrors{}
	if f.Title == "" {
		f.Errors["Title"] = "Please provide a title"
	}
	if f.Description == "" {
		f.Errors["Description"] = "Please enter a description"
	}
	return len(f.Errors) == 0
}

type CreateCommentForm struct {
	Content string
	Errors  FormErrors
}

func (f *CreateCommentForm) Validate() bool {

	f.Errors = FormErrors{}
	if f.Content == "" {
		f.Errors["Content"] = "Please enter a content"
	}
	return len(f.Errors) == 0
}

type UserRegistrationForm struct {
	Username      string
	Password      string
	UsernameTaken bool
	Errors        FormErrors
}

func (f *UserRegistrationForm) Validate() bool {

	f.Errors = FormErrors{}
	if f.Username == "" {
		f.Errors["Username"] = "Please enter a username"
	} else if f.UsernameTaken {
		f.Errors["Username"] = "This username is already taken"
	}
	if f.Password == "" {
		f.Errors["Password"] = "Please enter a password"
	} else if len(f.Password) < 8 {
		f.Errors["Password"] = "The password must be at least 8 characters long"
	}
	return len(f.Errors) == 0
}
