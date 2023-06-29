package server

import (
	"net/http"

	"github.com/exler/nurli/internal/core"
	"github.com/exler/nurli/internal/database"
	"github.com/go-chi/chi"
)

func (sh *ServerHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Get filters from query parameters
	queryParams := r.URL.Query()
	readFilter := queryParams.Get("read")
	favoriteFilter := queryParams.Get("favorite")
	noTagsFilter := queryParams.Get("no-tags")
	tagFilter := queryParams.Get("tag")

	var bookmarks []database.Bookmark
	// Find all bookmarks and load tags for each bookmark taking into account the filters
	if readFilter != "" {
		sh.DB.Preload("Tags").Where("read = ?", readFilter).Find(&bookmarks)
	} else if favoriteFilter != "" {
		sh.DB.Preload("Tags").Where("favorite = ?", favoriteFilter).Find(&bookmarks)
	} else if noTagsFilter != "" {
		sh.DB.Preload("Tags").Where("(SELECT COUNT(*) FROM bookmark_tags WHERE bookmark_id = bookmarks.id) = 0").Find(&bookmarks)
	} else if tagFilter != "" {
		var tag database.Tag
		sh.DB.Where("name = ?", tagFilter).First(&tag)
		sh.DB.Preload("Tags").Where("? IN (SELECT tag_id FROM bookmark_tags WHERE bookmark_id = bookmarks.id)", tag.ID).Find(&bookmarks)
	} else {
		sh.DB.Preload("Tags").Find(&bookmarks)
	}

	var tags []database.Tag
	sh.DB.Find(&tags)

	sh.renderTemplate(w, "bookmark/bookmark_list", map[string]interface{}{
		"Bookmarks": bookmarks,
		"Tags":      tags,
	})
}

func (sh *ServerHandler) AddBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		url := r.FormValue("url")
		read := r.FormValue("read") == "on"
		favorite := r.FormValue("favorite") == "on"
		tags := r.Form["tags[]"]
		tagObjects := []database.Tag{}

		for _, tag := range tags {
			tagObj := database.Tag{
				Name: tag,
			}
			sh.DB.Where("name = ?", tag).FirstOrCreate(&tagObj)
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
			Tags:        tagObjects,
		}
		sh.DB.Create(&bookmark)

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		var tags []database.Tag
		sh.DB.Find(&tags)
		sh.renderTemplate(w, "bookmark/bookmark_change_form", map[string]interface{}{
			"Tags": tags,
		})
	}
}

func (sh *ServerHandler) EditBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		url := r.FormValue("url")
		read := r.FormValue("read") == "on"
		favorite := r.FormValue("favorite") == "on"
		tags := r.Form["tags[]"]

		// Get the bookmark from the database
		var bookmark database.Bookmark
		sh.DB.Preload("Tags").Where("id = ?", chi.URLParam(r, "id")).First(&bookmark)

		tagObjects := []database.Tag{}
		for _, tag := range tags {
			tagObj := database.Tag{
				Name: tag,
			}
			sh.DB.Where("name = ?", tag).FirstOrCreate(&tagObj)
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

		// Update the bookmark as map because GORM only updates non-zero fields
		sh.DB.Model(&bookmark).Updates(map[string]interface{}{
			"URL":         url,
			"Title":       title,
			"Description": description,
			"Read":        read,
			"Favorite":    favorite,
		})

		// Update the tags
		sh.DB.Model(&bookmark).Association("Tags").Replace(tagObjects)

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		var bookmark database.Bookmark
		sh.DB.Where("id = ?", chi.URLParam(r, "id")).Preload("Tags").First(&bookmark)

		var tags []database.Tag
		sh.DB.Find(&tags)

		var initialTags []string
		for _, tag := range bookmark.Tags {
			initialTags = append(initialTags, tag.Name)
		}

		sh.renderTemplate(w, "bookmark/bookmark_change_form", map[string]interface{}{
			"Bookmark":    bookmark,
			"Tags":        tags,
			"InitialTags": initialTags,
		})
	}
}

func (sh *ServerHandler) DeleteBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	var bookmark database.Bookmark
	sh.DB.Where("id = ?", chi.URLParam(r, "id")).First(&bookmark)

	if r.Method == "POST" {
		sh.DB.Delete(&bookmark)

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		sh.renderTemplate(w, "bookmark/bookmark_confirm_delete", map[string]interface{}{
			"Bookmark": bookmark,
		})
	}
}
