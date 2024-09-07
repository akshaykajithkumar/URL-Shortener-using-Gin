package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/shorten", ShortenURL)
	router.GET("/:short_url", RetrieveURL)
	return router
}

func TestShortenURL(t *testing.T) {
	router := SetupRouter()

	t.Run("valid URL", func(t *testing.T) {

		requestBody := `{"url":"https://example.com"}`
		req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "short_url")
	})

	t.Run("invalid URL (no https)", func(t *testing.T) {

		requestBody := `{"url":"http://example.com"}`
		req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Please provide valid url")
	})
}

func TestRetrieveURL(t *testing.T) {
	router := SetupRouter()

	requestBody := `{"url":"https://example.com"}`
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	shortURL := ExtractShortURL(w.Body.String())

	t.Run("valid short URL", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/"+shortURL, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Location"))
	})

	t.Run("invalid short URL", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/invalidURL", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "Short URL not found")
	})
}

func ExtractShortURL(responseBody string) string {
	parts := strings.Split(responseBody, "/")
	if len(parts) > 1 {
		return strings.Trim(parts[len(parts)-1], "\"}")
	}
	return ""
}
