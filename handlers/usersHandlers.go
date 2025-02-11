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

type updateUserRequest struct {
	Name  string
	Email string
}

type userResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

// FindAll godoc
// @Tags users
// @Summary      Get users list
// @Accept       json
// @Produce      json
// @Success      200  {array} handlers.userResponse "OK"
// @Failure   	 500  {object} models.ApiError
// @Router       /users [get]
// @Security Bearer
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

// FindById godoc
// @Tags users
// @Summary      Find users by id
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Success      200  {array} handlers.userResponse "OK"
// @Failure   	 400  {object} models.ApiError "Invalid user id"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id} [get]
// @Security Bearer
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

// Create godoc
// @Tags users
// @Summary      Create user
// @Accept       json
// @Produce      json
// @Param request body handlers.createUserRequest true "User data"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 500  {object} models.ApiError
// @Router       /users [post]
// @Security Bearer
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

// Update godoc
// @Tags users
// @Summary      Update user
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Param request body handlers.updateUserRequest true "User data"
// @Success      200  {object} object{id=int} "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id} [put]
// @Security Bearer
func (h *UsersHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user id"))
		return
	}

	var request updateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Couldn't bind json"))
		return
	}

	user, err := h.userRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	user.Name = request.Name
	user.Email = request.Email

	if err := h.userRepo.Update(c, id, user); err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
		return
	}

	c.Status(http.StatusOK)
}

// ChangePassword godoc
// @Tags users
// @Summary      Change user password
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Param request body handlers.ChangePasswordRequest true "Password data"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id}/changePassword [patch]
// @Security Bearer
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

// Delete godoc
// @Tags users
// @Summary      Delete user
// @Accept       json
// @Produce      json
// @Param id path int true "User id"
// @Success      200  "OK"
// @Failure   	 400  {object} models.ApiError "Invalid data"
// @Failure   	 404  {object} models.ApiError "User not found"
// @Failure   	 500  {object} models.ApiError
// @Router       /users/{id} [delete]
// @Security Bearer
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
