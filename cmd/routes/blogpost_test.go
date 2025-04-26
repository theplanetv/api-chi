package routes

import (
	"api-chi/cmd/models"
	"api-chi/cmd/services"
	"api-chi/internal/message"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func Test_BlogPostRoutes(t *testing.T) {
	r := chi.NewRouter()
	BlogPostRoutes(r)
	id := ""
	slug := ""
	service := services.AuthService{}
	token, _ := service.GenerateToken(&models.Auth{Username: "admin"})

	authCookie := &http.Cookie{
		Name:  "auth-token",
		Value: token,
	}

	t.Run("Create success", func(t *testing.T) {
		input := models.BlogPostCreated{
			Title:     "new post",
			Content:   "## Hello new post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
			Tags:      []models.BlogTag{},
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest("POST", "/blog/posts", bytes.NewBuffer(body))
		req.AddCookie(authCookie)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var response message.Response
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, message.CREATE_DATA_SUCCESS, response.Message)
		assert.NotNil(t, response.Data)

		// Convert response.Data to map to extract BlogTag
		dataMap, ok := response.Data.(map[string]any)
		if !ok {
			t.Fatalf("Expected response.Data to be a map, got %T", response.Data)
		}

		// Extract id from the BlogPost data
		id = dataMap["id"].(string)
		assert.NotEmpty(t, id)

		// Extract slug from the BlogPost data
		slug = dataMap["slug"].(string)
		assert.NotEmpty(t, slug)
	})

	t.Run("Count success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/blog/posts/count", nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var response message.Response
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, message.GET_DATA_SUCCESS, response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("Get with slug success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/blog/posts/slug/"+slug, nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var response message.Response
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, message.GET_DATA_SUCCESS, response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("GetAll success", func(t *testing.T) {
		search := ""
		limit := 10
		page := 1

		req := httptest.NewRequest("GET", fmt.Sprintf("/blog/posts?search=%s&limit=%d&page=%d", search, limit, page), nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var response message.Response
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, message.GET_DATA_SUCCESS, response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("Update success", func(t *testing.T) {
		if id == "" {
			t.Fatal("ID must be set before running Update test")
		}

		input := models.BlogPostUpdated{
			Id:        id,
			Title:     "My test post",
			Content:   "## Hello my test post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest("PATCH", "/blog/posts", bytes.NewBuffer(body))
		req.AddCookie(authCookie)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var response message.Response
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, message.UPDATE_DATA_SUCCESS, response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("Remove success", func(t *testing.T) {
		if id == "" {
			t.Fatal("ID must be set before running Remove test")
		}

		req := httptest.NewRequest("DELETE", "/blog/posts/"+id, nil)
		req.AddCookie(authCookie)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var response message.Response
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, message.REMOVE_DATA_SUCCESS, response.Message)
		assert.NotNil(t, response.Data)
	})
}
