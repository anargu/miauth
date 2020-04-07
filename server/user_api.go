package server

import (
	"encoding/json"
	"errors"
	"github.com/anargu/miauth"
	"github.com/dgrijalva/jwt-go"
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
	Kind        string      `json:"kind" binding:"required"`
	UserRole    string      `json:"role" binding:"required"`
	Username    *string     `json:"username" valid:"miauth_username"`
	Email       *string     `json:"email" valid:"miauth_email"`
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

	var role miauth.Role
	if err := miauth.DB.Where(&miauth.Role{Name: input.UserRole}).First(&role).Error; err != nil {
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
		var flc miauth.FacebookLoginCredential
		if err := miauth.DB.Where(miauth.FacebookLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}).First(&flc).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "FB ID Account not found", "FB ID Account not found")
			return
		}
		var lc miauth.LoginCredential
		if err := miauth.DB.Where(miauth.LoginCredential{LoginCredentialID: flc.ID}).First(&lc).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "FB ID to User Data Relation not found", "User not found")
			return
		}
		var user miauth.User
		if err := miauth.DB.Where(&miauth.User{Base: miauth.Base{ID: lc.UserID}}).First(&user).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "User Data Relation not found", "User not found")
			return
		}

		accessToken, expString, err := miauth.TokenizeAccessToken(user.ID.String(), user.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		refreshToken, err := miauth.TokenizeRefreshToken(user.ID.String(), user.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		// create login session
		session := miauth.Session{
			UserID:       user.ID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expString,
		}
		if err := miauth.DB.Create(&session).Error; err != nil {
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
		var glc miauth.GoogleLoginCredential
		if err := miauth.DB.Where(miauth.GoogleLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}).First(&glc).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "FB ID Account not found", "FB ID Account not found")
			return
		}
		var lc miauth.LoginCredential
		if err := miauth.DB.Where(miauth.LoginCredential{LoginCredentialID: glc.ID}).First(&lc).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "FB ID to User Data Relation not found", "User not found")
			return
		}
		var user miauth.User
		if err := miauth.DB.Where(&miauth.User{Base: miauth.Base{ID: lc.UserID}}).First(&user).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "User Data Relation not found", "User not found")
			return
		}

		accessToken, expString, err := miauth.TokenizeAccessToken(user.ID.String(), user.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		refreshToken, err := miauth.TokenizeRefreshToken(user.ID.String(), user.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		// create login session
		session := miauth.Session{
			UserID:       user.ID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expString,
		}
		if err := miauth.DB.Create(&session).Error; err != nil {
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
		var userFound miauth.User
		if input.Username == nil || miauthCredential.Password == nil {
			ErrorResponse(c, http.StatusBadRequest, err, "Missed Params", "Missed parameters")
			return
		}
		if err := miauth.DB.Where(&miauth.User{Username: *input.Username}).First(&userFound).Error; err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "No Username found", "Username not found")
			return
		}
		lc, err := userFound.FindCredentialType(miauth.MiauthLC)
		if err != nil || lc == nil {
			if err == nil {
				err = errors.New("no credential found")
			}
			ErrorResponse(c, http.StatusInternalServerError, err, "No credential found", "Invalid user data")
			return
		}
		mlc := (*lc).(miauth.MiauthLoginCredential)
		if err := miauth.ComparePassword(*miauthCredential.Password, mlc.Hash); err != nil {
			ErrorResponse(c, http.StatusBadRequest, err, "Password mismatch", "Incorrect Password")
			return
		}

		accessToken, expString, err := miauth.TokenizeAccessToken(userFound.ID.String(), userFound.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		refreshToken, err := miauth.TokenizeRefreshToken(userFound.ID.String(), userFound.Email)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, err.Error(), err.Error())
			return
		}
		// create login session
		session := miauth.Session{
			UserID:       userFound.ID,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expString,
		}
		if err := miauth.DB.Create(&session).Error; err != nil {
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

	UserRole    string      `json:"role" binding:"required"`
	Username    *string     `json:"username" binding:"required" valid:"miauth_username"`
	Email       *string     `json:"email" binding:"required,email" valid:"miauth_email"`
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
	var role miauth.Role
	if err := miauth.DB.Where(&miauth.Role{Name: input.UserRole}).First(&role).Error; err != nil {
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
		tx := miauth.DB.Begin()
		var user miauth.User
		{
			user = miauth.User{Username: *input.Username, Email: *input.Email, Role: role}
			if err := miauth.DB.Create(&user).Error; err != nil {
				if _err := miauth.ValidateDuplicateErrorInField(err, "username", nil); _err != nil {
					SendError(c, http.StatusBadRequest, _err)
				} else if _err = miauth.ValidateDuplicateErrorInField(err, "email", nil); _err != nil {
					SendError(c, http.StatusBadRequest, _err)
				} else {
					ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User", "Cannot process your user information. Please try again in a moment")
				}
				tx.Rollback()
				return
			}
			flc := miauth.FacebookLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}
			if err := miauth.DB.Create(&flc).Error; err != nil {
				duplicateErrorMessage := "It seems that Facebook Account have been already signed up."
				if _err := miauth.ValidateDuplicateErrorInField(err, "account_id", &duplicateErrorMessage); _err != nil {
					SendError(c, http.StatusBadRequest, _err)
				} else {
					ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				}
				tx.Rollback()
				return
			}
			lc := miauth.LoginCredential{
				UserID:              user.ID,
				KindLoginCredential: miauth.FacebookLC,
				LoginCredentialID:   flc.ID}
			if err := miauth.DB.Create(&lc).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				tx.Rollback()
				return
			}
		}
		tx.Commit()
		c.JSON(http.StatusOK, user)
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

		tx := miauth.DB.Begin()
		var user miauth.User
		{
			user = miauth.User{Username: *input.Username, Email: *input.Email, Role: role}
			if err := miauth.DB.Create(&user).Error; err != nil {
				duplicateErrorMessage := "It seems that Google Account have been already signed up."
				if _err := miauth.ValidateDuplicateErrorInField(err, "username", &duplicateErrorMessage); _err != nil {
					SendError(c, http.StatusBadRequest, _err)
				} else if _err = miauth.ValidateDuplicateErrorInField(err, "email", nil); _err != nil {
					SendError(c, http.StatusBadRequest, _err)
				} else {
					ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User", "Cannot process your user information. Please try again in a moment")
				}
				tx.Rollback()
				return
			}
			glc := miauth.GoogleLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}
			if err := miauth.DB.Create(&glc).Error; err != nil {
				if _err := miauth.ValidateDuplicateErrorInField(err, "account_id", nil); _err != nil {
					SendError(c, http.StatusBadRequest, _err)
				} else {
					ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				}
				tx.Rollback()
				return
			}
			lc := miauth.LoginCredential{
				UserID:              user.ID,
				KindLoginCredential: miauth.GoogleLC,
				LoginCredentialID:   glc.ID}
			if err := miauth.DB.Create(&lc).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				tx.Rollback()
				return
			}
		}
		tx.Commit()
		c.JSON(http.StatusOK, user)
		return
	case "miauth":
		var miauthCredential MiauthCredential
		if err := json.Unmarshal(credentials, &miauthCredential); err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Cannot process input data", "Values sent are invalid")
			return
		}
		if ok := miauth.UsernamePatterns[miauth.OnlyAlphanumericNoSpaceValues].MatchString(*input.Username); !ok {
			ErrorResponse(c,
				http.StatusBadRequest,
				err,
				"Invalid Username pattern",
				miauth.Config.FieldValidations.Username.InvalidPatternErrorMessage)
			return
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(*miauthCredential.Password), 10)
		if err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Hashing Password Error", "Cannot process your password. Please try a different password or try again in a moment")
			return
		}

		tx := miauth.DB.Begin()
		var user miauth.User
		{
			user = miauth.User{Username: *input.Username, Email: *input.Email, Role: role}
			if err := miauth.DB.Create(&user).Error; err != nil {
				if _err := miauth.ValidateDuplicateErrorInField(err, "username", nil); _err != nil {
					SendError(c, http.StatusBadRequest, _err)
				} else if _err = miauth.ValidateDuplicateErrorInField(err, "email", nil); _err != nil {
					SendError(c, http.StatusBadRequest, _err)
				} else {
					ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User", "Cannot process your user information. Please try again in a moment")
				}
				tx.Rollback()
				return
			}
			mlc := miauth.MiauthLoginCredential{Hash: string(hashed)}
			if err := miauth.DB.Create(&mlc).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				tx.Rollback()
				return
			}
			lc := miauth.LoginCredential{
				UserID:              user.ID,
				KindLoginCredential: miauth.MiauthLC,
				LoginCredentialID:   mlc.ID}
			if err := miauth.DB.Create(&lc).Error; err != nil {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
				tx.Rollback()
				return
			}
		}
		tx.Commit()
		c.JSON(http.StatusOK, user)
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
	var tk *jwt.Token
	if tk, err = miauth.VerifyAccessToken(queryParam.AccessToken); err != nil || !tk.Valid {
		if err == nil {
			err = errors.New("invalid token")
		}
		ErrorResponse(c, http.StatusUnauthorized, err, "Invalid token", "Invalid token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"isOk":    true,
		"payload": tk.Claims,
	})
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

	if isOk, err := miauth.VerifyRefreshToken(input.RefreshToken); err != nil || !isOk {
		if err == nil {
			err = errors.New("invalid refresh token")
		}
		ErrorResponse(c, http.StatusUnauthorized, err, "Invalid token", "Invalid token")
		return
	}

	var session miauth.Session
	if err := miauth.DB.Where(&miauth.Session{RefreshToken: input.RefreshToken}).First(&session).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, err, "no session found", "Invalid Session")
		return
	}
	var user miauth.User
	if err := miauth.DB.Where(&miauth.User{Base: miauth.Base{ID: session.UserID}}).First(&user).Error; err != nil {
		ErrorResponse(c, http.StatusUnauthorized, err, "no user found", "Invalid User")
		return
	}
	accessToken, expiresIn, err := miauth.TokenizeAccessToken(user.ID.String(), user.Email)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "Cannot tokenize access token", "Wrong tokenize process")
		return
	}
	refreshToken, err := miauth.TokenizeRefreshToken(user.ID.String(), user.Email)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "Cannot tokenize access token", "Wrong tokenize process")
		return
	}
	newSession := &miauth.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		ExpiresIn:    expiresIn,
		Scope:        nil,
	}
	if err := miauth.DB.Create(newSession).Error; err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create session", "Session was not created. Please try again")
		return
	}
	miauth.DB.Delete(&session)

	c.JSON(http.StatusOK, newSession)
}
