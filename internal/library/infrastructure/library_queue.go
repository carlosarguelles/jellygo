package infrastructure

import (
	"context"
	"fmt"

	"github.com/carlosarguelles/jellygo/internal/library/application"
	"github.com/carlosarguelles/jellygo/internal/library/domain"
	movieapp "github.com/carlosarguelles/jellygo/internal/movie/application"
)

type LibraryQueue struct {
	q               chan domain.Library
	movieAPI        movieapp.MovieAPI
	movieRepository movieapp.MovieRepository
}

func NewLibraryQueue(movieAPI movieapp.MovieAPI, movieRepository movieapp.MovieRepository) *LibraryQueue {
	return &LibraryQueue{q: make(chan domain.Library), movieAPI: movieAPI, movieRepository: movieRepository}
}

func (lq *LibraryQueue) Run(ctx context.Context) {
	uc := application.NewCreateLibraryMetadataUseCase(lq.movieAPI, lq.movieRepository)
	for lib := range lq.q {
		fmt.Printf("New Library : %v", lib)
		err := uc.Handle(ctx, &lib)
		if err != nil {
			fmt.Printf("\nerror from use case : %v", err.Error())
		}
	}
}

func (lq *LibraryQueue) Add(lib domain.Library) {
	lq.q <- lib
}
