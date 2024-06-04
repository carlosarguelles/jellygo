package domain

type MovieMeta struct {
	ReleaseDate string
	Title       string
	ID          int
}

type Movie struct {
	Meta      *MovieMeta
	Path      string
	LibraryID int
	ID        int
}
