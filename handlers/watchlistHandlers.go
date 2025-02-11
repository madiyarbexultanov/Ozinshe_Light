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

// FindAll godoc
// @Summary      Add movie to watchlist
// @Tags watchlist
// @Accept       json
// @Produce      json
// @Param movieId path int true "Movie id"
// @Success      200 "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /watchlist/:movieId [post]
// @Security Bearer
func (h *WatchlistHandler) FindAll(c *gin.Context) {
	movies, err := h.watchlistRepo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, movies)
}

// AddToWatchlist godoc
// @Summary      Add movie to watchlist
// @Tags watchlist
// @Accept       json
// @Produce      json
// @Param movieId path int true "Movie id"
// @Success      200 "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /watchlist/:movieId [post]
// @Security Bearer
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

// Delete godoc
// @Summary      Remove movie from watchlist
// @Tags watchlist
// @Accept       json
// @Produce      json
// @Param movieId path int true "Movie id"
// @Success      200 "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /watchlist/:movieId [delete]
// @Security Bearer
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