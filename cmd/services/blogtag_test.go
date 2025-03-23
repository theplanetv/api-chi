package services

import (
	"api-chi/cmd/models"

	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BlogTagService(t *testing.T) {
	id := ""
	service := BlogTagService{}

	t.Run("Create success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		input := models.BlogTag{
			Name: "test tag",
		}

		// Create database
		value, err := service.Create(&input)
		assert.NoError(t, err)
		assert.NotEmpty(t, value)
		assert.NotEmpty(t, value.Id)
		assert.Equal(t, input.Name, value.Name)

		// Assign value to id
		id = value.Id
	})

	t.Run("Count success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		assert.NoError(t, err)
		defer service.Close()

		// Declare input
		search := ""

		// Count database
		count, err := service.Count(search)
		assert.NoError(t, err)
		assert.Equal(t, count, 1)
	})

	t.Run("GetAll success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		search := ""
		limit := 3
		page := 1

		// Get all database
		data, err := service.GetAll(search, limit, page)
		assert.NoError(t, err)
		assert.IsType(t, data[0], models.BlogTag{})
		count := 0
		for _, item := range data {
			count += 1
			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.Name)
		}
		assert.Equal(t, 1, count)
	})

	t.Run("Update success", func(t *testing.T) {
		// Connect database
		err := service.Open()
		defer service.Close()
		assert.NoError(t, err)

		// Declare input
		input := models.BlogTag{
			Id:   id,
			Name: "this is test tag",
		}

		// Update database
		value, err := service.Update(&input)
		assert.NoError(t, err)
		assert.NotEmpty(t, value)
		assert.Equal(t, input.Id, value.Id)
		assert.Equal(t, input.Name, value.Name)
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
