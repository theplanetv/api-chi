package services

import (
	"api-chi/cmd/models"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_BlogPostService(t *testing.T) {
	id := ""
	service := BlogPostService{}

	t.Run("Count success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		search := ""

		// Count database
		count, err := service.Count(search)
		assert.NoError(t, err)
		assert.Greater(t, count, 0)
	})

	t.Run("Count success with search", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		search := "post 1"

		// Count database
		count, err := service.Count(search)
		assert.NoError(t, err)
		assert.Greater(t, count, 0)
		assert.Less(t, count, 60)
	})

	t.Run("GetAll default success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		search := ""
		tags := ""
		limit := 10
		page := 1

		// Get all database
		data, err := service.GetAll(search, tags, limit, page)
		assert.NoError(t, err)
		assert.IsType(t, data[0], models.BlogPost{})
		count := 0
		for _, item := range data {
			count += 1
			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.Title)
			assert.NotEmpty(t, item.Slug)
			assert.NotEmpty(t, item.CreatedAt)
			assert.NotEmpty(t, item.UpdatedAt)
			assert.NotEmpty(t, item.IsDraft)
		}
		assert.Equal(t, count, 10)
	})

	t.Run("GetAll limit < 10 will return 10 success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		search := ""
		tags := ""
		limit := 9
		page := 1

		// Get all database
		data, err := service.GetAll(search, tags, limit, page)
		assert.NoError(t, err)
		assert.IsType(t, data[0], models.BlogPost{})
		count := 0
		for _, item := range data {
			count += 1
			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.Title)
			assert.NotEmpty(t, item.Slug)
			assert.NotEmpty(t, item.CreatedAt)
			assert.NotEmpty(t, item.UpdatedAt)
			assert.NotEmpty(t, item.IsDraft)
		}
		assert.Equal(t, count, 10)
	})

	t.Run("GetAll limit > 50 will return 50 success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		search := ""
		tags := ""
		limit := 51
		page := 1

		// Get all database
		data, err := service.GetAll(search, tags, limit, page)
		assert.NoError(t, err)
		assert.IsType(t, data[0], models.BlogPost{})
		count := 0
		for _, item := range data {
			count += 1
			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.Title)
			assert.NotEmpty(t, item.Slug)
			assert.NotEmpty(t, item.CreatedAt)
			assert.NotEmpty(t, item.UpdatedAt)
			assert.NotEmpty(t, item.IsDraft)
		}
		assert.Equal(t, count, 50)
	})

	t.Run("GetAll page < 1 will return 10 success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		search := ""
		tags := ""
		limit := 10
		page := 0

		// Get all database
		data, err := service.GetAll(search, tags, limit, page)
		assert.NoError(t, err)
		assert.IsType(t, data[0], models.BlogPost{})
		count := 0
		for _, item := range data {
			count += 1
			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.Title)
			assert.NotEmpty(t, item.Slug)
			assert.NotEmpty(t, item.CreatedAt)
			assert.NotEmpty(t, item.UpdatedAt)
			assert.NotEmpty(t, item.IsDraft)
		}
		assert.Equal(t, count, 10)
	})

	t.Run("GetAll with search is post 1", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		search := "post 1"
		tags := ""
		limit := 20
		page := 1

		// Get all database
		data, err := service.GetAll(search, tags, limit, page)
		assert.NoError(t, err)
		assert.IsType(t, data[0], models.BlogPost{})
		count := 0
		for _, item := range data {
			count += 1
			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.Title)
			assert.NotEmpty(t, item.Slug)
			assert.NotEmpty(t, item.CreatedAt)
			assert.NotEmpty(t, item.UpdatedAt)
			assert.NotEmpty(t, item.IsDraft)
		}
		assert.Less(t, count, 20)
	})

	t.Run("GetAll with tags is tag 1", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		search := ""
		tags := "tag 1"
		limit := 10
		page := 1

		// Get all database
		data, err := service.GetAll(search, tags, limit, page)
		assert.NoError(t, err)
		assert.IsType(t, data[0], models.BlogPost{})
		count := 0
		for _, item := range data {
			count += 1
			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.Title)
			assert.NotEmpty(t, item.Slug)
			assert.NotEmpty(t, item.CreatedAt)
			assert.NotEmpty(t, item.UpdatedAt)
			assert.NotEmpty(t, item.IsDraft)
		}
		assert.Greater(t, count, 0)
	})

	t.Run("GetAll with tags are tag 2 and tag 3", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		search := ""
		tags := "tag 2,tag 3"
		limit := 10
		page := 1

		// Get all database
		data, err := service.GetAll(search, tags, limit, page)
		assert.NoError(t, err)
		assert.IsType(t, data[0], models.BlogPost{})
		count := 0
		for _, item := range data {
			count += 1
			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.Title)
			assert.NotEmpty(t, item.Slug)
			assert.NotEmpty(t, item.CreatedAt)
			assert.NotEmpty(t, item.UpdatedAt)
			assert.NotEmpty(t, item.IsDraft)
		}
		assert.Greater(t, count, 0)
	})

	t.Run("Create success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		input := models.BlogPost{
			Title:     "new post",
			Slug:      "new_post",
			Content:   "## Hello new post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
		}

		// Create database
		value, err := service.Create(&input)
		assert.NoError(t, err)
		assert.NotEmpty(t, value)
		assert.Equal(t, value.Title, input.Title)
		assert.Equal(t, value.Slug, input.Slug)
		assert.WithinDuration(t, value.CreatedAt, input.CreatedAt, time.Millisecond)
		assert.WithinDuration(t, value.UpdatedAt, input.UpdatedAt, time.Millisecond)
		assert.Equal(t, value.IsDraft, input.IsDraft)

		// Assign value to id
		id = value.Id
	})

	t.Run("Update success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		input := models.BlogPost{
			Id:        id,
			Title:     "My test post",
			Slug:      "my_test_post",
			Content:   "## Hello my test post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
		}

		// Update database
		value, err := service.Update(&input)
		assert.NoError(t, err)
		assert.NotEmpty(t, value)
		assert.Equal(t, value.Title, input.Title)
		assert.Equal(t, value.Slug, input.Slug)
		assert.WithinDuration(t, value.CreatedAt, input.CreatedAt, time.Millisecond)
		assert.WithinDuration(t, value.UpdatedAt, input.UpdatedAt, time.Millisecond)
		assert.Equal(t, value.IsDraft, input.IsDraft)
	})

	t.Run("Remove success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Remove database
		value, err := service.Remove(id)
		assert.NoError(t, err)
		assert.NotEmpty(t, value)
	})
}
