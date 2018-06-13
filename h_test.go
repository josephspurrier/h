package h_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/josephspurrier/h"
	"github.com/stretchr/testify/assert"
)

func TestFResponseOK(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", h.F(func(w http.ResponseWriter, r *http.Request) (int, error) {
		return http.StatusOK, nil
	}))
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "")
}

func TestFResponseInternalServerError(t *testing.T) {
	mux := http.NewServeMux()
	mux.Handle("/", h.F(func(w http.ResponseWriter, r *http.Request) (int, error) {
		return http.StatusInternalServerError, errors.New("error happened")
	}))
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, r)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error happened")
}

func TestServeHTTPResponseOK(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r, http.StatusOK, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "")
}

func TestServeHTTPResponseIgnoreLessThan200(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r, http.StatusCreated, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "")
}

func TestServeHTTPResponseErrorExist(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r, http.StatusInternalServerError, errors.New("error happened"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error happened")
}

func TestServeHTTPResponseErrorEmpty(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r, http.StatusInternalServerError, nil)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), "\n")
}
