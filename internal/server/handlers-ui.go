package server

import (
	"net/http"

	"github.com/exler/nurli/internal/database"
)

func (sh *ServerHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	var bookmarks []database.Bookmark
	sh.DB.Find(&bookmarks)

	sh.renderTemplate(w, "index", bookmarks)
}
