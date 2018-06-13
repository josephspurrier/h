package h_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josephspurrier/h"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTPResponseOK(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r, http.StatusOK, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestServeHTTPResponseIgnoreLessThan200(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r, http.StatusCreated, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestServeHTTPResponseError(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r, http.StatusInternalServerError, errors.New("error happened"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error happened")
}
