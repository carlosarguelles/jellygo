package infrastructure

import (
	"context"

	"github.com/carlosarguelles/jellygo/internal/library/application"
	"github.com/carlosarguelles/jellygo/internal/library/domain"
	"github.com/carlosarguelles/jellygo/internal/media"
	movieapp "github.com/carlosarguelles/jellygo/internal/movie/application"
	"go.uber.org/zap"
)

type LibraryQueue struct {
	q               chan domain.Library
	movieAPI        movieapp.MovieAPI
	movieRepository movieapp.MovieRepository
	imageDownloader media.ImageManager
	logger          *zap.Logger
}

func NewLibraryQueue(movieAPI movieapp.MovieAPI, movieRepository movieapp.MovieRepository, imageDownloader media.ImageManager, logger *zap.Logger) *LibraryQueue {
	return &LibraryQueue{q: make(chan domain.Library), movieAPI: movieAPI, movieRepository: movieRepository, imageDownloader: imageDownloader, logger: logger}
}

func (lq *LibraryQueue) Run(ctx context.Context) {
	uc := application.NewCreateLibraryMetadataUseCase(lq.movieAPI, lq.movieRepository, lq.imageDownloader, lq.logger)
	for lib := range lq.q {
		lq.logger.Info("New Library", zap.Int("LibraryID", lib.ID))
		err := uc.Handle(ctx, &lib)
		if err != nil {
			lq.logger.Error("Error from use case NewCreateLibraryMetadataUseCase", zap.Error(err))
		}
	}
}

func (lq *LibraryQueue) Add(lib domain.Library) {
	lq.q <- lib
}

func (lq *LibraryQueue) Close() {
	close(lq.q)
}
