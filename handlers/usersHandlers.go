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
	Name 		string
	Email 		string
	Password 	string
 }

 type userResponse struct {
	Id 			int		`json:"id"`
	Name 		string	`json:"name"`
	Email 		string	`json:"email"`
 }

 type ChangePasswordRequest struct {
    Password string `json:"password"`
}

 func (h *UsersHandler) FindAll(c *gin.Context) {
	users, err := h.userRepo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("couldn't load users"))
		return
	}

	dtos := make([]userResponse, 0, len(users))
	for _, u := range users {
		r := userResponse{
			Id: 	u.Id,
			Name: 	u.Name,
			Email: 	u.Email,
		}

		dtos = append(dtos, r)
	}
 
	c.JSON(http.StatusOK, dtos)
 }

 func (h *UsersHandler) FindById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user id"))
		return
	}

	user, err := h.userRepo.FindById(c, id)
    if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }

	response := userResponse{
		Id: 	user.Id,
		Name: 	user.Name,
		Email: 	user.Email,
	}

	c.JSON(http.StatusOK, response)
}

 func (h *UsersHandler) Create(c *gin.Context) {
	var request createUserRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid payload"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed to hash password"))
		return
	}

	user := models.User {
		Name: 			request.Name,
		Email: 			request.Email,
		PasswordHash: 	string(passwordHash),
	}

	id, err := h.userRepo.Create(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Couldn't create user"))
		return
	}

	c.JSON(http.StatusOK,  gin.H{"id": id})

 }

 func (h *UsersHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user id"))
		return
	}

	var request models.User
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Couldn't bind json"))
		return
	}


	err = h.userRepo.Update(c, id, request)
	if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }

	c.Status(http.StatusOK)
 }

 func (h *UsersHandler) ChangePasswordHash(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user id"))
		return
	}

	_, err = h.userRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	var request ChangePasswordRequest

	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid payload"))
		return
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Failed to hash password"))
		return
	}

	err = h.userRepo.ChangePasswordHash(c, id, string(newPasswordHash))
	if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }
	c.Status(http.StatusOK)
 }


 func (h *UsersHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid user id"))
		return
	}

	_, err = h.userRepo.FindById(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.NewApiError("User not found"))
		return
	}

	err = h.userRepo.Delete(c, id)
	if err != nil {
        c.JSON(http.StatusNotFound, models.NewApiError(err.Error()))
        return
    }
	c.Status(http.StatusOK)
}
