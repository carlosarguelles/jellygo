package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/carlosarguelles/jellygo/internal/library/domain"
	moviedom "github.com/carlosarguelles/jellygo/internal/movie/domain"
)

type SqliteLibraryRepository struct {
	db *sql.DB
}

func (repo *SqliteLibraryRepository) GetAll(ctx context.Context) ([]*domain.Library, error) {
	rows, err := repo.db.QueryContext(ctx, "select id, path, type from libraries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	libs := make([]*domain.Library, 0)
	for rows.Next() {
		var lib domain.Library
		if err := rows.Scan(&lib.ID, &lib.Path, &lib.Type); err != nil {
			return nil, err
		}
		libs = append(libs, &lib)
	}
	return libs, nil
}

func (repo *SqliteLibraryRepository) GetByID(ctx context.Context, id int) (*domain.Library, error) {
	row := repo.db.QueryRowContext(ctx, "select id, path, type from libraries where id = ?", id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var lib domain.Library
	if err := row.Scan(&lib.ID, &lib.Path, &lib.Type); err != nil {
		return nil, err
	}
	return &lib, nil
}

func (repo *SqliteLibraryRepository) GetByIDWithMovies(ctx context.Context, id int) (*domain.Library, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		`select lib.id,
               lib.path,
               lib.type,
               mov.id,
               mov.path,
               mov.meta
          from libraries lib
         inner join movies mov
            on mov.library_id = lib.id
         where lib.id = ?`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var lib domain.Library
	movies := make([]*moviedom.Movie, 0)
	for rows.Next() {
		var movie moviedom.Movie
		var metaStr []byte
		if err := rows.Scan(
			&lib.ID,
			&lib.Path,
			&lib.Type,
			&movie.ID,
			&movie.Path,
			&metaStr,
		); err != nil {
			return nil, err
		}
		movie.LibraryID = lib.ID
		err := json.Unmarshal(metaStr, &movie.Meta)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}
	lib.Movies = movies
	return &lib, nil
}

func (repo *SqliteLibraryRepository) Create(ctx context.Context, library *domain.Library) error {
	row := repo.db.QueryRowContext(ctx, "insert into libraries (path, type) values (?, ?) returning id", library.Path, library.Type)
	if row.Err() != nil {
		return row.Err()
	}
	row.Scan(&library.ID)
	return nil
}

func NewSqliteLibraryRepository(db *sql.DB) *SqliteLibraryRepository {
	return &SqliteLibraryRepository{db}
}
