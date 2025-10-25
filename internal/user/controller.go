package user

import (
	"net/http"
	"shop-near-u/internal/middlewares"
	"shop-near-u/internal/models"
	"shop-near-u/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

func (ctrl *Controller) Register(c *gin.Context) {

	var userDTO UserRegisterDTO
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		utils.ErrorResponseSimple(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := ctrl.service.RegisterUser(&userDTO)
	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := utils.GenerateAccessToken(user.ID, models.RoleUser)
	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.SetCookie(token, 3600*24*30, c)

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"token": token,
	})
}

func (ctrl *Controller) Login(c *gin.Context) {

	var loginDTO UserLoginDTO
	if err := c.ShouldBindJSON(&loginDTO); err != nil {
		utils.ErrorResponseSimple(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := ctrl.service.AuthenticateUser(loginDTO.Email, loginDTO.Password)
	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := utils.GenerateAccessToken(user.ID, models.RoleUser)
	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.SetCookie(token, 3600*24*30, c)

	utils.SuccessResponse(c, http.StatusOK, "User logged in successfully", gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"token": token,
	})
}

func (ctrl *Controller) Me(c *gin.Context) {

	user, exists := c.Get("user")
	if !exists {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	u := user.(models.User)

	utils.SuccessResponse(c, http.StatusOK, "User profile retrieved successfully", gin.H{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
	})
}

func (ctrl *Controller) Logout(c *gin.Context) {

	utils.SetCookie("", -1, c)

	utils.SuccessResponse(c, http.StatusOK, "Successfully logged out", nil)
}

func (ctrl *Controller) ChangePassword(c *gin.Context) {
	var pwdDTO ChangePasswordDTO
	if err := c.ShouldBindJSON(&pwdDTO); err != nil {
		utils.ErrorResponseSimple(c, http.StatusBadRequest, err.Error())
		return
	}

	user, exists := c.Get("user")
	if !exists {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	u := user.(models.User)
	err := ctrl.service.ChangePassword(u.ID, pwdDTO.OldPassword, pwdDTO.NewPassword)
	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}

func (ctrl *Controller) DeleteAccount(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorResponseSimple(c, http.StatusUnauthorized, "unauthorized")
		c.Abort()
		return
	}

	u := user.(models.User)
	err := ctrl.service.DeleteUser(u.ID)
	if err != nil {
		utils.ErrorResponseSimple(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SetCookie("", -1, c)

	utils.SuccessResponse(c, http.StatusOK, "Account deleted successfully", nil)
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	ctrl := NewController(svc)

	users := r.Group("/auth")
	{
		users.POST("/register", ctrl.Register)
		users.POST("/login", ctrl.Login)
		users.GET("/me", middlewares.RequireUserAuth(db), ctrl.Me)
		users.POST("/logout", middlewares.RequireUserAuth(db), ctrl.Logout)
		users.POST("/change-password", middlewares.RequireUserAuth(db), ctrl.ChangePassword)
		users.DELETE("/delete-account", middlewares.RequireUserAuth(db), ctrl.DeleteAccount)
	}
}
