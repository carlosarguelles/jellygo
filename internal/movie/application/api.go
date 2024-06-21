package application

type MovieInfo struct {
	Title       string
	ReleaseDate string
	Overview    string
	ID          int
}

type MovieImages struct {
	Logo   string
	Banner string
	Bg     string
	Poster string
}

type MovieAPI interface {
	SearchMovieByName(string) (*MovieInfo, error)
	GetMovieImages(int) (*MovieImages, error)
}
