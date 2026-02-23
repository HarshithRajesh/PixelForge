// Package handler for calling auth function
package handler

import (
	"net/http"

	"github.com/HarshithRajesh/PixelForge/internal/models"
	"github.com/HarshithRajesh/PixelForge/internal/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var user models.User
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = h.userService.SignUp(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User Create Successfully"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var user models.Login
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = h.userService.Login(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Login Successfull"})
}
