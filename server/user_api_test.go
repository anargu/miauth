package server_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	miauthv2 "github.com/anargu/miauth"
	"github.com/anargu/miauth/server"
	"github.com/go-playground/assert/v2"
	"github.com/k0kubun/pp"
)

func TestMain(m *testing.M) {
	log.Println("=== setup before Tests ===")
	SetupConfig()
	miauthv2.InitConfig()
	exitVal := m.Run()
	log.Println("=== after all Tests ===")
	os.Exit(exitVal)
}

func TestBadParamsLogin(t *testing.T) {
	r := setupTestServer(http.MethodPost, "/login", server.LoginEndpoint)

	username := "abc"
	password := "1234"
	credential := server.MiauthCredential{
		Password: &password,
	}

	input := server.LoginInputPayload{Username: &username, Credentials: credential, Kind: "miauth", UserRole: "user"}
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

	/*
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
	*/
}

func createSomeUser(username string, email string, kind string, credentialValue string) error {
	role := miauthv2.Role{}
	if err := miauthv2.DB.Where(&miauthv2.Role{Name: "user"}).First(&role).Error; err != nil {
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
			UserID:              user.ID,
			LoginCredentialID:   mlc.ID,
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
			UserID:              user.ID,
			LoginCredentialID:   flc.ID,
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
			UserID:              user.ID,
			LoginCredentialID:   glc.ID,
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

	if err := createSomeUser("abcdef", "abc@abc.com", "miauth", "1234"); err != nil {
		t.Fatal(err)
	}

	r := setupTestServer(http.MethodPost, "/login", server.LoginEndpoint)

	username := "abcdef"
	//email := "abc@abc.com"
	password := "1234"
	input := server.LoginInputPayload{
		UserRole: "user",
		Kind:     "miauth",
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
		Kind:     "facebook",
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
		Kind:     "google",
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

func TestMiauthSignup(t *testing.T) {
	setupDBConfig(t)
	// if err := setupForSignup(); err != nil {
	// 	t.Fatal(err)
	// }
	r := setupTestServer(http.MethodPost, "/signup", server.SignupEndpoint)

	type ExpectedData struct {
		ShouldBeOk          bool
		MatchExpectedValues map[string]string
	}
	type TestCase struct {
		Name     string
		Input    map[string]string
		Expected ExpectedData
	}

	cases := []TestCase{
		TestCase{
			Name: "error case",
			Input: map[string]string{
				"email":    "abcdefg@abc.com",
				"username": "abc",
			},
			Expected: ExpectedData{
				ShouldBeOk:          false,
				MatchExpectedValues: map[string]string{},
			},
		},
		TestCase{
			Name: "success case",
			Input: map[string]string{
				"email":    "abcdefg@xyz.com",
				"username": "abcdefg",
			},
			Expected: ExpectedData{
				ShouldBeOk:          true,
				MatchExpectedValues: map[string]string{},
			},
		},
	}

	for i, testcase := range cases {
		email := testcase.Input["email"]
		password := "1234"
		username := testcase.Input["username"]
		input := server.SignupInputPayload{
			Kind:     "miauth",
			Username: &username,
			Email:    &email,
			UserRole: "user",
			Credentials: server.MiauthCredential{
				Password: &password,
			},
		}
		t.Run(fmt.Sprintf("case %s %v", testcase.Name, i), func(t *testing.T) {
			readio, err := passStructToReader(input)
			if err != nil {
				t.Fatal(err)
			}

			w := performRequest(r, http.MethodPost, "/signup", readio)
			var response miauthv2.User
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				log.Fatal(err)
			}

			if testcase.Expected.ShouldBeOk {
				assert.Equal(t, w.Code, http.StatusOK)
				pp.Printf("username: %v\n", response.Username)
				assert.NotEqual(t, response.Username, "")
				assert.NotEqual(t, response.Email, "")

				if err := miauthv2.DB.Preload("Credentials").Find(&response).Error; err != nil {
					log.Fatal(err)
				}
			} else {
				pp.Printf("payload: %v\n\n", w.Body.String())
				assert.Equal(t, w.Code, http.StatusBadRequest)
			}

		})
	}
}

type expectedResult struct {
	Success        bool
	ExpectedValues map[string]string
}

func TestThirdPartySignup(t *testing.T) {
	setupDBConfig(t)
	// if err := setupForSignup(); err != nil {
	// 	t.Fatal(err)
	// }
	r := setupTestServer(http.MethodPost, "/signup", server.SignupEndpoint)

	testCases := map[string]expectedResult{
		`{ 
				"kind": "miauth", "role": "user", "email": "juan@abc.com", "username": "anargu",
				"credential": {
					"password": "anargu"
				}
		}`: {
			Success:        true,
			ExpectedValues: map[string]string{},
		},
		`{ 
				"kind": "apple", "role": "user", "email": "juan0@abc.com", "username": "",
				"credential": {
					"password": "anargu"
				}
		}`: {
			Success:        false,
			ExpectedValues: map[string]string{},
		},
		`{ 
				"kind": "apple", "role": "user", "email": "juan1@abc.com", "username": "",
				"credential": {
					"account_id": "11111"
				}
		}`: {
			Success:        true,
			ExpectedValues: map[string]string{},
		},
		`{ 
				"kind": "facebook", "role": "user", "email": "juan2@abc.com", "username": "",
				"credential": {
					"account_id": "abcdef"
				}
		}`: {
			Success:        true,
			ExpectedValues: map[string]string{},
		},
		`{ 
				"kind": "google", "role": "user", "email": "juan3@abc.com", "username": "",
				"credential": {
					"account_id": "3456"
				}
		}`: {
			Success:        true,
			ExpectedValues: map[string]string{},
		},
		`{ 
				"kind": "apple", "role": "user", "email": "juan4@abc.com", "username": "anargu",
				"credential": {
					"account_id": "1234"
				}
		}`: {
			Success:        true,
			ExpectedValues: map[string]string{},
		},
		`{ 
				"kind": "miauth", "role": "user", "email": "juan5@abc.com", "username": "",
				"credential": {
					"password": "anargu"
				}
		}`: {
			Success:        false,
			ExpectedValues: map[string]string{},
		},
	}
	newUsernames := map[string]int{}
	for payload, expectedResponse := range testCases {

		t.Run(fmt.Sprintf("testCase okResponse %v", expectedResponse.Success), func(t *testing.T) {
			pp.Printf("input: %v\n", payload)
			w := performRequest(r, http.MethodPost, "/signup", strings.NewReader(payload))
			var response miauthv2.User
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				log.Fatal(err)
			}

			var input server.SignupInputPayload
			if err := json.Unmarshal([]byte(payload), &input); err != nil {
				log.Fatal(err)
			}

			if expectedResponse.Success {
				assert.Equal(t, w.Code, http.StatusOK)
				pp.Printf("payload:\nusername: %v\nemail: %v \n", response.Username, response.Email)
				if input.Kind == "miauth" {
					assert.Equal(t, response.Username, input.Username)
				} else {
					assert.NotEqual(t, response.Username, "")
					assert.MatchRegex(t, response.Username, regexp.MustCompile(`^user_(\d+)$`))
					if counter, ok := newUsernames[response.Username]; ok {
						newUsernames[response.Username] = counter + 1
					} else {
						newUsernames[response.Username] = 1
					}
				}
				assert.NotEqual(t, response.Email, "")
			} else {
				pp.Printf("error: %v\n", w.Body.String())
				assert.NotEqual(t, w.Code, http.StatusOK)
			}
		})
	}

	for u, v := range newUsernames {
		assert.Equal(t, v, 1)
		pp.Printf(">> unique username created : %v\n", u)
	}

}
