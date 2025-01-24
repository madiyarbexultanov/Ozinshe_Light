package handlers

import (
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MoviesHandler struct {
	moviesRepo *repositories.MoviesRepository
	genresRepo *repositories.GenresRepository
}

type createMovieRequest struct {
	Title 		string
	Description string
	ReleaseYear int
	Director 	string
	TrailerUrl 	string
	GenreIds 	[]int
}

type updateMovieRequest struct {
	Title 		string
	Description string
	ReleaseYear int
	Director 	string
	TrailerUrl 	string
	GenreIds 	[]int
}

func NewMoviesHandler(
	genresRepo *repositories.GenresRepository,
	moviesRepo *repositories.MoviesRepository) *MoviesHandler {
	return &MoviesHandler{
		moviesRepo: moviesRepo,
		genresRepo: genresRepo,
	}
}

func (h *MoviesHandler) FindAll(c *gin.Context) {
	movies, err := h.moviesRepo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MoviesHandler) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	movie, err := h.moviesRepo.FindById(c, id)
    if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }

	c.JSON(http.StatusOK, movie)
}

func (h *MoviesHandler) Create(c *gin.Context) {
	var request createMovieRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Couldn't bind json"))
		return
	}

	genres, err := h.genresRepo.FindAllByIds(c, request.GenreIds)
	if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }

	movie := models.Movie {
		Title: 			request.Title,
		Description:	request.Description,
		ReleaseYear:	request.ReleaseYear,
		Director:		request.Director,
		TrailerUrl: 	request.TrailerUrl,
		Genres: 		genres,
	}

	id, err := h.moviesRepo.Create(c, movie)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *MoviesHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var request updateMovieRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Couldn't bind json"))
		return
	}

	genres, err := h.genresRepo.FindAllByIds(c, request.GenreIds)
	if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }

	movie := models.Movie {
		Title: 			request.Title,
		Description:	request.Description,
		ReleaseYear:	request.ReleaseYear,
		Director:		request.Director,
		TrailerUrl: 	request.TrailerUrl,
		Genres: 		genres,
	}

	h.moviesRepo.Update(c, id, movie)

	c.Status(http.StatusOK)
}

func (h *MoviesHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	_, err = h.moviesRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}
	
	h.moviesRepo.Delete(c, id)
	c.Status(http.StatusOK)
}