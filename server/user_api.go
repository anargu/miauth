package server

import (
	"encoding/json"
	"errors"
	miauthv2 "github.com/anargu/miauth"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type ThidPartyCredential struct {
	ThirdPartyAccountID *string `json:"account_id" binding:"required"`
}

type MiauthCredential struct {
	Password *string `json:"password" binding:"required"`
}

type LoginInputPayload struct {
	Kind       string      `json:"kind" binding:"required"`
	UserRole   string      `json:"role" binding:"required"`
	Username *string `json:"username" valid:"miauth_username"`
	Email    *string `json:"email" valid:"miauth_email"`
	Credentials interface{} `json:"credential" binding:"required"`
}

func LoginEndpoint(c *gin.Context) {
	var credentials json.RawMessage
	input := LoginInputPayload{
		Credentials: &credentials,
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Bad Params", "Wrong parameters")
		return
	}

	var role miauthv2.Role
	if err := miauthv2.DB.Where(&miauthv2.Role{Name: input.UserRole}).First(&role).Error; err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Role does not match", "User was binding with unknown role")
		return
	}

	switch input.Kind {
	case "facebook":
		var thirdPartyCredential ThidPartyCredential
		if err := json.Unmarshal(credentials, &thirdPartyCredential); err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Cannot process input data", "Values sent are invalid")
			return
		}
		if thirdPartyCredential.ThirdPartyAccountID == nil {
			ErrorResponse(c, http.StatusBadRequest, err, "Missed Params", "Missed parameters")
			return
		}
		var flc miauthv2.FacebookLoginCredential
		if err := miauthv2.DB.Where(miauthv2.FacebookLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}).First(&flc).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "FB ID Account not found", "FB ID Account not found")
			return
		}
		var lc miauthv2.LoginCredential
		if err := miauthv2.DB.Where(miauthv2.LoginCredential{LoginCredentialID: flc.ID}).First(&lc).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "FB ID to User Data Relation not found", "User not found")
			return
		}
		var user miauthv2.User
		if err := miauthv2.DB.Where(&miauthv2.User{Base: miauthv2.Base{ID: lc.UserID}}).First(&user).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "User Data Relation not found", "User not found")
			return
		}

		accessToken, expString, err := miauthv2.TokenizeAccessToken(user.ID.String(), user.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		refreshToken, err := miauthv2.TokenizeRefreshToken(user.ID.String(), user.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		// create login session
		session := miauthv2.Session{
			UserID:       user.ID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expString,
		}
		if err := miauthv2.DB.Create(&session).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}

		c.JSON(http.StatusOK, session)
		return
	case "google":
		var thirdPartyCredential ThidPartyCredential
		if err := json.Unmarshal(credentials, &thirdPartyCredential); err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Cannot process input data", "Values sent are invalid")
			return
		}
		if thirdPartyCredential.ThirdPartyAccountID == nil {
			ErrorResponse(c, http.StatusBadRequest, err, "Missed Params", "Missed parameters")
			return
		}
		var glc miauthv2.GoogleLoginCredential
		if err := miauthv2.DB.Where(miauthv2.GoogleLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}).First(&glc).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "FB ID Account not found", "FB ID Account not found")
			return
		}
		var lc miauthv2.LoginCredential
		if err := miauthv2.DB.Where(miauthv2.LoginCredential{LoginCredentialID: glc.ID}).First(&lc).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "FB ID to User Data Relation not found", "User not found")
			return
		}
		var user miauthv2.User
		if err := miauthv2.DB.Where(&miauthv2.User{Base: miauthv2.Base{ID: lc.UserID}}).First(&user).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "User Data Relation not found", "User not found")
			return
		}

		accessToken, expString, err := miauthv2.TokenizeAccessToken(user.ID.String(), user.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		refreshToken, err := miauthv2.TokenizeRefreshToken(user.ID.String(), user.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		// create login session
		session := miauthv2.Session{
			UserID:       user.ID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expString,
		}
		if err := miauthv2.DB.Create(&session).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}

		c.JSON(http.StatusOK, session)
		return
		break
	case "miauth":
		var miauthCredential MiauthCredential
		if err := json.Unmarshal(credentials, &miauthCredential); err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Cannot process input data", "Values sent are invalid")
			return
		}
		var userFound miauthv2.User
		if input.Username == nil || miauthCredential.Password == nil {
			ErrorResponse(c, http.StatusBadRequest, err, "Missed Params", "Missed parameters")
			return
		}
		if err := miauthv2.DB.Where(&miauthv2.User{Username: *input.Username}).First(&userFound).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "No Username found", "Username not found")
			return
		}
		lc, err := userFound.FindCredentialType(miauthv2.MiauthLC)
		if err != nil || lc == nil {
			if err == nil {
				err = errors.New("no credential found")
			}
			ErrorResponse(c, http.StatusInternalServerError, err, "No credential found", "Invalid user data")
			return
		}
		mlc := (*lc).(miauthv2.MiauthLoginCredential)
		if err := miauthv2.ComparePassword(*miauthCredential.Password, mlc.Hash); err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "Password mismatch", "Incorrect Password")
			return
		}

		accessToken, expString, err := miauthv2.TokenizeAccessToken(userFound.ID.String(), userFound.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		refreshToken, err := miauthv2.TokenizeRefreshToken(userFound.ID.String(), userFound.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		// create login session
		session := miauthv2.Session{
			UserID:       userFound.ID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expString,
		}
		if err := miauthv2.DB.Create(&session).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}

		c.JSON(http.StatusOK, session)
		return
	default:
		ErrorResponse(c, http.StatusBadRequest, err, "Incorrect Credential Type", "Wrong parameters. Are you hacker?")
		return
	}
}

