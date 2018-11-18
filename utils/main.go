package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "seed_db",
			Aliases: []string{"sdu"},
			Usage:   "seed the database with users",
			Action: func(c *cli.Context) error {
				seedUsersInDB()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
