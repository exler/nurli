package server

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/exler/nurli/internal"
)

func GetPageFromQueryParams(queryParams url.Values) int {
	pageParam := queryParams.Get("page")
	var page int
	if pageParam == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(pageParam)
	}

	return page
}

func GetPageSizeFromQueryParams(queryParams url.Values) int {
	pageSizeParam := queryParams.Get("page-size")
	var pageSize int
	if pageSizeParam == "" {
		pageSize = internal.DEFAULT_PAGE_SIZE
	} else {
		pageSize, _ = strconv.Atoi(pageSizeParam)
		switch {
		case pageSize > internal.MAX_PAGE_SIZE:
			pageSize = internal.MAX_PAGE_SIZE
		case pageSize <= 0:
			pageSize = internal.DEFAULT_PAGE_SIZE
		}
	}

	return pageSize
}

func UpdateSingleParamInURL(r *http.Request, key, value string) string {
	queryParams := r.URL.Query()
	queryParams.Set(key, value)
	r.URL.RawQuery = queryParams.Encode()

	return r.URL.String()
}
