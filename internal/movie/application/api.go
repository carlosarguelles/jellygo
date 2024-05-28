package application

type MovieInfo struct {
	Title       string
	ReleaseDate string
}

type MovieAPI interface {
	SearchMovieByName(string) (*MovieInfo, error)
}
