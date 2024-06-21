package domain

import (
	"fmt"
	"time"
)

type MoviePictures struct {
	Banner string
	Logo   string
}

type MovieMeta struct {
	Pictures    *MoviePictures
	ReleaseDate string
	Title       string
	ID          int
}

func (m *MovieMeta) Year() string {
	t, _ := time.Parse(time.DateOnly, m.ReleaseDate)
	return fmt.Sprint(t.Year())
}

type Movie struct {
	Meta      *MovieMeta
	Path      string
	LibraryID int
	ID        int
}
