package handlers

import (
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	userRepo *repositories.UsersRepository
}

func NewUsersHandler(repo *repositories.UsersRepository) *UsersHandler {
	return &UsersHandler{userRepo: repo}
}

type createUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type userResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

func (h *UsersHandler) FindAll(c *gin.Context) {
	users, err := h.userRepo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("couldn't load users"))
		return
	}
	dtos := make([]userResponse, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, userResponse{Id: u.Id, Name: u.Name, Email: u.Email})
	}
	c.JSON(http.StatusOK, dtos)
}

func (h *UsersHandler) FindById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user id"))
		return
	}

	user, err := h.userRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, userResponse{Id: user.Id, Name: user.Name, Email: user.Email})
}

func (h *UsersHandler) Create(c *gin.Context) {
	var request createUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed to hash password"))
		return
	}

	id, err := h.userRepo.Create(c, models.User{
		Name: request.Name, Email: request.Email, PasswordHash: string(passwordHash),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Couldn't create user"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *UsersHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user id"))
		return
	}

	var request models.User
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Couldn't bind json"))
		return
	}

	if err := h.userRepo.Update(c, id, request); err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

func (h *UsersHandler) ChangePasswordHash(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user id"))
		return
	}

	if _, err := h.userRepo.FindById(c, id); err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	var request ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed to hash password"))
		return
	}

	if err := h.userRepo.ChangePasswordHash(c, id, string(newPasswordHash)); err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

func (h *UsersHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user id"))
		return
	}

	if _, err := h.userRepo.FindById(c, id); err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	if err := h.userRepo.Delete(c, id); err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
		return
	}
	c.Status(http.StatusOK)
}
