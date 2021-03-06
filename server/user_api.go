package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/anargu/miauth"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
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
	Username    *string     `json:"username" binding:"miauth_username"`
	// Email       *string     `json:"email" binding:"miauth_email"`
	Credentials interface{} `json:"credential" binding:"required"`
}

func LoginEndpoint(c *gin.Context) {
	var credentials json.RawMessage
	input := LoginInputPayload{
		Credentials: &credentials,
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		if strings.Contains(err.Error(), "username")  {
			ErrorResponse(
				c,
				http.StatusBadRequest,
				err,
				"Bad Params",
				miauth.Config.FieldValidations.Username.InvalidPatternErrorMessage)
		} else if strings.Contains(err.Error(), "email")  {
			ErrorResponse(
				c,
				http.StatusBadRequest,
				err,
				"Bad Params",
				miauth.Config.FieldValidations.Email.InvalidPatternErrorMessage)
		} else {
			ErrorResponse(c, http.StatusBadRequest, err, "Bad Params", err.Error())
		}
		return
	}

	var role miauth.Role
	if err := miauth.DB.Where(&miauth.Role{Name: input.UserRole}).First(&role).Error; err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Role does not match", "User was binding with unknown role")
		return
	}

	switch input.Kind {
	case "apple":
		handleThirdPartyAction(
			c,
			credentials,
			"Apple ID Account not found",
			func(thirdPartyAccountID *string) (*uuid.UUID, error) {
				var alc miauth.AppleLoginCredential
				if err := miauth.DB.Where(miauth.AppleLoginCredential{AccountID: *thirdPartyAccountID}).First(&alc).Error; err != nil {
					return nil, err
				}
				return &alc.ID, nil
			},
		)
	case "facebook":
		handleThirdPartyAction(
			c,
			credentials,
			"Facebook ID Account not found",
			func(thirdPartyAccountID *string) (*uuid.UUID, error) {
				var flc miauth.FacebookLoginCredential
				if err := miauth.DB.Where(miauth.FacebookLoginCredential{AccountID: *thirdPartyAccountID}).First(&flc).Error; err != nil {
					return nil, err
				}
				return &flc.ID, nil
			},
		)
	case "google":
		handleThirdPartyAction(
			c,
			credentials,
			"Google ID Account not found",
			func(thirdPartyAccountID *string) (*uuid.UUID, error) {
				var glc miauth.GoogleLoginCredential
				if err := miauth.DB.Where(miauth.GoogleLoginCredential{AccountID: *thirdPartyAccountID}).First(&glc).Error; err != nil {
					return nil, err
				}

				return &glc.ID, nil
			},
		)
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

func handleThirdPartyAction(
	c *gin.Context,
	credentials json.RawMessage,
	notFoundMessage string,
	getLoginCredentialID func(thirdPartyAccountID *string) (*uuid.UUID, error)) {

	var err error
	var thirdPartyCredential ThidPartyCredential

	if err := json.Unmarshal(credentials, &thirdPartyCredential); err != nil {
		ErrorResponse(c, http.StatusInternalServerError, err, "Cannot process input data", "Values sent are invalid")
		return
	}
	if thirdPartyCredential.ThirdPartyAccountID == nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Missed Params", "Missed parameters")
		return
	}
	tlcUUID, err := getLoginCredentialID(thirdPartyCredential.ThirdPartyAccountID)
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, notFoundMessage, notFoundMessage)
		return
	}
	var lc miauth.LoginCredential
	if err := miauth.DB.Where(miauth.LoginCredential{LoginCredentialID: *tlcUUID}).First(&lc).Error; err != nil {
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
}

type SignupInputPayload struct {
	Kind string `json:"kind" binding:"required"`

	UserRole    string      `json:"role" binding:"required"`
	Username    *string     `json:"username" binding:"miauth_username"`
	Email       *string     `json:"email" binding:"required,email,miauth_email"`
	Credentials interface{} `json:"credential" binding:"required"`
}

