package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/carlosarguelles/jellygo/internal/movie/domain"
)

type SqliteMovieRepository struct {
	db *sql.DB
}

func (repo *SqliteMovieRepository) GetAll(ctx context.Context) ([]*domain.Movie, error) {
	rows, err := repo.db.QueryContext(ctx, "select id, path, library_id, meta from movies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	movies := make([]*domain.Movie, 0)
	for rows.Next() {
		var movie domain.Movie
		var metaStr []byte
		if err := rows.Scan(&movie.ID, &movie.Path, &movie.LibraryID, &metaStr); err != nil {
			return nil, err
		}
		err := json.Unmarshal(metaStr, &movie.Meta)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	return movies, nil
}

func (repo *SqliteMovieRepository) GetByID(ctx context.Context, id int) (*domain.Movie, error) {
	row := repo.db.QueryRowContext(ctx, "select id, path, library_id, meta from movies where id = ?", id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var movie domain.Movie
	var metaStr []byte
	err := row.Scan(&movie.ID, &movie.Path, &movie.LibraryID, &metaStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	err = json.Unmarshal(metaStr, &movie.Meta)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (repo *SqliteMovieRepository) GetByPath(ctx context.Context, path string) (*domain.Movie, error) {
	row := repo.db.QueryRowContext(ctx, "select id, path, library_id, meta from movies where path = ?", path)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var movie domain.Movie
	var metaStr []byte
	err := row.Scan(&movie.ID, &movie.Path, &movie.LibraryID, &metaStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	err = json.Unmarshal(metaStr, &movie.Meta)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (repo *SqliteMovieRepository) Create(ctx context.Context, movie *domain.Movie) error {
	metaJSON, err := json.Marshal(movie.Meta)
	if err != nil {
		return err
	}
	row := repo.db.QueryRowContext(
		ctx,
		"insert into movies (path, library_id, meta) values (?, ?, ?) returning id",
		movie.Path,
		movie.LibraryID,
		metaJSON,
	)
	if row.Err() != nil {
		return row.Err()
	}
	row.Scan(&movie.ID)
	return nil
}

func (repo *SqliteMovieRepository) Update(ctx context.Context, movie *domain.Movie) error {
	row := repo.db.QueryRowContext(
		ctx,
		"select id from movies where id = ?",
		movie.ID,
	)
	if err := row.Scan(&movie.ID); err != nil {
		return err
	}
	metaJSON, err := json.Marshal(movie.Meta)
	if err != nil {
		return err
	}
	row = repo.db.QueryRowContext(
		ctx,
		"update movies set meta = ? where movies.id = ?",
		metaJSON,
		movie.ID,
	)
	return row.Err()
}

func (repo *SqliteMovieRepository) Delete(ctx context.Context, id int) error {
	_, err := repo.db.ExecContext(ctx, "delete from movies where movies.id = ?", id)
	return err
}

func NewSqliteMovieRepository(db *sql.DB) *SqliteMovieRepository {
	return &SqliteMovieRepository{db}
}
