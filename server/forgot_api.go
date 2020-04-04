package server

import (
	"errors"
	"fmt"
	miauthv2 "github.com/anargu/miauth"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type ForgotRequestInputPayload struct {
	Username string
	Email string
}
func ForgotRequestEndpoint(c *gin.Context) {
	var input ForgotRequestInputPayload
	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Bad Params", "Wrong parameters")
		return
	}
	if input.Username == "" && input.Email == "" {
		ErrorResponse(c, http.StatusBadRequest, errors.New("empty input"), "Bad Params", "Wrong parameters")
		return
	}
	var user miauthv2.User
	if err := miauthv2.DB.Where(&miauthv2.User{
		Username: input.Username,
		Email: input.Email,
	}).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusBadRequest, errors.New("user not found"), "User Not Found", "User not found")
		return
	}

	resetEmailToken, err := miauthv2.TokenizeResetEmailToken(user.ID.String(), user.Email)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "cannot generate reset email token", "")
		return
	}
	resetLink := fmt.Sprintf("%s?token=%s", miauthv2.Config.PublicForgotPasswordURL, resetEmailToken)
	if err := miauthv2.SendResetPassword(
		&user, resetLink, "Reset Password"); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "cannot generate reset email token", "")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email sent to the user with the instructions to reset password.",
	})
}

type ForgotResutInputPayload struct {
	Token string `form:"token"`
	NewPassword string `json:"new_password" binding:"required"`
	Retypedpassword string `json:"retyped_password" binding:"required"`
}
func forgotResetEndpoint(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "../public/reset_password.html", nil)
		return
	} else if c.Request.Method == http.MethodPost {
		var input ForgotResutInputPayload
		if err := c.ShouldBind(&input); err != nil {
			c.HTML(http.StatusBadRequest, "../public/reset_password_result_error.html", nil)
			return
		}
		if input.NewPassword != input.Retypedpassword {
			c.JSON(http.StatusBadRequest, errors.New("password and retyped-password mismatch"))
			return
		}
		if userId, err := miauthv2.VerifyResetEmailToken(input.Token); err != nil {
			c.JSON(http.StatusUnauthorized, err)
			return
		} else {
			// REMOVING ALL SESSIONS LOGGED IN OF USER
			id, err := uuid.FromString(*userId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
			if err := miauthv2.DB.Where(&miauthv2.Session{UserID:  id}).Delete([]miauthv2.Session{}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}

			var loginCredential miauthv2.LoginCredential
			if err := miauthv2.DB.Where(&miauthv2.LoginCredential{
				UserID: id, KindLoginCredential: miauthv2.MiauthLC}).First(&loginCredential).Error; err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
			hashed, err := miauthv2.HashPassword(input.NewPassword)
			if err != nil || hashed == nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}

			// UPDATING PASSWORD
			if err := miauthv2.DB.Where(&miauthv2.MiauthLoginCredential{
				Base: miauthv2.Base{ID: loginCredential.LoginCredentialID },
			}).Update(miauthv2.MiauthLoginCredential{Hash: *hashed}).Error; err != nil {
				c.HTML(http.StatusBadRequest, "../public/reset_password_result_error.html", nil)
				return
			}
		}

		c.HTML(http.StatusOK, "../public/reset_password_result_success.html", nil)
		return
	}
}