func handleThirdParthSignUp(
	c *gin.Context,
	role miauth.Role,
	input SignupInputPayload,
	thirdPartyCredential ThidPartyCredential,
	kindLoginCredential int,
) *miauth.User {

	tx := miauth.DB.Begin()

	usernameGenerated, err := miauth.GenerateGenericUsername()
	if err != nil {
		tx.Rollback()
	}
	var user miauth.User
	{
		user = miauth.User{Username: *usernameGenerated, Email: *input.Email, Role: role}
		if err := miauth.DB.Create(&user).Error; err != nil {
			duplicateUsernameErrorMessage := "It seems that username have been taken."
			if _err := miauth.ValidateDuplicateErrorInField(err, "username", &duplicateUsernameErrorMessage); _err != nil {
				SendError(c, http.StatusBadRequest, _err)
			} else if _err = miauth.ValidateDuplicateErrorInField(err, "email", nil); _err != nil {
				SendError(c, http.StatusBadRequest, _err)
			} else {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User", "Cannot process your user information. Please try again in a moment")
			}
			tx.Rollback()
			return nil
		}

		lc := miauth.LoginCredential{
			UserID:              user.ID,
			KindLoginCredential: miauth.AppleLC,
		}
		thirdPartyName := ""
		var errorAtCreateThirdPartyLC error
		switch kindLoginCredential {
		case miauth.AppleLC:
			thirdPartyName = "Apple"
			alc := miauth.AppleLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}
			errorAtCreateThirdPartyLC = miauth.DB.Create(&alc).Error
			if errorAtCreateThirdPartyLC == nil {
				lc.LoginCredentialID = alc.ID
			}
			break
		case miauth.GoogleLC:
			thirdPartyName = "Google"
			glc := miauth.GoogleLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}
			errorAtCreateThirdPartyLC = miauth.DB.Create(&glc).Error
			if errorAtCreateThirdPartyLC == nil {
				lc.LoginCredentialID = glc.ID
			}
			break
		case miauth.FacebookLC:
			thirdPartyName = "Facebook"
			flc := miauth.FacebookLoginCredential{AccountID: *thirdPartyCredential.ThirdPartyAccountID}
			errorAtCreateThirdPartyLC = miauth.DB.Create(&flc).Error
			if errorAtCreateThirdPartyLC == nil {
				lc.LoginCredentialID = flc.ID
			}
			break
		}

		duplicateThirdPartyAccountErrorMessage := fmt.Sprintf("It seems that %s Account have been already signed up.", thirdPartyName)

		if errorAtCreateThirdPartyLC != nil {
			if _err := miauth.ValidateDuplicateErrorInField(err, "account_id", &duplicateThirdPartyAccountErrorMessage); _err != nil {
				SendError(c, http.StatusBadRequest, _err)
			} else {
				ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
			}
			tx.Rollback()
			return nil
		}

		if err := miauth.DB.Create(&lc).Error; err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Cannot create User Credentials", "Cannot process your user credentials. Please try a different password or try again in a moment")
			tx.Rollback()
			return nil
		}
	}
	tx.Commit()

	return &user
}

func SignupEndpoint(c *gin.Context) {
	var credentials json.RawMessage
	input := SignupInputPayload{
		Credentials: &credentials,
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		if strings.Contains(err.Error(), "username")  {
			ErrorResponse(
				c,
				http.StatusBadRequest,
				err,
				"Bad Params",
				miauth.Config.FieldValidations.Username.InvalidPatternErrorMessage)
		} else if strings.Contains(err.Error(), "email")  {
			ErrorResponse(
				c,
				http.StatusBadRequest,
				err,
				"Bad Params",
				miauth.Config.FieldValidations.Email.InvalidPatternErrorMessage)
		} else {
			ErrorResponse(c, http.StatusBadRequest, err, "Bad Params", err.Error())
		}
		return
	}
	var role miauth.Role
	if err := miauth.DB.Where(&miauth.Role{Name: input.UserRole}).First(&role).Error; err != nil {
		ErrorResponse(c, http.StatusBadRequest, err, "Role does not match", "User was binding with unknown role")
		return
	}
	switch input.Kind {
	case "apple":
		var thirdPartyCredential ThidPartyCredential
		if err := json.Unmarshal(credentials, &thirdPartyCredential); err != nil {
			ErrorResponse(c, http.StatusInternalServerError, err, "Cannot process input data", "Values sent are invalid")
			return
		}
		if thirdPartyCredential.ThirdPartyAccountID == nil {
			ErrorResponse(c, http.StatusBadRequest, err, "Missed Params", "Missed parameters")
			return
		}

		user := handleThirdParthSignUp(
			c, role, input, thirdPartyCredential, miauth.AppleLC,
		)
		if user != nil {
			c.JSON(http.StatusOK, *user)
		}
		return
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
		user := handleThirdParthSignUp(
			c, role, input, thirdPartyCredential, miauth.FacebookLC,
		)
		if user != nil {
			c.JSON(http.StatusOK, *user)
		}
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
		user := handleThirdParthSignUp(
			c, role, input, thirdPartyCredential, miauth.GoogleLC,
		)
		if user != nil {
			c.JSON(http.StatusOK, *user)
		}
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
