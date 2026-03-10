// Package handler for calling auth function
package handler

import (
	"context"
	"net/http"

	"github.com/HarshithRajesh/PixelForge/internal/config"
	"github.com/HarshithRajesh/PixelForge/internal/middleware"
	"github.com/HarshithRajesh/PixelForge/internal/models"
	"github.com/HarshithRajesh/PixelForge/internal/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService user.UserService
	rds         *config.Redis
}

func NewUserHandler(userService user.UserService, rds *config.Redis) *UserHandler {
	return &UserHandler{userService: userService, rds: rds}
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
	token, err := h.userService.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	middleware.SetAuthCookies(c, token)
	c.JSON(http.StatusOK, gin.H{"message": "Login Successfull"})
}

func (h *UserHandler) Logout(c *gin.Context) {
	acc, _ := c.Cookie("access_token")
	ref, _ := c.Cookie("refresh_token")
	ctx := context.Background()

	if acc != "" {
		if claims, err := middleware.ParseAccess(acc); err == nil {
			_ = h.rds.DelJTI(ctx, "access:"+claims.ID)
		}
	}
	if ref != "" {
		if claims, err := middleware.ParseRefresh(ref); err == nil {
			_ = h.rds.DelJTI(ctx, "refresh:"+claims.ID)
		}
	}
	middleware.ClearAuthCookies(c)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
