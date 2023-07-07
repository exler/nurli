package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/exler/nurli/internal/database"
	"github.com/urfave/cli/v2"
)

var (
	importCmd = &cli.Command{
		Name:  "import",
		Usage: "Import bookmarks from a file",
		Subcommands: []*cli.Command{
			raindropImportCmd,
		},
	}

	raindropImportCmd = &cli.Command{
		Name:  "raindrop",
		Usage: "Import bookmarks from Raindrop.io CSV export",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "read",
				Usage: "Set all imported bookmarks as read",
			},
		},
		Action: func(cCtx *cli.Context) error {
			db, err := openDatabase(cCtx)
			if err != nil {
				return err
			}

			setAsRead := cCtx.Bool("read")

			// Read the file
			filepath := cCtx.Args().First()
			if filepath == "" {
				return fmt.Errorf("missing filepath")
			} else if !strings.HasSuffix(filepath, ".csv") {
				return fmt.Errorf("only CSV files are supported")
			}

			// Open the CSV file
			file, err := os.Open(filepath) // #nosec
			if err != nil {
				return err
			}
			defer file.Close()

			createdRecords := 0

			// Parse the CSV file
			// Format: title, note, excerpt, url, tags, created, cover, highlights
			reader := csv.NewReader(file)
			// Skip the headers
			if _, err = reader.Read(); err != nil {
				logger.Error().Err(err).Msg("failed to read first line")
				return err
			}

			for {
				// Read each record from csv
				record, err := reader.Read()
				if err != nil {
					if err == io.EOF {
						break
					}
					logger.Error().Err(err).Msg("failed to read record")
					return err
				}

				// Parse the created date
				// Format: 2023-06-20T07:20:45.310Z
				createdAt, err := time.Parse("2006-01-02T15:04:05.000Z", record[5])
				if err != nil {
					logger.Error().Err(err).Msg("failed to parse created date")
				}

				// Get or create the bookmark
				bookmark := database.Bookmark{
					URL:         record[3],
					Title:       record[0],
					Description: record[2],
					CreatedAt:   createdAt,
					Read:        setAsRead,
				}
				db.FirstOrCreate(&bookmark, database.Bookmark{
					URL: record[3],
				})

				// Get or create the tags
				if tags := record[4]; tags != "" {
					tags := strings.Split(record[4], ",")
					for _, tagName := range tags {
						tagName = strings.TrimSpace(tagName)
						tag := database.Tag{
							Name: tagName,
						}
						db.FirstOrCreate(&tag, database.Tag{
							Name: tagName,
						})
						bookmark.Tags = append(bookmark.Tags, tag)
					}
				}

				// Save the bookmark
				db.Save(&bookmark)

				createdRecords++
			}

			logger.Info().Msgf("Imported %d records", createdRecords)

			return nil
		},
	}
)
