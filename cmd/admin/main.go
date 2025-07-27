package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/r-scheele/zero/ent/user"
	"github.com/r-scheele/zero/pkg/log"
	"github.com/r-scheele/zero/pkg/services"
)

// main creates a new admin user with the phone number passed in via the flag.
func main() {
	// Start a new container.
	c := services.NewContainer()
	defer func() {
		// Gracefully shutdown all services.
		if err := c.Shutdown(); err != nil {
			log.Default().Error("shutdown failed", "error", err)
		}
	}()

	var phone string
	flag.StringVar(&phone, "phone", "", "phone number for the admin user (E.164 format, e.g., +1234567890)")
	flag.Parse()

	if len(phone) == 0 {
		invalid("phone number is required")
	}

	// Generate a password.
	pw, err := c.Auth.RandomToken(10)
	if err != nil {
		invalid("failed to generate a random password")
	}

	// Create the admin user.
	err = c.ORM.User.
		Create().
		SetPhoneNumber(phone).
		SetName("Admin").
		SetAdmin(true).
		SetVerified(true).
		SetPassword(pw).
		SetRegistrationMethod(user.RegistrationMethodWeb).
		Exec(context.Background())

	if err != nil {
		invalid(err.Error())
	}

	fmt.Println("")
	fmt.Println("-- ADMIN USER CREATED --")
	fmt.Printf("Phone: %s\n", phone)
	fmt.Printf("Password: %s\n", pw)
	fmt.Println("----")
	fmt.Println("")
}

func invalid(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
	os.Exit(1)
}
