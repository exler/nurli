package server

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/exler/nurli/internal/core"
	"github.com/exler/nurli/internal/database"
)

func (sh *ServerHandler) HealthAPIHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sh *ServerHandler) URLDetailAPIHandler(w http.ResponseWriter, r *http.Request) {
	// Get the encoded URL from the `url` parameter
	encodedUrl := r.URL.Query().Get("url")
	if encodedUrl == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// Decode the URL
	decodedURL, err := url.QueryUnescape(encodedUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove all query parameters from the URL
	parsedURL, err := url.Parse(decodedURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	parsedURL.RawQuery = ""
	finalUrl := parsedURL.String()

	// Check if URL exists in database
	var bookmark database.Bookmark
	sh.DB.Preload("Tags").Where("url = ?", finalUrl).First(&bookmark)

	// If URL does not exist in database, simply return 404
	if bookmark.ID == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// If URL exists in database, return the bookmark as JSON
	if err := WriteAsJSON(w, &bookmark); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (sh *ServerHandler) SaveBookmarkAPIHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		URL      string   `json:"URL"`
		Read     bool     `json:"Read"`
		Favorite bool     `json:"Favorite"`
		Tags     []string `json:"Tags"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if requestData.URL == "" {
		http.Error(w, "URL is missing", http.StatusBadRequest)
		return
	}

	tagObjects := []database.Tag{}

	for _, tag := range requestData.Tags {
		tagObj := database.Tag{
			Name: tag,
		}
		sh.DB.Where("name = ?", tag).FirstOrCreate(&tagObj)
		tagObjects = append(tagObjects, tagObj)
	}

	var bookmark database.Bookmark
	sh.DB.Where("url = ?", requestData.URL).First(&bookmark)

	var responseCode int

	if bookmark.ID == 0 {
		page_html, err := core.GetPageHTML(requestData.URL)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		title := core.GetTitleFromHTML(page_html)
		description := core.TrimString(core.GetDescriptionFromHTML(page_html), core.DESCRIPTION_TRIM_LENGTH)

		bookmark := database.Bookmark{
			URL:         requestData.URL,
			Title:       title,
			Description: description,
			Read:        requestData.Read,
			Favorite:    requestData.Favorite,
			Tags:        tagObjects,
		}
		sh.DB.Create(&bookmark)

		responseCode = http.StatusCreated
	} else {
		bookmark.Read = requestData.Read
		bookmark.Favorite = requestData.Favorite

		if err := sh.DB.Model(&bookmark).Association("Tags").Replace(tagObjects); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sh.DB.Save(&bookmark)

		responseCode = http.StatusOK
	}

	w.WriteHeader(responseCode)
}

func (sh *ServerHandler) DeleteBookmarkAPIHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		URL string `json:"URL"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if requestData.URL == "" {
		http.Error(w, "URL is missing", http.StatusBadRequest)
		return
	}

	// Get the bookmark from the database
	var bookmark database.Bookmark
	sh.DB.Where("url = ?", requestData.URL).First(&bookmark)

	// Delete the bookmark
	sh.DB.Delete(&bookmark)

	// Return 200 OK
	w.WriteHeader(http.StatusOK)
}
