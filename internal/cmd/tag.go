package cmd

import (
	"fmt"
	"strconv"

	"github.com/exler/nurli/internal/database"
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
			fmt.Print("Name: ")
			name := GetUserInput()

			db.Create(&database.Tag{
				Name: name,
			})

			fmt.Println("Tag created successfully")

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

			var tags []database.Tag
			db.Find(&tags)

			if len(tags) == 0 {
				fmt.Println("No tags found")
				return nil
			}

			PrintListHeader("ID", "Name")
			for _, tag := range tags {
				fmt.Printf("%d\t%s\n", tag.ID, tag.Name)
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

			var tag database.Tag
			db.First(&tag, tagID)
			if tag.ID == 0 {
				fmt.Println("Tag not found")
			}

			db.Delete(&tag)

			fmt.Println("Tag removed successfully")

			return nil
		},
	}
)
