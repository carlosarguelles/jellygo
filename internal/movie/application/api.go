package application

type MovieInfo struct {
	Title       string
	ReleaseDate string
	ID          int
}

type MovieAPI interface {
	SearchMovieByName(string) (*MovieInfo, error)
}
