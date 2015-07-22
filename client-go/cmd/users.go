package cmd

import (
	"fmt"

	"github.com/brendangibat/deis/client-go/controller/client"
	"github.com/brendangibat/deis/client-go/controller/models/users"
)

// UsersList lists users registered with the controller.
func UsersList() error {
	c, err := client.New()

	if err != nil {
		return err
	}

	users, err := users.List(c)

	if err != nil {
		return err
	}

	fmt.Println("=== Users")

	for _, user := range users {
		fmt.Println(user.Username)
	}
	return nil
}
