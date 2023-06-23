package cmd

import (
	"fmt"

	"github.com/exler/nurli/internal/core"
	"github.com/urfave/cli/v2"
)

var (
	userCmd = &cli.Command{
		Name:  "user",
		Usage: "Manage users",
		Subcommands: []*cli.Command{
			userAddCmd,
			userListCmd,
			userRemoveCmd,
		},
	}

	userAddCmd = &cli.Command{
		Name:  "add",
		Usage: "Add a new user",
		Action: func(cCtx *cli.Context) error {
			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}
			fmt.Print("Username: ")
			username := GetUserInput()
			fmt.Print("Password: ")
			password := GetUserSecureInput()

			err = db.CreateUser(cCtx.Context, core.User{
				Username: username,
				Password: password,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	userListCmd = &cli.Command{
		Name:  "list",
		Usage: "List all users",
		Action: func(cCtx *cli.Context) error {
			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}

			users, err := db.ListUsers(cCtx.Context)
			if err != nil {
				return err
			}

			if len(users) == 0 {
				fmt.Println("No users found")
				return nil
			}

			PrintListHeader("ID", "Username")
			for _, user := range users {
				fmt.Printf("%d\t%s\n", user.ID, user.Username)
			}

			return nil
		},
	}

	userRemoveCmd = &cli.Command{
		Name:  "remove",
		Usage: "Remove a user",
		Action: func(cCtx *cli.Context) error {
			username := cCtx.Args().Get(0)
			if username == "" {
				return fmt.Errorf("username is required")
			}

			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}

			err = db.DeleteUserByUsername(cCtx.Context, username)
			if err != nil {
				return err
			}

			return nil
		},
	}
)
