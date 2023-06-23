package cmd

import (
	"fmt"
	"strconv"

	"github.com/exler/nurli/internal/core"
	"github.com/urfave/cli/v2"
)

var (
	tagCmd = &cli.Command{
		Name:  "tag",
		Usage: "Manage tags",
		Subcommands: []*cli.Command{
			tagAddCmd,
			tagListCmd,
			tagRemoveCmd,
		},
	}

	tagAddCmd = &cli.Command{
		Name:  "add",
		Usage: "Add a new tag",
		Action: func(cCtx *cli.Context) error {
			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}
			fmt.Print("User ID: ")
			userID, _ := strconv.Atoi(GetUserInput())
			fmt.Print("Name: ")
			name := GetUserInput()

			err = db.CreateTag(cCtx.Context, core.Tag{
				Name:    name,
				OwnerID: userID,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	tagListCmd = &cli.Command{
		Name:  "list",
		Usage: "List all tags",
		Action: func(cCtx *cli.Context) error {
			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}

			tags, err := db.ListTags(cCtx.Context)
			if err != nil {
				return err
			}

			if len(tags) == 0 {
				fmt.Println("No tags found")
				return nil
			}

			PrintListHeader("ID", "Owner ID", "Name")
			for _, tag := range tags {
				fmt.Printf("%d\t%d\t%s\n", tag.ID, tag.OwnerID, tag.Name)
			}

			return nil
		},
	}

	tagRemoveCmd = &cli.Command{
		Name:  "remove",
		Usage: "Remove a tag",
		Action: func(cCtx *cli.Context) error {
			tagID, _ := strconv.Atoi(cCtx.Args().First())

			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}

			err = db.DeleteTag(cCtx.Context, tagID)
			if err != nil {
				return err
			}

			return nil
		},
	}
)
