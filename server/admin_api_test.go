package server_test

import (
	"encoding/json"
	miauthv2 "github.com/anargu/miauth"
	. "github.com/anargu/miauth/server"
	"github.com/go-playground/assert/v2"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"testing"
)



func TestBadParamsRevokeAll(t *testing.T) {
	path := "/revokeall"
	r := setupTestServer(http.MethodPost,path, RevokeAllEndpoint)

	body := &RevokeAllInputPayload{}
	buff, _ := passStructToReader(body)
	w := performRequest(r, http.MethodPost, path, buff)

	_body := w.Body
	var errorExpected ErrorResponsePayload
	if err := json.Unmarshal(_body.Bytes(), &errorExpected); err != nil {
		panic(err)
	}
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Equal(t, errorExpected.ErrorDescription, "Bad Params")

	body.UserID = "abc"
	buff, _ = passStructToReader(body)
	w = performRequest(r, http.MethodPost, path, buff)
	_body = w.Body
	var errorExpected2 ErrorResponsePayload
	if err := json.Unmarshal(_body.Bytes(), &errorExpected2); err != nil {
		panic(err)
	}
	assert.Equal(t, w.Code, http.StatusInternalServerError)
	assert.Equal(t, errorExpected2.ErrorDescription, "User ID Corrupted")
}

func TestSuccessRevokeAll(t *testing.T) {
	setupDBConfig(t)

	path := "/revokeall"
	r := setupTestServer(http.MethodPost, path, RevokeAllEndpoint)

	body := &RevokeAllInputPayload{}
	body.UserID = uuid.NewV4().String()
	buff, _ := passStructToReader(body)
	w := performRequest(r, http.MethodPost, path, buff)
	_body := w.Body

	var result = struct{
		SessionDeleted []miauthv2.Session `json:"sessions_deleted"`
	}{}
	if err := json.Unmarshal(_body.Bytes(), &result); err != nil {
		panic(err)
	}
	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, len(result.SessionDeleted), 0)
}
