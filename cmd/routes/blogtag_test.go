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

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func Test_BlogTagRoutes(t *testing.T) {
	r := chi.NewRouter()
	BlogTagRoutes(r)
	id := ""
	service := services.AuthService{}
	token, _ := service.GenerateToken(&models.Auth{Username: "admin"})

	authCookie := &http.Cookie{
		Name:  "auth-token",
		Value: token,
	}

	t.Run("Count success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/blog/tags/count", nil)
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

		req := httptest.NewRequest("GET", fmt.Sprintf("/blog/tags?search=%s&limit=%d&page=%d", search, limit, page), nil)
		res := httptest.NewRecorder()

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
		var response message.Response
		err := json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, message.GET_DATA_SUCCESS, response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("Create success", func(t *testing.T) {
		input := models.BlogTag{Name: "test tag"}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest("POST", "/blog/tags", bytes.NewBuffer(body))
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
		dataMap, ok := response.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected response.Data to be a map, got %T", response.Data)
		}

		// Extract id from the BlogTag data
		id = dataMap["id"].(string)
		assert.NotEmpty(t, id)
	})

	t.Run("Update success", func(t *testing.T) {
		if id == "" {
			t.Fatal("ID must be set before running Update test")
		}

		input := models.BlogTag{
			Id:   id,
			Name: "updated tag",
		}
		body, _ := json.Marshal(input)

		req := httptest.NewRequest("PATCH", "/blog/tags", bytes.NewBuffer(body))
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

		req := httptest.NewRequest("DELETE", "/blog/tags/"+id, nil)
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
