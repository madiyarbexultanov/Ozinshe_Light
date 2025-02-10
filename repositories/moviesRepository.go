package repositories

import (
	"context"
	"fmt"
	"goozinshe/logger"
	"goozinshe/models"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type MoviesRepository struct {
	db *pgxpool.Pool
}

func NewMoviesRepository(conn *pgxpool.Pool) *MoviesRepository {
	return &MoviesRepository{db: conn}
}

func (r *MoviesRepository) FindById(c context.Context, id int) (models.Movie, error) {
	sql :=
		`
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
join movies_genres mg on mg.movie_id = m.id
join genres g on mg.genre_id  = g.id
where m.id = $1
	`

	logger := logger.GetLogger()

	rows, err := r.db.Query(c, sql, id)
	defer rows.Close()
	if err != nil {
		logger.Error("Could not query database", zap.String("db_msg", err.Error()))
		return models.Movie{}, err
	}

	var movie *models.Movie

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
			logger.Error("Could not scan query row", zap.String("db_msg", err.Error()))
			return models.Movie{}, err
		}

		if movie != nil {
			m = *movie
		}

		m.Genres = append(m.Genres, g)
		movie = &m
	}

	err = rows.Err()
	if err != nil {
		logger.Error(err.Error())
		return models.Movie{}, err
	}

	return *movie, nil
}

