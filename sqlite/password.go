package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// get hashed password for seed data
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hash))
}
