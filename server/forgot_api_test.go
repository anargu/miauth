package server_test

import (
	"encoding/json"
	"github.com/anargu/miauth/server"
	"github.com/go-playground/assert/v2"
	"net/http"
	"testing"
)

func TestBadParamsForgotRequest(t *testing.T) {

	r := setupTestServer(http.MethodPost, "/forgot", server.ForgotRequestEndpoint)

	input := server.ForgotRequestInputPayload{}
	readio, err := passStructToReader(input)
	if err != nil {
		t.Fatal(err)
	}
	w := performRequest(r, http.MethodPost, "/forgot", readio)
	var errorResponse server.ErrorResponsePayload
	if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, w.Code,  http.StatusBadRequest)
	assert.Equal(t, errorResponse.ErrorDescription, "Bad Params")
}

//func TestSuccessForgotRequest(t *testing.T) {
//	setupDBConfig(t)
//
//	username := "toby"
//	email := "toby@ebay.com"
//	if err := createSomeUser(username, email, "ybot"); err != nil {
//		t.Fatal(err)
//	}
//	r := setupTestServer(http.MethodPost, "/forgot", server.ForgotRequestEndpoint)
//
//	input := server.ForgotRequestInputPayload{
//		Email: email,
//		//Username: username,
//	}
//	readio, err := passStructToReader(input)
//	if err != nil {
//		t.Fatal(err)
//	}
//	w := performRequest(r, http.MethodPost, "/forgot", readio)
//	var response struct{
//		Message string
//	}
//	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
//		t.Fatal(err)
//	}
//
//	assert.Equal(t, w.Code, http.StatusOK)
//	assert.Equal(t, response.Message, "Email sent to the user with the instructions to reset password.")
//}