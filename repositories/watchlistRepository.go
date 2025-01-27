package repositories

import (
	"context"
	"goozinshe/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WatchlistRepository struct {
	db *pgxpool.Pool
}

func NewWatchlistRepository(conn *pgxpool.Pool) *WatchlistRepository {
	return &WatchlistRepository{db: conn}
}

func (r *WatchlistRepository) FindAll(c context.Context) ([]models.Movie, error) {
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
		return nil, err
	}

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
			return nil, err
		}

		if _, exists := moviesMap[m.Id]; !exists {
			moviesMap[m.Id] = &m
			movies = append(movies, &m)
		}

		moviesMap[m.Id].Genres = append(moviesMap[m.Id].Genres, g)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	watchlistMovies := make([]models.Movie, 0, len(movies))
	for _, v := range movies {
		watchlistMovies = append(watchlistMovies, *v)
	}

	return watchlistMovies, nil
}


func (r *WatchlistRepository) AddToWatchlist(c context.Context, movieId int) error {
	var exists bool
	err := r.db.QueryRow(c, "SELECT EXISTS(SELECT 1 FROM movies WHERE id = $1)", movieId).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("movie not found")
	}

	_, err = r.db.Exec(c, "INSERT INTO watchlist (movie_id) VALUES ($1)", movieId)
	if err != nil {
		return err
	}

	return nil
}

func (r *WatchlistRepository) Delete(c context.Context, movieId int) error {
	tx, err := r.db.Begin(c)
	if err != nil {
		return err
	}

	_, err = tx.Exec(c, "delete from watchlist where movie_id = $1", movieId)
	if err != nil {
		return err
	}

	err = tx.Commit(c)
	if err != nil {
		return err
	}

	return nil
}