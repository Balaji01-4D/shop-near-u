package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponseWithMessage(c *gin.Context, statusCode int, message string, error string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Error:   error,
	})
}

func ErrorResponseSimple(c *gin.Context, statusCode int, error string) {
	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Message: "Request failed",
		Error:   error,
	})
}
