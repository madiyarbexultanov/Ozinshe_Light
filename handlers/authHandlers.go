package handlers

import (
	"goozinshe/config"
	"goozinshe/models"
	"goozinshe/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandlers struct {
	usersRepo *repositories.UsersRepository
}

func NewAuthHandlers(usersRepo *repositories.UsersRepository) *AuthHandlers {
	return &AuthHandlers{usersRepo: usersRepo}
}

type SignInRequest struct {
	Email 		string
	Password 	string
}

// SignIn godoc
// @Summary      Authenticate user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body SignInRequest true "User credentials"
// @Success      200 {object} map[string]string "JWT Token"
// @Failure      400 {object} models.ApiError "Invalid payload"
// @Failure      401 {object} models.ApiError "Invalid credentials"
// @Failure      500 {object} models.ApiError "Internal server error"
// @Router       /auth/signin [post]
func (h *AuthHandlers) SignIn(c *gin.Context) {
	var request SignInRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid payload"))
		return
	}

	user, err := h.usersRepo.FindByEmail(c, request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewApiError("Invalid credentials"))
		return
	}

	claims := jwt.RegisteredClaims {
		Subject: strconv.Itoa(user.Id),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.JwtExpiresIn)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Couldn't generate JWT token"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// SignOut godoc
// @Summary      Logout user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200 "OK"
// @Router       /auth/signout [post]
func (h *AuthHandlers) SignOut(c *gin.Context) {
	c.Status(http.StatusOK)
}

// GetUserInfo godoc
// @Summary      Get authenticated user info
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200 {object} userResponse "User info"
// @Failure      500 {object} models.ApiError "Internal server error"
// @Router       /auth/user [get]
// @Security     Bearer
func (h *AuthHandlers) GetUserInfo(c *gin.Context) {
	userId := c.GetInt("userId")
	user, err := h.usersRepo.FindById(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, userResponse{
		Id:		user.Id,
		Email: 	user.Email,
		Name: 	user.Name,
	})
}