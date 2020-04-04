package server_test

import (
	"bytes"
	"encoding/json"
	miauthv2 "github.com/anargu/miauth"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"
)

func setupTestServer(method string, path string, handler gin.HandlerFunc) *gin.Engine {
	r := gin.Default()
	r.Handle(method, path, handler)
	return r
}

func setupDBConfig(t *testing.T) {
	data, err := ioutil.ReadFile("../miauth.config.v2.yml")
	if err != nil {
		t.Fatal(err)
	}
	if err := miauthv2.ReadConfig(string(data)); err != nil {
		t.Fatal(err)
	}
	miauthv2.Config.DB.Postgres = "user=miauth password=miauth DB.name=miauth port=9910 sslmode=disable"

	miauthv2.InitDB()
	if err := miauthv2.DB.DropTableIfExists(miauthv2.Tables...).Error; err != nil {
		t.Fatal(err)
	}
	miauthv2.RunMigration()
}


func SetupDatabase(t *testing.T, testHandler func()) {
	cmd := exec.Command("/bin/sh","-c", "docker-compose -f ../test-compose.yml up -d")
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1000 * time.Millisecond)

	setupDBConfig(t)

	testHandler()


	defer func() {
		cmd = exec.Command("/bin/sh","-c", "docker container stop miauth_postgres_db_test")
		if err = cmd.Run(); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command("/bin/sh","-c", "docker container rm miauth_postgres_db_test")
		if err = cmd.Run(); err != nil {
			t.Fatal(err)
		}
	}()
}

func performRequest(r http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func passStructToReader(v interface{}) (*bytes.Reader, error) {
	bodyBytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	buff := bytes.NewReader(bodyBytes)
	return buff, nil
}