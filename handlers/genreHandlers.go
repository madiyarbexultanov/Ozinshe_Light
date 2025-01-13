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

type createGenreRequest struct {
	Id 		int
	Title 	string
}

type updateGenreRequest struct {
	Id 		int
	Title 	string
}

func NewGenreHandler(
	genresRepo *repositories.GenresRepository,) *GenresHandler{
	return &GenresHandler{
		genresRepo: genresRepo,
	}
}

func (h *GenresHandler) FindAll(c *gin.Context) {
	genres := h.genresRepo.FindAll(c)
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
	var request createGenreRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Couldn't bind json"))
		return
	}

	genre := models.Genre {
		Id: 			request.Id,
		Title: 			request.Title,
	}

	id := h.genresRepo.Create(c, genre)

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

	var request updateGenreRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Couldn't bind json"))
		return
	}

	genre := models.Genre {
		Id:				request.Id,
		Title: 			request.Title,
	}

	h.genresRepo.Update(c, id, genre)

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
	
	h.genresRepo.Delete(c, id)
	c.Status(http.StatusOK)
}


