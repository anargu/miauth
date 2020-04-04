package server

import (
	"errors"
	"fmt"
	miauthv2 "github.com/anargu/miauth"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"path/filepath"
)

func updateTemplatesEndpoint(c *gin.Context) {
	templates := map[string]string{
		"reset_password": "reset_password.html",
		"reset_password_result_success": "reset_password_result_success.html",
		"reset_password_result_error": "reset_password_result_error.html",
	}
	for k, v := range templates {
		file, err := c.FormFile(k)
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.New("bad input payload"))
			continue
		}

		extType := filepath.Ext(file.Filename)
		if extType != ".html" {
			c.JSON(http.StatusBadRequest, errors.New("only html is allowed"))
			return
		}
		if file.Filename != v {
			c.JSON(http.StatusBadRequest,
				errors.New(fmt.Sprintf("filename %s does not match with key %s", file.Filename, k)))
			return
		}

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, errors.New("file corrupted"))
			return
		}

		if err := c.SaveUploadedFile(file, fmt.Sprintf("../public/%s", file.Filename)); err != nil {
			c.JSON(http.StatusInternalServerError, errors.New("cannot be possible upload template file "))
			return
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"files_uploaded": templates,
	})
}

type RevokeAllInputPayload struct {
	UserID string `json:"userId" binding:"required"`
}
func RevokeAllEndpoint(c *gin.Context) {
	var input RevokeAllInputPayload
	err := c.ShouldBindJSON(&input)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Bad Params", "Wrong parameters")
		return
	}

	uuidString, err := uuid.FromString(input.UserID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "User ID Corrupted", "User cannot be identified")
		return
	}
	session := miauthv2.Session{UserID: uuidString}
	if err := miauthv2.DB.Where(&session).Delete(&miauthv2.Session{}).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "Cannot delete all sessions", "Sessions weren't deleted.")
		return
	}
	var deletedSessions []miauthv2.Session
	if err := miauthv2.DB.Unscoped().Where(&session).Find(&deletedSessions).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "Cannot retrieve all deleted sessions", "Sessions were deleted but not retrieved.")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"sessions_deleted": deletedSessions,
	})
}
