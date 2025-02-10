package repositories

import (
	"context"

	"goozinshe/models"
	"goozinshe/logger"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type GenresRepository struct {
	db *pgxpool.Pool
}

func NewGenresRepository(conn *pgxpool.Pool) *GenresRepository {
	return &GenresRepository{db: conn}
} 

func (r *GenresRepository) FindAllByIds(c context.Context, ids []int) ([]models.Genre, error) {
	logger := logger.GetLogger()
	logger.Info("Fetching genres by IDs", zap.Ints("genre_ids", ids))

	rows, err := r.db.Query(c, "select id, title from genres where id = any($1)", ids)
	if err != nil {
		logger.Error("Could not fetch genres", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	genres := make([]models.Genre, 0)
	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.Id, &genre.Title)
		if err != nil {
			logger.Error("Could not scan genre", zap.Error(err))
			return nil, err
		}
		genres = append(genres, genre)
	}

	logger.Info("Successfully fetched genres", zap.Int("count", len(genres)))
	return genres, nil
}

func (r *GenresRepository) FindById(c context.Context, id int) (models.Genre, error) {
	logger := logger.GetLogger()
	logger.Info("Fetching genre by ID", zap.Int("genre_id", id))

	var genre models.Genre
	row := r.db.QueryRow(c, "select id, title from genres where id = $1", id)
	err := row.Scan(&genre.Id, &genre.Title)
	if err != nil {
		logger.Error("Could not fetch genre", zap.Error(err))
		return models.Genre{}, err
	}

	logger.Info("Successfully fetched genre", zap.Int("genre_id", genre.Id))
	return genre, nil
}

func (r *GenresRepository) FindAll(c context.Context) ([]models.Genre, error) {
	logger := logger.GetLogger()
	logger.Info("Fetching all genres")

	rows, err := r.db.Query(c, "select id, title from genres")
	if err != nil {
		logger.Error("Could not fetch genres", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	genres := make([]models.Genre, 0)
	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.Id, &genre.Title)
		if err != nil {
			logger.Error("Could not scan genre", zap.Error(err))
			return nil, err
		}
		genres = append(genres, genre)
	}

	logger.Info("Successfully fetched all genres", zap.Int("count", len(genres)))
	return genres, nil
}

func (r *GenresRepository) Create(c context.Context, genre models.Genre) (int, error) {
	logger := logger.GetLogger()
	logger.Info("Creating new genre", zap.String("title", genre.Title))

	var id int
	row := r.db.QueryRow(c, "insert into genres (title) values ($1) returning id", genre.Title)
	err := row.Scan(&id)
	if err != nil {
		logger.Error("Could not create genre", zap.Error(err))
		return 0, err
	}

	logger.Info("Successfully created genre", zap.Int("genre_id", id))
	return id, nil
}

func (r *GenresRepository) Update(c context.Context, id int, genre models.Genre) error {
	logger := logger.GetLogger()
	logger.Info("Updating genre", zap.Int("genre_id", id), zap.String("title", genre.Title))

	_, err := r.db.Exec(c, "update genres set title=$1 where id=$2", genre.Title, id)
	if err != nil {
		logger.Error("Could not update genre", zap.Error(err))
		return err
	}

	logger.Info("Successfully updated genre", zap.Int("genre_id", id))
	return nil
}

func (r *GenresRepository) Delete(c context.Context, id int) error {
	logger := logger.GetLogger()
	logger.Info("Deleting genre", zap.Int("genre_id", id))

	_, err := r.db.Exec(c, "delete from genres where id=$1", id)
	if err != nil {
		logger.Error("Could not delete genre", zap.Error(err))
		return err
	}

	logger.Info("Successfully deleted genre", zap.Int("genre_id", id))
	return nil
}
