package application

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"sync"

	libdom "github.com/carlosarguelles/jellygo/internal/library/domain"
	movieapp "github.com/carlosarguelles/jellygo/internal/movie/application"
	moviedom "github.com/carlosarguelles/jellygo/internal/movie/domain"
	"go.uber.org/zap"

	"github.com/carlosarguelles/jellygo/internal/media"
)

type CreateLibraryMetadataUseCase struct {
	movieAPI        movieapp.MovieAPI
	movieRepository movieapp.MovieRepository
	imageManager    media.ImageManager
	logger          *zap.Logger
}

func (uc *CreateLibraryMetadataUseCase) Handle(ctx context.Context, lib *libdom.Library) error {
	dirs, err := os.ReadDir(lib.Path)
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	for _, d := range dirs {
		if lib.Type == libdom.TypeMovies {
			wg.Add(1)
			go func(d fs.DirEntry) {
				storedMovie, err := uc.movieRepository.GetByPath(ctx, d.Name())
				if err != nil {
					uc.logger.Error("Error executing movieRepository.`GetByPath`", zap.Error(err))
				}
				if storedMovie != nil {
					uc.logger.Info("Skipping metadata fetching for movie", zap.String("title", storedMovie.Meta.Title))

					wg.Done()
					return
				}
				movieInfo, err := uc.movieAPI.SearchMovieByName(d.Name())
				if err != nil {
					fmt.Println(err)
				}
				pictures, err := uc.movieAPI.GetMovieImages(movieInfo.ID)
				logoPathID, err1 := uc.imageManager.Download(ctx, pictures.Logo)
				if err1 != nil {
					fmt.Println(err)
				}
				bannerPathID, err2 := uc.imageManager.Download(ctx, pictures.Banner)
				if err2 != nil {
					fmt.Println(err)
				}
				movie := moviedom.Movie{
					Path:      d.Name(),
					LibraryID: lib.ID,
					Meta: &moviedom.MovieMeta{
						ID:          movieInfo.ID,
						Title:       movieInfo.Title,
						ReleaseDate: movieInfo.ReleaseDate,
						Pictures: &moviedom.MoviePictures{
							Banner: bannerPathID,
							Logo:   logoPathID,
						},
					},
				}
				if err := uc.movieRepository.Create(ctx, &movie); err != nil {
					fmt.Println(err)
				}
				uc.logger.Info("Metadata fetched for movie", zap.String("title", movie.Meta.Title))
				wg.Done()
			}(d)
		}
	}
	wg.Wait()
	return nil
}

func NewCreateLibraryMetadataUseCase(movieAPI movieapp.MovieAPI, movieRepository movieapp.MovieRepository, imageDownloader media.ImageManager, logger *zap.Logger) *CreateLibraryMetadataUseCase {
	return &CreateLibraryMetadataUseCase{movieAPI, movieRepository, imageDownloader, logger}
}
