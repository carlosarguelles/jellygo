package domain

import moviedom "github.com/carlosarguelles/jellygo/internal/movie/domain"

const (
	TypeMovies  = "movies"
	TypeTVShows = "tvshows"
)

type Library struct {
	Name   string
	Path   string
	Type   string
	Movies []*moviedom.Movie
	ID     int
}

type LibraryEntry struct {
	PathName string
	Name     string
}

func GetLibraryTypes() []string {
	return []string{TypeMovies, TypeTVShows}
}

func NewLibrary(path string, _type string, name string) *Library {
	return &Library{Path: path, Type: _type, Name: name}
}
