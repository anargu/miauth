package server_test

import (
	"encoding/json"
	"errors"
	miauthv2 "github.com/anargu/miauth"
	"github.com/anargu/miauth/server"
	"github.com/go-playground/assert/v2"
	"github.com/k0kubun/pp"
	"log"
	"net/http"
	"testing"
)

func TestBadParamsLogin(t *testing.T) {
	r := setupTestServer(http.MethodPost, "/login", server.LoginEndpoint)

	input := server.LoginInputPayload{}
	readio, err := passStructToReader(input)
	if err != nil {
		log.Fatal(err)
	}
	w := performRequest(r, http.MethodPost, "/login", readio)
	errorResponse := server.ErrorResponsePayload{}
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Equal(t, errorResponse.ErrorDescription, "Bad Params")

	// mising required role param should throw BadRequest Error
	input = server.LoginInputPayload{
		Kind: "miauth",
		//UserRole: "",
	}
	readio, err = passStructToReader(input)
	if err != nil {
		log.Fatal(err)
	}
	w = performRequest(r, http.MethodPost, "/login", readio)
	errorResponse = server.ErrorResponsePayload{}
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Equal(t, errorResponse.ErrorDescription, "Bad Params")


	// mising required kind param should throw BadRequest Error
	input = server.LoginInputPayload{
		//Kind: "miauth",
		UserRole: "user",
	}
	readio, err = passStructToReader(input)
	if err != nil {
		log.Fatal(err)
	}
	w = performRequest(r, http.MethodPost, "/login", readio)
	errorResponse = server.ErrorResponsePayload{}
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Equal(t, errorResponse.ErrorDescription, "Bad Params")
}

func createSomeUser(username string, email string, kind string, credentialValue string) error {
	role := miauthv2.Role{Name: "user"}
	if err := miauthv2.DB.Create(&role).Error; err != nil {
		return err
	}
	user := miauthv2.User{Email: email, Username: username, Role: role}
	if err := miauthv2.DB.Create(&user).Error; err != nil {
		return err
	}

	switch kind {
	case "miauth":
		hash, err := miauthv2.HashPassword(credentialValue)
		if err != nil {
			return err
		}
		mlc := miauthv2.MiauthLoginCredential{Hash: *hash}
		if err := miauthv2.DB.Create(&mlc).Error; err != nil {
			return err
		}
		lc := miauthv2.LoginCredential{
			UserID: user.ID,
			LoginCredentialID: mlc.ID,
			KindLoginCredential: miauthv2.MiauthLC,
		}
		if err := miauthv2.DB.Create(&lc).Error; err != nil {
			return err
		}
		return nil
	case "facebook":
		flc := miauthv2.FacebookLoginCredential{AccountID: credentialValue}
		if err := miauthv2.DB.Create(&flc).Error; err != nil {
			return err
		}
		lc := miauthv2.LoginCredential{
			UserID: user.ID,
			LoginCredentialID: flc.ID,
			KindLoginCredential: miauthv2.FacebookLC,
		}
		if err := miauthv2.DB.Create(&lc).Error; err != nil {
			return err
		}
		return nil
	case "google":
		glc := miauthv2.GoogleLoginCredential{AccountID: credentialValue}
		if err := miauthv2.DB.Create(&glc).Error; err != nil {
			return err
		}
		lc := miauthv2.LoginCredential{
			UserID: user.ID,
			LoginCredentialID: glc.ID,
			KindLoginCredential: miauthv2.GoogleLC,
		}
		if err := miauthv2.DB.Create(&lc).Error; err != nil {
			return err
		}
		return nil
	default:
		return errors.New("Incorrect Kind of Credential selected")
	}
}

