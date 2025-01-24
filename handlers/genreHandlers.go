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


func NewGenreHandler(
	genresRepo *repositories.GenresRepository,) *GenresHandler{
	return &GenresHandler{
		genresRepo: genresRepo,
	}
}

func (h *GenresHandler) FindAll(c *gin.Context) {
	genres, err := h.genresRepo.FindAll(c)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, genres)
}

func (h *GenresHandler) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Genre Id"))
		return
	}

	genre, err := h.genresRepo.FindById(c, id)
    if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }

	c.JSON(http.StatusOK, genre)
}

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

func (h *GenresHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Genre Id"))
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

func (h *GenresHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Genre Id"))
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


