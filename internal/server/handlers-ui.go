package server

import (
	"net/http"
	"strings"

	"github.com/exler/nurli/internal/core"
	"github.com/exler/nurli/internal/database"
)

func (sh *ServerHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequest(r)

	// Get filters from query parameters
	queryParams := r.URL.Query()
	readFilter := queryParams.Get("read")
	favoriteFilter := queryParams.Get("favorite")
	noTagsFilter := queryParams.Get("no-tags")
	tagFilter := queryParams.Get("tag")

	var bookmarks []database.Bookmark
	// Find all bookmarks for the current user and load tags for each bookmark
	// taking into account the filters
	if readFilter != "" {
		sh.DB.Preload("Tags").Where("owner_id = ? AND read = ?", user.ID, readFilter).Find(&bookmarks)
	} else if favoriteFilter != "" {
		sh.DB.Preload("Tags").Where("owner_id = ? AND favorite = ?", user.ID, favoriteFilter).Find(&bookmarks)
	} else if noTagsFilter != "" {
		sh.DB.Preload("Tags").Where("owner_id = ? AND (SELECT COUNT(*) FROM bookmark_tags WHERE bookmark_id = bookmarks.id) = 0", user.ID).Find(&bookmarks)
	} else if tagFilter != "" {
		var tag database.Tag
		sh.DB.Where("name = ? AND owner_id = ?", tagFilter, user.ID).First(&tag)
		sh.DB.Preload("Tags").Where("owner_id = ? AND ? IN (SELECT tag_id FROM bookmark_tags WHERE bookmark_id = bookmarks.id)", user.ID, tag.ID).Find(&bookmarks)
	} else {
		sh.DB.Preload("Tags").Where("owner_id = ?", user.ID).Find(&bookmarks)
	}

	// Find all tags for the current user
	var tags []database.Tag
	sh.DB.Where("owner_id = ?", user.ID).Find(&tags)

	sh.renderTemplate(w, "index", map[string]interface{}{
		"Bookmarks": bookmarks,
		"Tags":      tags,
	})
}

func (sh *ServerHandler) AddBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequest(r)

	if r.Method == "POST" {
		r.ParseForm()
		url := r.FormValue("url")
		tags := r.Form["tags[]"]
		tagObjects := []database.Tag{}

		for _, tag := range tags {
			var tagObj database.Tag
			if strings.HasPrefix(tag, "NEW:") {
				tagObj = database.Tag{
					Name:    strings.TrimPrefix(tag, "NEW:"),
					OwnerID: user.ID,
				}
				sh.DB.Create(&tagObj)
			} else {
				sh.DB.Where("id = ? AND owner_id = ?", tag, user.ID).First(&tagObj)
			}
			tagObjects = append(tagObjects, tagObj)
		}

		page_html, err := core.GetPageHTML(url)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		title := core.GetTitleFromHTML(page_html)
		description := core.TrimString(core.GetDescriptionFromHTML(page_html), core.DESCRIPTION_TRIM_LENGTH)

		bookmark := database.Bookmark{
			URL:         url,
			Title:       title,
			Description: description,
			OwnerID:     user.ID,
			Tags:        tagObjects,
		}
		sh.DB.Create(&bookmark)

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		var tags []database.Tag
		sh.DB.Model(&user).Association("Tags").Find(&tags)
		sh.renderTemplate(w, "add-bookmark", map[string]interface{}{
			"Tags": tags,
		})
	}
}
