package application

import (
	"context"

	"github.com/carlosarguelles/jellygo/internal/movie/domain"
)

type MovieRepository interface {
	GetByID(context.Context, int) (*domain.Movie, error)
	GetByPath(context.Context, string) (*domain.Movie, error)
	GetAll(context.Context) ([]*domain.Movie, error)
	Create(context.Context, *domain.Movie) error
	Update(context.Context, *domain.Movie) error
	Delete(context.Context, int) error
}
