package cmd

import (
	"fmt"
	"strconv"

	"github.com/exler/nurli/internal/database"
	"github.com/urfave/cli/v2"
)

var (
	bookmarkCmd = &cli.Command{
		Name:  "bookmark",
		Usage: "Manage bookmarks",
		Subcommands: []*cli.Command{
			bookmarkAddCmd,
			bookmarkListCmd,
			bookmarkMarkAsReadCmd,
			bookmarkMarkAsFavoriteCmd,
			bookmarkRemoveCmd,
		},
	}

	bookmarkAddCmd = &cli.Command{
		Name:  "add",
		Usage: "Add a new bookmark",
		Action: func(cCtx *cli.Context) error {
			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}
			fmt.Print("User ID: ")
			userID, _ := strconv.Atoi(GetUserInput())
			fmt.Print("URL: ")
			url := GetUserInput()
			fmt.Print("Title: ")
			title := GetUserInput()

			db.Create(&database.Bookmark{
				URL:     url,
				Title:   title,
				OwnerID: uint(userID),
			})

			return nil
		},
	}

	bookmarkListCmd = &cli.Command{
		Name:  "list",
		Usage: "List all bookmarks",
		Action: func(cCtx *cli.Context) error {
			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}

			var bookmarks []database.Bookmark
			db.Select("id", "title", "read", "favorite").Find(&bookmarks)

			if len(bookmarks) == 0 {
				fmt.Println("No bookmarks found.")
				return nil
			}

			PrintListHeader("ID", "Title", "Read", "Favorite")
			for _, bookmark := range bookmarks {
				fmt.Printf("%d\t%s\t%t\t%t\n", bookmark.ID, bookmark.Title, bookmark.Read, bookmark.Favorite)
			}

			return nil
		},
	}

	bookmarkMarkAsReadCmd = &cli.Command{
		Name:  "mark-read",
		Usage: "Mark a bookmark as read (or unread)",
		Action: func(cCtx *cli.Context) error {
			bookmarkID, _ := strconv.Atoi(cCtx.Args().Get(0))

			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}

			tx := db.WithContext(cCtx.Context)
			var bookmark database.Bookmark
			tx.First(&bookmark, bookmarkID)
			tx.Model(&database.Bookmark{}).Where("id = ?", bookmarkID).Update("read", !bookmark.Read)

			return nil
		},
	}

	bookmarkMarkAsFavoriteCmd = &cli.Command{
		Name:  "mark-favorite",
		Usage: "Mark a bookmark as favorite (or not favorite)",
		Action: func(cCtx *cli.Context) error {
			bookmarkID, _ := strconv.Atoi(cCtx.Args().Get(0))

			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}

			tx := db.WithContext(cCtx.Context)
			var bookmark database.Bookmark
			tx.First(&bookmark, bookmarkID)
			tx.Model(&database.Bookmark{}).Where("id = ?", bookmarkID).Update("favorite", !bookmark.Favorite)

			return nil
		},
	}

	bookmarkRemoveCmd = &cli.Command{
		Name:  "remove",
		Usage: "Remove a bookmark",
		Action: func(cCtx *cli.Context) error {
			bookmarkID, _ := strconv.Atoi(cCtx.Args().Get(0))

			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}

			db.Delete(&database.Bookmark{}, bookmarkID)

			return nil
		},
	}
)
