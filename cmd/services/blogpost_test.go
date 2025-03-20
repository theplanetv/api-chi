package services

import (
	"api-chi/cmd/models"
	"testing"
	"time"

	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
)

func Test_BlogPostService(t *testing.T) {
	id := ""
	tagService := BlogTagService{}
	tagService.Open()
	defer tagService.Close()
	postService := BlogPostService{}
	tag1 := models.BlogTag{
		Name: "website",
	}
	tag2 := models.BlogTag{
		Name: "technology",
	}
	tag3 := models.BlogTag{
		Name: "life",
	}
	tagValue1, _ := tagService.Create(&tag1)
	tagValue2, _ := tagService.Create(&tag2)
	tagValue3, _ := tagService.Create(&tag3)
	defer tagService.Remove(tagValue1.Id)
	defer tagService.Remove(tagValue2.Id)
	defer tagService.Remove(tagValue3.Id)

	t.Run("Create success all attributes", func(t *testing.T) {
		// Connect database
		err := postService.Open()
		defer postService.Close()
		assert.NoError(t, err)

		// Declare input
		tags := []models.BlogTag{tagValue1, tagValue2}
		input := models.BlogPostCreated{
			Title:     "new post",
			Content:   "## Hello new post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
			Tags:      tags,
		}

		// Create post
		value, err := postService.Create(&input)
		assert.NoError(t, err)
		assert.NotEmpty(t, value)
		assert.Equal(t, value.Title, input.Title)
		assert.Equal(t, value.Slug, slug.Make(input.Title))
		assert.WithinDuration(t, value.CreatedAt, input.CreatedAt, time.Millisecond)
		assert.WithinDuration(t, value.UpdatedAt, input.UpdatedAt, time.Millisecond)
		assert.Equal(t, value.IsDraft, input.IsDraft)
		for _, item := range value.Tags {
			assert.NotEmpty(t, item.Id)
			assert.NotEmpty(t, item.Name)
		}

		// Assign value to id
		id = value.Id
	})

	t.Run("Update success", func(t *testing.T) {
		// Connect database
		err := postService.Open()
		defer postService.Close()
		assert.NoError(t, err)

		// Declare input
		tags := []models.BlogTag{tagValue1, tagValue3}
		input := models.BlogPostUpdated{
			Id:        id,
			Title:     "My test post",
			Content:   "## Hello my test post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
			Tags:      tags,
		}

		// Update post
		value, err := postService.Update(&input)
		assert.NoError(t, err)
		assert.NotEmpty(t, value)
		assert.Equal(t, value.Title, input.Title)
		assert.Equal(t, value.Slug, slug.Make(input.Title))
		assert.WithinDuration(t, value.CreatedAt, input.CreatedAt, time.Millisecond)
		assert.WithinDuration(t, value.UpdatedAt, input.UpdatedAt, time.Millisecond)
		assert.Equal(t, value.IsDraft, input.IsDraft)
	})

	t.Run("Remove success", func(t *testing.T) {
		// Connect database
		err := postService.Open()
		defer postService.Close()
		assert.NoError(t, err)

		// Remove post
		value, err := postService.Remove(id)
		assert.NoError(t, err)
		assert.NotEmpty(t, value)
	})

	t.Run("GetAll default success", func(t *testing.T) {
		// Connect database
		err := postService.Open()
		defer postService.Close()
		assert.NoError(t, err)

		// Create data
		tagsPost1 := []models.BlogTag{tagValue1, tagValue2}
		inputPost1 := models.BlogPostCreated{
			Title:     "new post",
			Content:   "## Hello new post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
			Tags:      tagsPost1,
		}
		tagsPost2 := []models.BlogTag{tagValue2, tagValue3}
		inputPost2 := models.BlogPostCreated{
			Title:     "My test post",
			Content:   "## Hello my test post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
			Tags:      tagsPost2,
		}
		valuePost1, _ := postService.Create(&inputPost1)
		valuePost2, _ := postService.Create(&inputPost2)
		defer postService.Remove(valuePost1.Id)
		defer postService.Remove(valuePost2.Id)

		// Declare input
		search := ""
		tagsSearch := []models.BlogTag{
			tagValue1,
			tagValue2,
		}
		limit := 10
		page := 1

		// Get all database
		data, err := postService.GetAll(search, tagsSearch, limit, page)
		assert.NoError(t, err)

		assert.IsType(t, data[0], models.BlogPostWithTags{})
		count := 0
		for _, postItem := range data {
			count += 1
			assert.NotEmpty(t, postItem.Id)
			assert.NotEmpty(t, postItem.Title)
			assert.NotEmpty(t, postItem.Slug)
			assert.NotEmpty(t, postItem.CreatedAt)
			assert.NotEmpty(t, postItem.UpdatedAt)
			assert.NotEmpty(t, postItem.IsDraft)
		}
		assert.Equal(t, count, 1)
	})

	t.Run("Count success", func(t *testing.T) {
		// Connect database
		err := postService.Open()
		defer postService.Close()
		assert.NoError(t, err)

		// Create data
		tagsPost1 := []models.BlogTag{tagValue1, tagValue2}
		inputPost1 := models.BlogPostCreated{
			Title:     "new post",
			Content:   "## Hello new post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
			Tags:      tagsPost1,
		}
		tagsPost2 := []models.BlogTag{tagValue2, tagValue3}
		inputPost2 := models.BlogPostCreated{
			Title:     "My test post",
			Content:   "## Hello my test post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
			Tags:      tagsPost2,
		}
		valuePost1, _ := postService.Create(&inputPost1)
		valuePost2, _ := postService.Create(&inputPost2)
		defer postService.Remove(valuePost1.Id)
		defer postService.Remove(valuePost2.Id)

		// Declare input
		search := ""
		tagsSearch := []models.BlogTag{
			tagValue1,
			tagValue2,
		}

		// Count database
		count, err := postService.Count(search, tagsSearch)
		assert.NoError(t, err)
		assert.Greater(t, count, 0)
	})
}
