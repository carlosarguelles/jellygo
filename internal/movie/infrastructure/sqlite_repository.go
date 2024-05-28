package infrastructure

import (
	"context"
	"database/sql"

	"github.com/carlosarguelles/jellygo/internal/movie/domain"
)

type SqliteMovieRepository struct {
	db *sql.DB
}

func (repo *SqliteMovieRepository) GetAll(ctx context.Context) ([]*domain.Movie, error) {
	rows, err := repo.db.QueryContext(ctx, "select id, title, release_date, path, library_id from movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	movies := make([]*domain.Movie, 0)
	for rows.Next() {
		var movie domain.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Path, &movie.LibraryID); err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	return movies, nil
}

func (repo *SqliteMovieRepository) GetByID(ctx context.Context, id int) (*domain.Movie, error) {
	row := repo.db.QueryRowContext(ctx, "select id, title, release_date, path, library_id from movies where id = ?", id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var movie domain.Movie
	if err := row.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.Path, &movie.LibraryID); err != nil {
		return nil, err
	}
	return &movie, nil
}

func (repo *SqliteMovieRepository) Create(ctx context.Context, movie *domain.Movie) error {
	row := repo.db.QueryRowContext(
		ctx,
		"insert into movies (title, release_date, path, library_id) values (?, ?, ?, ?) returning id",
		movie.Title,
		movie.ReleaseDate,
		movie.Path,
		movie.LibraryID,
	)
	if row.Err() != nil {
		return row.Err()
	}
	row.Scan(&movie.ID)
	return nil
}

func NewSqliteMovieRepository(db *sql.DB) *SqliteMovieRepository {
	return &SqliteMovieRepository{db}
}