func (r *MoviesRepository) FindAll(c context.Context, filters models.MovieFilters) ([]models.Movie, error) {
	logger := logger.GetLogger()

	sql := `select m.id, m.title, m.description, m.release_year, m.director, m.rating, m.is_watched, m.trailer_url, m.poster_url, g.id, g.title from movies m join movies_genres mg on mg.movie_id = m.id join genres g on mg.genre_id  = g.id where 1=1`
	params := pgx.NamedArgs{}

	if filters.SearchTerm != "" {
		sql = fmt.Sprintf("%s and m.title ilike @s", sql)
		params["s"] = fmt.Sprintf("%%%s%%", filters.SearchTerm)
	}

	if filters.GenreId != "" {
		sql = fmt.Sprintf("%s and g.id = @genreId", sql)
		params["genreId"] = filters.GenreId
	}

	if filters.IsWatched != "" {
		isWatched, _ := strconv.ParseBool(filters.IsWatched)
		sql = fmt.Sprintf("%s and m.is_watched = @isWatched", sql)
		params["isWatched"] = isWatched
	}

	if filters.Sort != "" {
		identifier := pgx.Identifier{filters.Sort}
		sql = fmt.Sprintf("%s order by m.%s", sql, identifier.Sanitize())
	}

	rows, err := r.db.Query(c, sql, params)
	if err != nil {
		logger.Error("Could not query database", zap.String("db_msg", err.Error()))
		return nil, err
	}
	defer rows.Close()

	movies := make([]*models.Movie, 0)
	moviesMap := make(map[int]*models.Movie)

	for rows.Next() {
		var m models.Movie
		var g models.Genre

		err := rows.Scan(&m.Id, &m.Title, &m.Description, &m.ReleaseYear, &m.Director, &m.Rating, &m.IsWatched, &m.TrailerUrl, &m.PosterUrl, &g.Id, &g.Title)
		if err != nil {
			logger.Error("Could not scan row", zap.String("db_msg", err.Error()))
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
		logger.Error("Error iterating rows", zap.String("db_msg", err.Error()))
		return nil, err
	}

	concreteMovies := make([]models.Movie, 0, len(movies))
	for _, v := range movies {
		concreteMovies = append(concreteMovies, *v)
	}

	logger.Info("Successfully retrieved movies", zap.Int("count", len(concreteMovies)))
	return concreteMovies, nil
}

func (r *MoviesRepository) Create(c context.Context, movie models.Movie) (int, error) {
	logger := logger.GetLogger()

	tx, err := r.db.Begin(c)
	if err != nil {
		logger.Error("Could not begin transaction", zap.String("db_msg", err.Error()))
		return 0, err
	}

	var id int
	row := tx.QueryRow(c, "insert into movies(title, description, release_year, director, trailer_url, poster_url) values($1, $2, $3, $4, $5, $6) returning id", movie.Title, movie.Description, movie.ReleaseYear, movie.Director, movie.TrailerUrl, movie.PosterUrl)
	err = row.Scan(&id)
	if err != nil {
		logger.Error("Could not insert movie", zap.String("db_msg", err.Error()))
		return 0, err
	}

	for _, genre := range movie.Genres {
		_, err = tx.Exec(c, "insert into movies_genres(movie_id, genre_id) values($1, $2)", id, genre.Id)
		if err != nil {
			logger.Error("Could not insert movie genre", zap.String("db_msg", err.Error()))
			return 0, err
		}
	}

	err = tx.Commit(c)
	if err != nil {
		logger.Error("Could not commit transaction", zap.String("db_msg", err.Error()))
		return 0, err
	}

	logger.Info("Successfully created movie", zap.Int("movie_id", id))
	return id, nil
}

func (r *MoviesRepository) Update(c context.Context, id int, updatedMovie models.Movie) error {
	logger := logger.GetLogger()
	logger.Info("Starting transaction for updating movie", zap.Int("movie_id", id))

	tx, err := r.db.Begin(c)
	if err != nil {
		logger.Error("Could not begin transaction", zap.Error(err))
		return err
	}

	_, err = tx.Exec(
		c,
		`
update movies
set 
title = $1,
description = $2,
release_year = $3,
director = $4,
trailer_url = $5,
poster_url = $6
where id = $7
		`,
		updatedMovie.Title,
		updatedMovie.Description,
		updatedMovie.ReleaseYear,
		updatedMovie.Director,
		updatedMovie.TrailerUrl,
		updatedMovie.PosterUrl,
		id)
	if err != nil {
		logger.Error("Could not update movie", zap.Error(err))
		return err
	}

	_, err = tx.Exec(c, "delete from movies_genres where movie_id = $1", id)
	if err != nil {
		logger.Error("Could not delete movie genres", zap.Error(err))
		return err
	}
	for _, genre := range updatedMovie.Genres {
		_, err = tx.Exec(c, "insert into movies_genres(movie_id, genre_id) values($1, $2)", id, genre.Id)
		if err != nil {
			logger.Error("Could not insert movie genre", zap.Error(err))
			return err
		}
	}

	err = tx.Commit(c)
	if err != nil {
		logger.Error("Could not commit transaction", zap.Error(err))
		return err
	}

	logger.Info("Successfully updated movie", zap.Int("movie_id", id))
	return nil
}

func (r *MoviesRepository) Delete(c context.Context, id int) error {
	logger := logger.GetLogger()
	logger.Info("Starting transaction for deleting movie", zap.Int("movie_id", id))

	tx, err := r.db.Begin(c)
	if err != nil {
		logger.Error("Could not begin transaction", zap.Error(err))
		return err
	}

	_, err = tx.Exec(c, "delete from movies_genres where movie_id = $1", id)
	if err != nil {
		logger.Error("Could not delete movie genres", zap.Error(err))
		return err
	}

	_, err = tx.Exec(c, "delete from movies where id = $1", id)
	if err != nil {
		logger.Error("Could not delete movie", zap.Error(err))
		return err
	}

	err = tx.Commit(c)
	if err != nil {
		logger.Error("Could not commit transaction", zap.Error(err))
		return err
	}

	logger.Info("Successfully deleted movie", zap.Int("movie_id", id))
	return nil
}

func (r *MoviesRepository) SetRating(c context.Context, id int, rating int) error {
	logger := logger.GetLogger()
	logger.Info("Updating movie rating", zap.Int("movie_id", id), zap.Int("rating", rating))

	_, err := r.db.Exec(c, "update movies set rating = $1 where id = $2", rating, id)
	if err != nil {
		logger.Error("Could not update movie rating", zap.Error(err))
		return err
	}

	logger.Info("Successfully updated movie rating", zap.Int("movie_id", id), zap.Int("rating", rating))
	return nil
}

func (r *MoviesRepository) SetWatched(c context.Context, id int, isWatched bool) error {
	logger := logger.GetLogger()
	logger.Info("Updating movie watch status", zap.Int("movie_id", id), zap.Bool("is_watched", isWatched))

	_, err := r.db.Exec(c, "update movies set is_watched = $1 where id = $2", isWatched, id)
	if err != nil {
		logger.Error("Could not update movie watch status", zap.Error(err))
		return err
	}

	logger.Info("Successfully updated movie watch status", zap.Int("movie_id", id), zap.Bool("is_watched", isWatched))
	return nil
}