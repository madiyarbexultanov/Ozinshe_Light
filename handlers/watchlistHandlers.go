package handlers

import (
	"goozinshe/repositories"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
	"goozinshe/models"
)

type WatchlistHandler struct {
	watchlistRepo *repositories.WatchlistRepository
}

func NewWatchlistHandler(
	watchlistRepo *repositories.WatchlistRepository) *WatchlistHandler {
	return &WatchlistHandler{
		watchlistRepo: watchlistRepo,
	}
}

func (h *WatchlistHandler) FindAll(c *gin.Context) {
	movies, err := h.watchlistRepo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *WatchlistHandler) AddToWatchlist(c *gin.Context) {
	idStr := c.Param("movieId")
	movieId, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	err = h.watchlistRepo.AddToWatchlist(c, movieId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Movie added to watchlist"})
}

func (h *WatchlistHandler) Delete(c *gin.Context) {
	idStr := c.Param("movieId")
	movieid, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid Movie Id"))
		return
	}

	err = h.watchlistRepo.Delete(c, movieid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}