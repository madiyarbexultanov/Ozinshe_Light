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

// FindById godoc
// @Summary      Find genre by id
// @Tags genres
// @Accept       json
// @Produce      json
// @Param id path int true "Genre ID"
// @Success      200  {object} models.Genre "OK"
// @Failure   	 400  {object} models.ApiError "Validation error"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres/{id} [get]
// @Security Bearer
func (h *GenresHandler) FindAll(c *gin.Context) {
	genres, err := h.genresRepo.FindAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, genres)
}

// FindAll godoc
// @Tags genres
// @Summary      Get genres list
// @Accept       json
// @Produce      json
// @Success      200  {array} models.Genre "OK"
// @Failure   	 400  {object} models.ApiError "Validation error"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres [get]
// @Security Bearer
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
// @Summary      Create genre
// @Tags genres
// @Accept       json
// @Produce      json
// @Param request body models.Genre true "Genre model"
// @Success      200  {object} object{id=int}  "OK"
// @Failure   	 400  {object} models.ApiError "Validation error"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres [post]
// @Security Bearer
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
// @Summary      Update genre
// @Tags genres
// @Accept       json
// @Produce      json
// @Param id path int true "Genre id"
// @Param request body models.Genre true "Genre model"
// @Success      200
// @Failure   	 400  {object} models.ApiError "Validation error"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres/{id} [put]
// @Security Bearer
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
// @Summary      Delete genre
// @Tags genres
// @Accept       json
// @Produce      json
// @Param id path int true "Genre id"
// @Success      200
// @Failure   	 400  {object} models.ApiError "Validation error"
// @Failure   	 500  {object} models.ApiError
// @Router       /genres/{id} [delete]
// @Security Bearer
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


