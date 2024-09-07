package main

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// In-memory storage to map short URLs to original URLs
var urlStore = make(map[string]string)
var mu sync.Mutex

const baseURL = "http://localhost:8080/"

// URL validation function (ensures URL starts with "https://")
func isValidURL(inputURL string) bool {
	if !strings.HasPrefix(inputURL, "https://") {
		return false
	}
	if strings.Contains(inputURL, " ") || len(inputURL) <= len("https://") || !strings.Contains(inputURL[len("https://"):], ".") {
		return false
	}
	return true
}

// Simple hash function to generate a shortened string
func shortenString(input string, length int) string {
	if length >= len(input) {
		return input
	}

	// This will store our "hash" value to generate the shortened string
	var hashValue int
	for i, char := range input {
		// Simple hash function: multiply hash value by a prime and add the char value
		hashValue = (hashValue*31 + int(char)) % 1000000007
		// Adding positional influence by mixing character index
		hashValue = hashValue ^ (int(char) * i)
	}

	// Convert the hash value to a string
	result := make([]byte, length)
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Filling the result array using the hash value and modulo operation to pick characters
	for i := 0; i < length; i++ {
		index := (hashValue + i) % len(chars) // Ensure within bounds of chars set
		result[i] = chars[index]
	}

	return string(result)
}

// Handler to shorten the URL
func ShortenURL(c *gin.Context) {
	var requestBody struct {
		URL string `json:"url"`
	}

	// Parsing request body
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validating the URL
	if !isValidURL(requestBody.URL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide valid url"})
		return
	}

	// Generating a short URL using the hash function
	shortURL := shortenString(requestBody.URL, 5)

	// Storing the original URL against the short URL
	mu.Lock()
	urlStore[shortURL] = requestBody.URL
	mu.Unlock()

	// Responding with the shortened URL
	shortenedURL := baseURL + shortURL
	c.JSON(http.StatusOK, gin.H{"short_url": shortenedURL})
}

// Handler to retrieve the original URL and redirect
func RetrieveURL(c *gin.Context) {
	// Getting the short URL from the path
	shortURL := c.Param("short_url")

	// Lookup the original URL
	mu.Lock()
	originalURL, exists := urlStore[shortURL]
	mu.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	// Redirecting to the original URL
	c.Redirect(http.StatusFound, originalURL)
}

func main() {

	router := gin.Default()

	router.POST("/shorten", ShortenURL)

	router.GET("/:short_url", RetrieveURL)

	router.Run(":8080")
}
