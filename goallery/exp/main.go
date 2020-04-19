package main

import (
	"devisions.org/goallery/commons/rand"
	"fmt"

	"devisions.org/goallery/models"
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

	ur, err := models.NewUserRepo(dbConnInfo)
	if err != nil {
		panic(err)
	}
	defer ur.Close()
	_ = ur.DestructiveReset()

	user := models.User{
		Name:  "Joe Black",
		Email: "joe@black.com",
	}

	// Addition test
	if err := ur.Add(&user); err != nil {
		panic(err)
	}

	// GetByEmail test
	foundUser, err := ur.GetByEmail("joe@black.com")
	if err != nil {
		panic(err)
	}
	fmt.Printf(">>> Found user: %+v\n", foundUser)

	// Update test
	user.Name = "Joe Black Updated"
	if err := ur.Update(&user); err != nil {
		panic(err)
	}

	// Delete test
	if err := ur.Delete(user.ID); err != nil {
		panic(">>> Error deleting user: " + err.Error())
	}

	// Testing of remember token value generator.
	remToken, _ := rand.RememberToken()
	fmt.Printf(">>> Remember Token value: %+v\n", remToken)

}
