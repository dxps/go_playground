package main

import (
	"devisions.org/goallery/users"
	"devisions.org/goallery/utils/rand"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 54321
	user     = "goallery"
	password = "goallery"
	dbname   = "goallery"
)

func main() {

	dbConnInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	userSvc, err := users.NewUserService(dbConnInfo)
	if err != nil {
		panic(err)
	}
	defer userSvc.Close()

	// Dangerous!!! Watch out!!!
	//_ = userSvc.DestructiveReset()

	user := users.User{
		Name:  "Joe Black",
		Email: "joe@black.com",
	}

	// Addition test
	if err := userSvc.Create(&user); err != nil {
		panic(err)
	}

	// GetByEmail test
	foundUser, err := userSvc.GetByEmail("joe@black.com")
	if err != nil {
		panic(err)
	}
	fmt.Printf(">>> Found user: %+v\n", foundUser)

	// Update test
	user.Name = "Joe Black Updated"
	if err := userSvc.Update(&user); err != nil {
		panic(err)
	}

	// Delete test
	if err := userSvc.Delete(user.ID); err != nil {
		panic(">>> Error deleting user: " + err.Error())
	}

	// Testing of remember token value generator.
	remToken, _ := rand.RememberToken()
	fmt.Printf(">>> Remember Token value: %+v\n", remToken)

}
