package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmcshane/hipchat-openshift/service"
)

var (
	testVal = "abcdefg"
)

func TestPostToCredsHandler(t *testing.T) {
	jsonStr := []byte(`{"Token":"abcdefg"}`)
	req, err := http.NewRequest("POST", "/creds", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	tokenService := service.NewTokenService()
	handler := http.Handler(NewOcCredsHandler(tokenService))
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if testVal != tokenService.Token {
		t.Errorf("handler set token to unexpected value: got %v want %v",
			tokenService.Token, testVal)
	}
}
