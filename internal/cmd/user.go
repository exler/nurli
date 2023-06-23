package cmd

import (
	"fmt"

	"github.com/exler/nurli/internal/core"
	"github.com/exler/nurli/internal/database"
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

			hashed_password, err := core.HashPassword(password)
			if err != nil {
				return err
			}

			db.Create(&database.User{
				Username: username,
				Password: hashed_password,
			})

			fmt.Println("User created successfully")

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

			var users []database.User
			db.Select("id, username").Find(&users)

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

			var user database.User
			db.Where("username = ?", username).First(&user)
			if user.ID == 0 {
				return fmt.Errorf("user not found")
			}

			db.Delete(&user)

			fmt.Println("User removed successfully")

			return nil
		},
	}
)
