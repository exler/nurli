package core

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func GetDomainFromURL(urlToParse string) string {
	parsedUrl, err := url.Parse(urlToParse)
	if err != nil {
		return ""
	}

	return strings.TrimPrefix(parsedUrl.Hostname(), "www.")
}

func GetPageHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url) // #nosec
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func GetTitleFromHTML(doc *html.Node) string {
	var title string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			title = n.FirstChild.Data
		}

		if title == "" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(doc)

	return title
}

func GetDescriptionFromHTML(doc *html.Node) string {
	var description string

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, a := range n.Attr {
				if a.Key == "name" && a.Val == "description" {
					for _, a := range n.Attr {
						if a.Key == "content" {
							description = a.Val
						}
					}
				}
			}
		}

		if description == "" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	f(doc)

	return description
}

func TrimString(s string, max_length int) string {
	if len(s) > max_length {
		return s[:max_length] + "..."
	}

	return s
}

func StringIn(s string, slice []string) bool {
	for _, val := range slice {
		if val == s {
			return true
		}
	}

	return false
}

func GetPageFromQueryParams(queryParams url.Values) int {
	pageParam := queryParams.Get("page")
	var page int
	if pageParam == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(pageParam)
		if page <= 0 {
			page = 1
		}
	}

	return page
}

func GetPageSizeFromQueryParams(queryParams url.Values) int {
	pageSizeParam := queryParams.Get("page-size")
	var pageSize int
	if pageSizeParam == "" {
		pageSize = DEFAULT_PAGE_SIZE
	} else {
		pageSize, _ = strconv.Atoi(pageSizeParam)
		switch {
		case pageSize > MAX_PAGE_SIZE:
			pageSize = MAX_PAGE_SIZE
		case pageSize <= 0:
			pageSize = DEFAULT_PAGE_SIZE
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
