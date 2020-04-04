package server


import (
	"errors"
	"fmt"
	"github.com/anargu/miauth"
	"github.com/gin-gonic/gin"
	"log"
)

var r *gin.Engine

func InitServer() {
	r = gin.Default()
	r.Static("/public", "../public")
	BindAPIRoutes()

	err := r.Run(fmt.Sprintf(":%s", miauth.Config.Port))
	if err != nil {
		log.Fatal(err)
	}
}

func BindAPIRoutes() {
	authAPI := r.Group("/auth")
	{
		authAPI.POST("/login", LoginEndpoint)
		authAPI.POST("/signup", SignupEndpoint)
		authAPI.GET("/verify", verifyEndpoint)
		authAPI.POST("/token/refresh", tokenRefreshEndpoint)
	}
	forgotAPI := r.Group("/forgot")
	{
		forgotAPI.POST("/request", ForgotRequestEndpoint)
		forgotAPI.GET("/reset", forgotResetEndpoint)
		forgotAPI.POST("/reset", forgotResetEndpoint)
	}
	adminAPI := r.Group("/admin")
	{
		adminAPI.PUT("/update/templates", updateTemplatesEndpoint)
		adminAPI.POST("/revoke_all", RevokeAllEndpoint)
	}
}

type ErrorResponsePayload struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	UserMessage      string `json:"user_message"`
}

func ErrorResponse(c *gin.Context, code int, err error, description string, userMessage string) {
	var _err = err
	if _err == nil {
		_err = errors.New(description)
	}
	c.JSON(code, ErrorResponsePayload{
		Error:            _err.Error(),
		ErrorDescription: description,
		UserMessage:      userMessage,
	})
	return
}
