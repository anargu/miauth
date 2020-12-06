package miauth_test

import (
	"io/ioutil"
	"testing"

	miauthv2 "github.com/anargu/miauth"
	"github.com/go-playground/assert/v2"
)

func init() {
	LoadConfig()
}

func LoadConfig() {
	data, err := ioutil.ReadFile("./miauth.config.v2.yml")
	if err != nil {
		panic(err)
	}
	if err := miauthv2.ReadConfig(string(data)); err != nil {
		panic(err)
	}
}

func TestLoadingConfig(t *testing.T) {

	if miauthv2.Config == nil {
		t.Fatal("miauthv2.Config not loaded")
	}

	assert.Equal(t, miauthv2.Config.ResetPassword.ExpiresIn, "600")
}
