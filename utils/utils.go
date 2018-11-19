package main

import (
	"fmt"
	"github.com/SoftwareFactory-GYN/nehalem/rest_api/user"
)

func seedUsersInDB() {

	users := map[string]string{
		"nsm":  "nsm",
		"gvl":  "gvl",
		"test": "test",
	}

	for username, password := range users {

		u := user.User{
			username,
			password,
		}

		fmt.Printf("Creating user %s with password %s\n", username, password)

		u.Create()

	}

}
