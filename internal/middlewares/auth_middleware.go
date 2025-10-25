package middlewares

import (
	"net/http"
	"shop-near-u/internal/models"
	"shop-near-u/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func requireUserAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	userID, role, err := utils.ParseToken(tokenString)

	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "invalid token")
		c.Abort()
		return
	}

	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil || user.ID == 0 {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "user not found")
		c.Abort()
		return
	}
	if role != models.RoleUser {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "insufficient permissions")
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Set("role", role)

	c.Next()
}

func requireShopOwnerAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	shopID, role, err := utils.ParseToken(tokenString)

	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "invalid token")
		c.Abort()
		return
	}

	var shop models.Shop
	if err := db.Where("id = ?", shopID).First(&shop).Error; err != nil || shop.ID == 0 {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "shop not found")
		c.Abort()
		return
	}
	if role != models.RoleShopOwner {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "insufficient permissions")
		c.Abort()
		return
	}

	c.Set("shop", shop)
	c.Set("role", role)

	c.Next()
}

func requireAdminAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	userID, role, err := utils.ParseToken(tokenString)

	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "invalid token")
		c.Abort()
		return
	}

	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil || user.ID == 0 {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "user not found")
		c.Abort()
		return
	}

	if role != models.RoleAdmin {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "insufficient permissions")
		c.Abort()
		return
	}

	c.Set("user", user)
	c.Set("role", role)

	c.Next()
}
func RequireUserAuth(gormDB *gorm.DB) gin.HandlerFunc {
	db = gormDB
	return requireUserAuth
}

func RequireShopOwnerAuth(gormDB *gorm.DB) gin.HandlerFunc {
	db = gormDB
	return requireShopOwnerAuth
}

func RequireAdminAuth(gormDB *gorm.DB) gin.HandlerFunc {
	db = gormDB
	return requireAdminAuth
}
