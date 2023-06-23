package cmd

import (
	"fmt"
	"strconv"

	"github.com/exler/nurli/internal/core"
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

			err = db.CreateBookmark(cCtx.Context, core.Bookmark{
				URL:     url,
				Title:   title,
				OwnerID: userID,
			})
			if err != nil {
				return err
			}

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

			bookmarks, err := db.ListBookmarks(cCtx.Context)
			if err != nil {
				return err
			}

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

			err = db.ToggleBookmarkRead(cCtx.Context, bookmarkID)
			if err != nil {
				return err
			}

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

			err = db.ToggleBookmarkFavorite(cCtx.Context, bookmarkID)
			if err != nil {
				return err
			}

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

			err = db.DeleteBookmark(cCtx.Context, bookmarkID)
			if err != nil {
				return err
			}

			return nil
		},
	}
)