func TestSuccessMiauthLogin(t *testing.T) {
	setupDBConfig(t)

	if err := createSomeUser("abc", "abc@abc.com", "miauth", "1234"); err != nil {
		t.Fatal(err)
	}

	r := setupTestServer(http.MethodPost, "/login", server.LoginEndpoint)

	username := "abc"
	//email := "abc@abc.com"
	password := "1234"
	input := server.LoginInputPayload{
		UserRole: "user",
		Kind: "miauth",
		Username: &username,
		//Email: &email,
		Credentials: server.MiauthCredential{
			Password: &password,
		},
	}
	readio, err := passStructToReader(input)
	if err != nil {
		log.Fatal(err)
	}
	w := performRequest(r, http.MethodPost, "/login", readio)
	var session miauthv2.Session
	if err := json.Unmarshal(w.Body.Bytes(), &session); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, w.Code, http.StatusOK)
	assert.NotEqual(t, session.AccessToken, "")
	assert.NotEqual(t, session.RefreshToken, "")
}

func TestSuccessFBLogin(t *testing.T) {
	setupDBConfig(t)

	if err := createSomeUser("abc", "abc@abc.com", "facebook", "1234"); err != nil {
		t.Fatal(err)
	}

	r := setupTestServer(http.MethodPost, "/login", server.LoginEndpoint)

	username := "abc"
	//email := "abc@abc.com"
	facebookID := "1234"
	input := server.LoginInputPayload{
		UserRole: "user",
		Kind: "facebook",
		Username: &username,
		//Email: &email,
		Credentials: server.ThidPartyCredential{
			ThirdPartyAccountID: &facebookID,
		},
	}
	readio, err := passStructToReader(input)
	if err != nil {
		log.Fatal(err)
	}
	w := performRequest(r, http.MethodPost, "/login", readio)
	var session miauthv2.Session
	if err := json.Unmarshal(w.Body.Bytes(), &session); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, w.Code, http.StatusOK)
	assert.NotEqual(t, session.AccessToken, "")
	assert.NotEqual(t, session.RefreshToken, "")
}


func TestSuccessGoogleLogin(t *testing.T) {
	setupDBConfig(t)

	if err := createSomeUser("abc", "abc@abc.com", "google", "1234"); err != nil {
		t.Fatal(err)
	}

	r := setupTestServer(http.MethodPost, "/login", server.LoginEndpoint)

	username := "abc"
	//email := "abc@abc.com"
	googleID := "1234"
	input := server.LoginInputPayload{
		UserRole: "user",
		Kind: "google",
		Username: &username,
		//Email: &email,
		Credentials: server.ThidPartyCredential{
			ThirdPartyAccountID: &googleID,
		},
	}
	readio, err := passStructToReader(input)
	if err != nil {
		log.Fatal(err)
	}
	w := performRequest(r, http.MethodPost, "/login", readio)
	var session miauthv2.Session
	if err := json.Unmarshal(w.Body.Bytes(), &session); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, w.Code, http.StatusOK)
	assert.NotEqual(t, session.AccessToken, "")
	assert.NotEqual(t, session.RefreshToken, "")
}

func setupForSignup() error {
	role := miauthv2.Role{Name: "user"}
	if err := miauthv2.DB.Create(&role).Error; err != nil {
		return err
	}
	return nil
}

func TestSuccessSignup(t *testing.T) {
	setupDBConfig(t)
	if err := setupForSignup(); err != nil {
		t.Fatal(err)
	}
 	r := setupTestServer(http.MethodPost, "/signup", server.SignupEndpoint)

	username := "abc"
	email := "abc@abc.com"
	password := "1234"
	input := server.SignupInputPayload{
		Kind: "miauth",
		Username: &username,
		Email: &email,
		UserRole: "user",
		Credentials: server.MiauthCredential{
			Password: &password,
		},
	}
	readio, err := passStructToReader(input)
	if err != nil {
		t.Fatal(err)
	}
	w := performRequest(r, http.MethodPost, "/signup", readio)
	var response struct{
		Data miauthv2.User `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, w.Code, http.StatusOK)
	assert.NotEqual(t, response.Data.Username, "")
	assert.NotEqual(t, response.Data.Email, "")

	if err := miauthv2.DB.Preload("Credentials").Find(&response.Data).Error; err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, len(response.Data.Credentials), 1)
	lc := response.Data.Credentials[0]
	assert.Equal(t, lc.KindLoginCredential, miauthv2.MiauthLC)
	_,_ = pp.Printf("LoginCredentialID :: %s\n", lc.LoginCredentialID.String())
	assert.NotEqual(t, lc.LoginCredentialID.String(), "")
}
