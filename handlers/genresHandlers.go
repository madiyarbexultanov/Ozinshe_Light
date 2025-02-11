package handlers

import (
	"goozinshe/models"
	"goozinshe/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GenresHandler struct {
	genresRepo *repositories.GenresRepository
}


func NewGenresHandler(
	genresRepo *repositories.GenresRepository,) *GenresHandler{
	return &GenresHandler{
		genresRepo: genresRepo,
	}
}

// FindAll godoc
// @Summary Get all genres
// @Description Retrieves a list of all available genres
// @Tags genres
// @Produce json
// @Success 200 {array} models.Genre
// @Failure 500 {object} models.ApiError
// @Router /genres [get]
func (h *GenresHandler) FindAll(c *gin.Context) {
	genres, err := h.genresRepo.FindAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, genres)
}

// FindById godoc
// @Summary Get a genre by ID
// @Description Retrieves a single genre by its ID
// @Tags genres
// @Produce json
// @Param id path int true "Genre ID"
// @Success 200 {object} models.Genre
// @Failure 400 {object} models.ApiError "Invalid genre id"
// @Failure 404 {object} models.ApiError "Genre not found"
// @Router /genres/{id} [get]
func (h *GenresHandler) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid genre id"))
		return
	}

	genre, err := h.genresRepo.FindById(c, id)
    if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }

	c.JSON(http.StatusOK, genre)
}

// Create godoc
// @Summary Create a new genre
// @Description Adds a new genre to the database
// @Tags genres
// @Accept json
// @Produce json
// @Param genre body models.Genre true "Genre data"
// @Success 200 {object} map[string]int "Created genre ID"
// @Failure 400 {object} models.ApiError "Couldn't bind json"
// @Failure 404 {object} models.ApiError
// @Router /genres [post]
func (h *GenresHandler) Create(c *gin.Context) {
	var createGenre models.Genre

	err := c.BindJSON(&createGenre)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Couldn't bind json"))
		return
	}

	id, err := h.genresRepo.Create(c, createGenre)
	if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// Update godoc
// @Summary Update a genre
// @Description Updates an existing genre by ID
// @Tags genres
// @Accept json
// @Produce json
// @Param id path int true "Genre ID"
// @Param genre body models.Genre true "Updated genre data"
// @Success 200
// @Failure 400 {object} models.ApiError "Invalid genre id or couldn't bind json"
// @Failure 404 {object} models.ApiError "Genre not found"
// @Router /genres/{id} [put]
func (h *GenresHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid genre id"))
		return
	}

	_, err = h.genresRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	var updateGenre models.Genre
	err = c.BindJSON(&updateGenre)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Couldn't bind json"))
		return
	}


	err = h.genresRepo.Update(c, id, updateGenre)
	if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }

	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary Delete a genre
// @Description Deletes a genre by ID
// @Tags genres
// @Param id path int true "Genre ID"
// @Success 200
// @Failure 400 {object} models.ApiError "Invalid genre id"
// @Failure 404 {object} models.ApiError "Genre not found"
// @Router /genres/{id} [delete]
func (h *GenresHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid genre id"))
		return
	}

	_, err = h.genresRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}
	
	err = h.genresRepo.Delete(c, id)
	if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }
	c.Status(http.StatusOK)
}


