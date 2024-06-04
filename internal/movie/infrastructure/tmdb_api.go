package infrastructure

import (
	"github.com/carlosarguelles/jellygo/internal/movie/application"
	"github.com/ryanbradynd05/go-tmdb"
)

type TmdbMovieAPI struct {
	tmdbAPI *tmdb.TMDb
}

func NewTmdbMovieAPI(tmdbAPI *tmdb.TMDb) *TmdbMovieAPI {
	return &TmdbMovieAPI{tmdbAPI}
}

func (api *TmdbMovieAPI) SearchMovieByName(name string) (*application.MovieInfo, error) {
	movieSearchResults, err := api.tmdbAPI.SearchMovie(name, map[string]string{})
	if err != nil {
		return nil, err
	}
	if len(movieSearchResults.Results) > 0 {
		apiMovie := movieSearchResults.Results[0]
		return &application.MovieInfo{
			ID:          apiMovie.ID,
			Title:       apiMovie.Title,
			ReleaseDate: apiMovie.ReleaseDate,
		}, nil
	}
	return nil, nil
}
