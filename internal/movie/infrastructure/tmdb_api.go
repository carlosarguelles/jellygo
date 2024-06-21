package infrastructure

import (
	"fmt"

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
			Overview:    apiMovie.Overview,
		}, nil
	}
	return nil, nil
}

func (api *TmdbMovieAPI) GetMovieImages(id int) (*application.MovieImages, error) {
	imagesTiteless, err := api.tmdbAPI.GetMovieImages(id, map[string]string{})
	if err != nil {
		return nil, err
	}
	images, err := api.tmdbAPI.GetMovieImages(id, map[string]string{"language": "en"})
	if err != nil {
		return nil, err
	}
	return &application.MovieImages{
		Logo:   fmt.Sprintf("https://image.tmdb.org/t/p/original%s", images.Logos[0].FilePath),
		Banner: fmt.Sprintf("https://image.tmdb.org/t/p/original%s", images.Backdrops[0].FilePath),
		Bg:     fmt.Sprintf("https://image.tmdb.org/t/p/original%s", imagesTiteless.Backdrops[0].FilePath),
		Poster: fmt.Sprintf("https://image.tmdb.org/t/p/original%s", images.Posters[0].FilePath),
	}, nil
}