type SignupInputPayload struct {
	Kind string `json:"kind" binding:"required"`

	UserRole string `json:"role" binding:"required"`
	Username *string `json:"username" binding:"required" valid:"miauth_username"`
	Email    *string `json:"email" binding:"required,email" valid:"miauth_email"`
	Credentials interface{} `json:"credential" binding:"required"`
	//Password *string `json:"password"`
	//ThirdPartyAccountID *string `json:"account_id"`
}

func SignupEndpoint(c *gin.Context) {
	var credentials json.RawMessage
	input := SignupInputPayload{
		Credentials: &credentials,
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Bad Params", "Wrong parameters")
		return
	}
	var role miauthv2.Role
	if err := miauthv2.DB.Where(&miauthv2.Role{Name: input.UserRole}).First(&role).Error; err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Role does not match", "User was binding with unknown role")
		return
	}
	switch input.Kind {
	case "facebook":
		var thirdPartyCredential ThidPartyCredential
		if err := json.Unmarshal(credentials, &thirdPartyCredential); err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Cannot process input data", "Values sent are invalid")
			return
		}
		if thirdPartyCredential.ThirdPartyAccountID == nil {
			ErrorResponse(c, http.StatusBadRequest, err, "Missed Params", "Missed parameters")
			return
		}
		tx := miauthv2.DB.Begin()
		var user miauthv2.User
		{
			user = miauthv2.User{Username: *input.Username, Email: *input.Email, Role: role}
			if err := miauthv2.DB.Create(&user).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User", "Cannot process your user information. Please try again in a moment")
				tx.Rollback()
				return
			}
			flc := miauthv2.FacebookLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}
			if err := miauthv2.DB.Create(&flc).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				tx.Rollback()
				return
			}
			lc := miauthv2.LoginCredential{
				UserID:              user.ID,
				KindLoginCredential: miauthv2.FacebookLC,
				LoginCredentialID:   flc.ID}
			if err := miauthv2.DB.Create(&lc).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				tx.Rollback()
				return
			}
		}
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"data": user})
		return
	case "google":
		var thirdPartyCredential ThidPartyCredential
		if err := json.Unmarshal(credentials, &thirdPartyCredential); err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Cannot process input data", "Values sent are invalid")
			return
		}
		if thirdPartyCredential.ThirdPartyAccountID == nil {
			ErrorResponse(c, http.StatusBadRequest, err, "Missed Params", "Missed parameters")
			return
		}

		tx := miauthv2.DB.Begin()
		var user miauthv2.User
		{
			user = miauthv2.User{Username: *input.Username, Email: *input.Email, Role: role}
			if err := miauthv2.DB.Create(&user).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User", "Cannot process your user information. Please try again in a moment")
				tx.Rollback()
				return
			}
			glc := miauthv2.GoogleLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}
			if err := miauthv2.DB.Create(&glc).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				tx.Rollback()
				return
			}
			lc := miauthv2.LoginCredential{
				UserID:              user.ID,
				KindLoginCredential: miauthv2.GoogleLC,
				LoginCredentialID:   glc.ID}
			if err := miauthv2.DB.Create(&lc).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				tx.Rollback()
				return
			}
		}
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"data": user})
		return
	case "miauth":
		var miauthCredential MiauthCredential
		if err := json.Unmarshal(credentials, &miauthCredential); err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Cannot process input data", "Values sent are invalid")
			return
		}
		if ok := miauthv2.UsernamePatterns[miauthv2.OnlyAlphanumericNoSpaceValues].MatchString(*input.Username); !ok {
			ErrorResponse(c,
				http.StatusBadRequest,
				err,
				"Invalid Username pattern",
				miauthv2.Config.FieldValidations.Username.InvalidPatternErrorMessage)
			return
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(*miauthCredential.Password), 10)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Hashing Password Error", "Cannot process your password. Please try a different password or try again in a moment")
			return
		}

		tx := miauthv2.DB.Begin()
		var user miauthv2.User
		{
			user = miauthv2.User{Username: *input.Username, Email: *input.Email, Role: role}
			if err := miauthv2.DB.Create(&user).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User", "Cannot process your user information. Please try again in a moment")
				tx.Rollback()
				return
			}
			mlc := miauthv2.MiauthLoginCredential{Hash: string(hashed)}
			if err := miauthv2.DB.Create(&mlc).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				tx.Rollback()
				return
			}
			lc := miauthv2.LoginCredential{
				UserID:              user.ID,
				KindLoginCredential: miauthv2.MiauthLC,
				LoginCredentialID:   mlc.ID}
			if err := miauthv2.DB.Create(&lc).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				tx.Rollback()
				return
			}
		}
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{"data": user})
		return
	default:
		ErrorResponse(c, http.StatusBadRequest, err, "Incorrect Credential Type", "Wrong parameters. Are you hacker?")
	}
}

