package views

import (
	"fmt"
	"log"

	"github.com/dionvu/spogo/auth"
	"github.com/dionvu/spogo/spotify"
)

type SearchResultView struct {
	items      *spotify.SearchResult
	query      string
	searchType string
}

func NewSearchResultView(searchQuery string, searchType string, s *auth.Session) *SearchResultView {
	srv := SearchResultView{
		query:      searchQuery,
		searchType: searchType,
	}

	if searchQuery == "" {
		log.Fatal("Empty search query")
	}

	fmt.Println(searchType)
	switch searchType {
	case "track":
		results, err := spotify.Search(searchQuery, []string{"track"}, s)
		if err != nil {
			log.Fatal("Error getting results")
		}

		srv.items = results

		fmt.Println(srv.items.Tracks)
	default:
		log.Fatal("Unknown search type passed")
	}

	return &srv
}

func (srv *SearchResultView) i() *spotify.SearchResult {
	return srv.items
}
