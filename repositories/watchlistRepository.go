package repositories

import (
	"context"
	"fmt"
	"goozinshe/logger"
	"goozinshe/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type WatchlistRepository struct {
	db *pgxpool.Pool
}

func NewWatchlistRepository(conn *pgxpool.Pool) *WatchlistRepository {
	return &WatchlistRepository{db: conn}
}

func (r *WatchlistRepository) FindAll(c context.Context) ([]models.Movie, error) {
	logger := logger.GetLogger()
	logger.Info("Fetching all movies from watchlist")

	sql := `
    select 
        m.id,
        m.title,
        m.description,
        m.release_year,
        m.director,
        m.rating,
        m.is_watched,
        m.trailer_url,
        m.poster_url,
        g.id,
        g.title
    from movies m
    join watchlist w on w.movie_id = m.id
    join movies_genres mg on mg.movie_id = m.id
    join genres g on mg.genre_id = g.id
    `

	rows, err := r.db.Query(c, sql)
	if err != nil {
		logger.Error("Error querying watchlist movies", zap.Error(err))
		return nil, err
	}

defer rows.Close()

	movies := make([]*models.Movie, 0)
	moviesMap := make(map[int]*models.Movie)

	for rows.Next() {
		var m models.Movie
		var g models.Genre

		err := rows.Scan(
			&m.Id,
			&m.Title,
			&m.Description,
			&m.ReleaseYear,
			&m.Director,
			&m.Rating,
			&m.IsWatched,
			&m.TrailerUrl,
			&m.PosterUrl,
			&g.Id,
			&g.Title,
		)
		if err != nil {
			logger.Error("Error scanning movie row", zap.Error(err))
			return nil, err
		}

		if _, exists := moviesMap[m.Id]; !exists {
			moviesMap[m.Id] = &m
			movies = append(movies, &m)
		}

		moviesMap[m.Id].Genres = append(moviesMap[m.Id].Genres, g)
	}

	if err = rows.Err(); err != nil {
		logger.Error("Error iterating over watchlist movies", zap.Error(err))
		return nil, err
	}

	watchlistMovies := make([]models.Movie, 0, len(movies))
	for _, v := range movies {
		watchlistMovies = append(watchlistMovies, *v)
	}

	logger.Info("Successfully fetched watchlist movies", zap.Int("count", len(watchlistMovies)))
	return watchlistMovies, nil
}

func (r *WatchlistRepository) AddToWatchlist(c context.Context, movieId int) error {
	logger := logger.GetLogger()
	logger.Info("Adding movie to watchlist", zap.Int("movie_id", movieId))

	var exists bool
	err := r.db.QueryRow(c, "SELECT EXISTS(SELECT 1 FROM movies WHERE id = $1)", movieId).Scan(&exists)
	if err != nil {
		logger.Error("Error checking if movie exists", zap.Error(err))
		return err
	}
	if !exists {
		logger.Warn("Movie not found", zap.Int("movie_id", movieId))
		return fmt.Errorf("movie not found")
	}

	_, err = r.db.Exec(c, "INSERT INTO watchlist (movie_id) VALUES ($1)", movieId)
	if err != nil {
		logger.Error("Error inserting movie into watchlist", zap.Error(err))
		return err
	}

	logger.Info("Successfully added movie to watchlist", zap.Int("movie_id", movieId))
	return nil
}

func (r *WatchlistRepository) Delete(c context.Context, movieId int) error {
	logger := logger.GetLogger()
	logger.Info("Removing movie from watchlist", zap.Int("movie_id", movieId))

	tx, err := r.db.Begin(c)
	if err != nil {
		logger.Error("Error beginning transaction", zap.Error(err))
		return err
	}

	_, err = tx.Exec(c, "delete from watchlist where movie_id = $1", movieId)
	if err != nil {
		tx.Rollback(c)
		logger.Error("Error deleting movie from watchlist", zap.Error(err))
		return err
	}

	err = tx.Commit(c)
	if err != nil {
		logger.Error("Error committing transaction", zap.Error(err))
		return err
	}

	logger.Info("Successfully removed movie from watchlist", zap.Int("movie_id", movieId))
	return nil
}
