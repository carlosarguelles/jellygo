package application

import (
	"context"

	"github.com/carlosarguelles/jellygo/internal/library/domain"
)

type LibraryRepository interface {
	GetByID(context.Context, int) (*domain.Library, error)
	GetByIDWithMovies(context.Context, int) (*domain.Library, error)
	GetAll(context.Context) ([]*domain.Library, error)
	Create(context.Context, *domain.Library) error
}