type VerifyInputPayload struct {
	AccessToken string `form:"access_token" binding:"required"`
}

func verifyEndpoint(c *gin.Context) {
	var queryParam VerifyInputPayload
	err := c.ShouldBindQuery(&queryParam)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Bad Params", "Wrong parameters")
		return
	}
	if isValid, err := miauthv2.VerifyAccessToken(queryParam.AccessToken); err != nil || !isValid {
		if err == nil {
			err = errors.New("invalid token")
		}
		ErrorResponse(c, http.StatusUnauthorized, err, "Invalid token", "Invalid token")
		return
	}
}

type RefreshTokenInputPayload struct {
	GrantType    string `json:"grant_type" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func tokenRefreshEndpoint(c *gin.Context) {
	var input RefreshTokenInputPayload
	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Bad Params", "Wrong parameters")
		return
	}
	if input.GrantType != "refresh_token" {
		ErrorResponse(c, http.StatusBadRequest, errors.New("invalid grant_type"), "invalid grant_type", "Wrong parameters")
		return
	}

	if isOk, err := miauthv2.VerifyRefreshToken(input.RefreshToken); err != nil || !isOk {
		if err == nil {
			err = errors.New("invalid refresh token")
		}
		ErrorResponse(c, http.StatusUnauthorized, err, "Invalid token", "Invalid token")
		return
	}

	var session miauthv2.Session
	if err := miauthv2.DB.Where(&miauthv2.Session{RefreshToken: input.RefreshToken}).First(&session).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, err, "no session found", "Invalid Session")
		return
	}
	var user miauthv2.User
	if err := miauthv2.DB.Where(&miauthv2.User{Base: miauthv2.Base{ID: session.UserID}}).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, err, "no user found", "Invalid User")
		return
	}
	accessToken, expiresIn, err := miauthv2.TokenizeAccessToken(user.ID.String(), user.Email)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "Cannot tokenize access token", "Wrong tokenize process")
		return
	}
	refreshToken, err := miauthv2.TokenizeRefreshToken(user.ID.String(), user.Email)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "Cannot tokenize access token", "Wrong tokenize process")
		return
	}
	newSession := &miauthv2.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		ExpiresIn:    expiresIn,
		Scope:        nil,
	}
	if err := miauthv2.DB.Create(newSession).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create session", "Session was not created. Please try again")
		return
	}
	miauthv2.DB.Delete(&session)

	c.JSON(http.StatusOK, newSession)
}
