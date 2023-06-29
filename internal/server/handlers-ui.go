package server

import (
	"net/http"
	"strings"

	"github.com/exler/nurli/internal/core"
	"github.com/exler/nurli/internal/database"
	"github.com/go-chi/chi"
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

	sh.renderTemplate(w, "bookmark/bookmark_list", map[string]interface{}{
		"Bookmarks": bookmarks,
		"Tags":      tags,
	})
}

func (sh *ServerHandler) AddBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequest(r)

	if r.Method == "POST" {
		r.ParseForm()
		url := r.FormValue("url")
		read := r.FormValue("read") == "on"
		favorite := r.FormValue("favorite") == "on"
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
			Read:        read,
			Favorite:    favorite,
			OwnerID:     user.ID,
			Tags:        tagObjects,
		}
		sh.DB.Create(&bookmark)

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		var tags []database.Tag
		sh.DB.Model(&user).Association("Tags").Find(&tags)
		sh.renderTemplate(w, "bookmark/bookmark_change_form", map[string]interface{}{
			"Tags": tags,
		})
	}
}

func (sh *ServerHandler) EditBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequest(r)

	if r.Method == "POST" {
		r.ParseForm()
		url := r.FormValue("url")
		read := r.FormValue("read") == "on"
		favorite := r.FormValue("favorite") == "on"
		tags := r.Form["tags[]"]

		// Get the bookmark from the database
		var bookmark database.Bookmark
		sh.DB.Preload("Tags").Where("id = ? AND owner_id = ?", chi.URLParam(r, "id"), user.ID).First(&bookmark)

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

		title := bookmark.Title
		description := bookmark.Description

		if url != bookmark.URL {
			title = core.GetTitleFromHTML(page_html)
			description = core.TrimString(core.GetDescriptionFromHTML(page_html), core.DESCRIPTION_TRIM_LENGTH)
		}

		// Update the bookmark
		sh.DB.Model(&bookmark).Updates(database.Bookmark{
			URL:         url,
			Title:       title,
			Description: description,
			Read:        read,
			Favorite:    favorite,
		})

		// Update the tags
		sh.DB.Model(&bookmark).Association("Tags").Replace(tagObjects)

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		var bookmark database.Bookmark
		sh.DB.Where("id = ? AND owner_id = ?", chi.URLParam(r, "id"), user.ID).Preload("Tags").First(&bookmark)

		var tags []database.Tag
		sh.DB.Model(&user).Association("Tags").Find(&tags)

		var initialTags []uint
		for _, tag := range bookmark.Tags {
			initialTags = append(initialTags, tag.ID)
		}

		sh.renderTemplate(w, "bookmark/bookmark_change_form", map[string]interface{}{
			"Bookmark":    bookmark,
			"Tags":        tags,
			"InitialTags": initialTags,
		})
	}
}

func (sh *ServerHandler) DeleteBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequest(r)

	var bookmark database.Bookmark
	sh.DB.Where("id = ? AND owner_id = ?", chi.URLParam(r, "id"), user.ID).First(&bookmark)

	if r.Method == "POST" {
		sh.DB.Delete(&bookmark)

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		sh.renderTemplate(w, "bookmark/bookmark_confirm_delete", map[string]interface{}{
			"Bookmark": bookmark,
		})
	}
}

func (sh *ServerHandler) SettingsHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromRequest(r)

	if r.Method == "POST" {
		var message string

		oldPassword := r.FormValue("old_password")
		newPassword := r.FormValue("new_password")
		if oldPassword != "" && newPassword != "" {
			if !core.CheckPasswordHash(oldPassword, user.Password) {
				message = "Old password is incorrect!"
			} else {
				hashedPassword, err := core.HashPassword(r.FormValue("password"))
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				sh.DB.Model(&user).Update("Password", hashedPassword)
				message = "Password updated!"
			}
		}

		sh.renderTemplate(w, "settings", map[string]interface{}{
			"Message": message,
			"User":    user,
		})
	} else {
		sh.renderTemplate(w, "settings", map[string]interface{}{
			"User": user,
		})
	}
}
