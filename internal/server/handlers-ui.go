package server

import (
	"net/http"

	"github.com/exler/nurli/internal/core"
	"github.com/exler/nurli/internal/database"
)

func (sh *ServerHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequest(r)

	// Find all bookmarks for the current user
	var bookmarks []database.Bookmark
	sh.DB.Model(&user).Association("Bookmarks").Find(&bookmarks)

	sh.renderTemplate(w, "index", bookmarks)
}

func (sh *ServerHandler) AddBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequest(r)

	if r.Method == "POST" {
		url := r.FormValue("url")
		// tags := r.FormValue("tags")
		page_html, err := core.GetPageHTML(url)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		title := core.GetTitleFromHTML(page_html)
		description := core.GetDescriptionFromHTML(page_html)

		bookmark := database.Bookmark{
			URL:         url,
			Title:       title,
			Description: description,
			OwnerID:     user.ID,
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
