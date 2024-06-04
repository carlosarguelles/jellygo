package application

import (
	"context"
	"os"

	libdom "github.com/carlosarguelles/jellygo/internal/library/domain"
	movieapp "github.com/carlosarguelles/jellygo/internal/movie/application"
	moviedom "github.com/carlosarguelles/jellygo/internal/movie/domain"
)

type CreateLibraryMetadataUseCase struct {
	movieAPI        movieapp.MovieAPI
	movieRepository movieapp.MovieRepository
}

func (u *CreateLibraryMetadataUseCase) Handle(ctx context.Context, lib *libdom.Library) error {
	dirs, err := os.ReadDir(lib.Path)
	if err != nil {
		return err
	}
	for _, d := range dirs {
		if lib.Type == libdom.TypeMovies {
			movieInfo, err := u.movieAPI.SearchMovieByName(d.Name())
			if err != nil {
				return err
			}
			movie := moviedom.Movie{
				Path:      d.Name(),
				LibraryID: lib.ID,
				Meta: &moviedom.MovieMeta{
					ID:          movieInfo.ID,
					Title:       movieInfo.Title,
					ReleaseDate: movieInfo.ReleaseDate,
				},
			}
			err = u.movieRepository.Create(ctx, &movie)
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

func NewCreateLibraryMetadataUseCase(movieAPI movieapp.MovieAPI, movieRepository movieapp.MovieRepository) *CreateLibraryMetadataUseCase {
	return &CreateLibraryMetadataUseCase{movieAPI, movieRepository}
}
