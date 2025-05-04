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
	err := tagService.Open()
	assert.NoError(t, err)
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
	tagValue1, err := tagService.Create(&tag1)
	assert.NoError(t, err)
	tagValue2, err := tagService.Create(&tag2)
	assert.NoError(t, err)
	tagValue3, err := tagService.Create(&tag3)
	assert.NoError(t, err)
	defer func() {
		_, err = tagService.Remove(tagValue1.Id)
		assert.NoError(t, err)
		_, err = tagService.Remove(tagValue2.Id)
		assert.NoError(t, err)
		_, err = tagService.Remove(tagValue3.Id)
		assert.NoError(t, err)
	}()

	t.Run("Create success", func(t *testing.T) {
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
		assert.IsType(t, value, models.BlogPostContentWithTags{})
		assert.Equal(t, value.Title, input.Title)
		assert.Equal(t, value.Slug, slug.Make(input.Title))
		assert.Equal(t, value.Content, input.Content)
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
		assert.IsType(t, value, models.BlogPostContentWithTags{})
		assert.Equal(t, value.Title, input.Title)
		assert.Equal(t, value.Slug, slug.Make(input.Title))
		assert.Equal(t, value.Content, input.Content)
		assert.WithinDuration(t, value.CreatedAt, input.CreatedAt, time.Millisecond)
		assert.WithinDuration(t, value.UpdatedAt, input.UpdatedAt, time.Millisecond)
		assert.Equal(t, value.IsDraft, input.IsDraft)
		for _, tag := range value.Tags {
			assert.NotEmpty(t, tag.Id)
			assert.NotEmpty(t, tag.Name)
		}
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

	t.Run("Get with slug success", func(t *testing.T) {
		// Connect database
		err := postService.Open()
		defer postService.Close()
		assert.NoError(t, err)

		// Create data
		tagsPost := []models.BlogTag{tagValue1, tagValue2}
		inputPost := models.BlogPostCreated{
			Title:     "new post",
			Content:   "## Hello new post!",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDraft:   true,
			Tags:      tagsPost,
		}
		valuePost, _ := postService.Create(&inputPost)
		defer func() {
			_, err = postService.Remove(valuePost.Id)
			assert.NoError(t, err)
		}()

		// Get all database
		data, err := postService.GetWithSlug(valuePost.Slug)
		assert.NoError(t, err)

		assert.IsType(t, data, models.BlogPostContentWithTags{})
		assert.NotEmpty(t, data.Id)
		assert.NotEmpty(t, data.Title)
		assert.NotEmpty(t, data.Slug)
		assert.NotEmpty(t, data.CreatedAt)
		assert.NotEmpty(t, data.UpdatedAt)
		assert.NotEmpty(t, data.IsDraft)
		for _, tag := range data.Tags {
			assert.NotEmpty(t, tag.Id)
			assert.NotEmpty(t, tag.Name)
		}
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
		defer func() {
			_, err = postService.Remove(valuePost1.Id)
			assert.NoError(t, err)
			_, err = postService.Remove(valuePost2.Id)
			assert.NoError(t, err)
		}()

		// Declare input
		search := ""
		tagsSearch := []models.BlogTag{}
		limit := 10
		page := 1

		// Get all database
		data, err := postService.GetAll(search, tagsSearch, limit, page)
		assert.NoError(t, err)

		assert.IsType(t, data[0], models.BlogPostWithTags{})
		count := 0
		for _, post := range data {
			count += 1
			assert.NotEmpty(t, post.Id)
			assert.NotEmpty(t, post.Title)
			assert.NotEmpty(t, post.Slug)
			assert.NotEmpty(t, post.CreatedAt)
			assert.NotEmpty(t, post.UpdatedAt)
			assert.NotEmpty(t, post.IsDraft)
			for _, tag := range post.Tags {
				assert.NotEmpty(t, tag.Id)
				assert.NotEmpty(t, tag.Name)
			}
		}
		assert.Equal(t, count, 2)
	})

	t.Run("GetAll success with search", func(t *testing.T) {
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
		defer func() {
			_, err = postService.Remove(valuePost1.Id)
			assert.NoError(t, err)
			_, err = postService.Remove(valuePost2.Id)
			assert.NoError(t, err)
		}()

		// Declare input
		search := "TEST"
		tagsSearch := []models.BlogTag{}
		limit := 10
		page := 1

		// Get all database
		data, err := postService.GetAll(search, tagsSearch, limit, page)
		assert.NoError(t, err)

		assert.IsType(t, data[0], models.BlogPostWithTags{})
		count := 0
		for _, post := range data {
			count += 1
			assert.NotEmpty(t, post.Id)
			assert.NotEmpty(t, post.Title)
			assert.NotEmpty(t, post.Slug)
			assert.NotEmpty(t, post.CreatedAt)
			assert.NotEmpty(t, post.UpdatedAt)
			assert.NotEmpty(t, post.IsDraft)
			for _, tag := range post.Tags {
				assert.NotEmpty(t, tag.Id)
				assert.NotEmpty(t, tag.Name)
			}
		}
		assert.Equal(t, count, 1)
	})

	t.Run("GetAll success with tags", func(t *testing.T) {
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
		defer func() {
			_, err = postService.Remove(valuePost1.Id)
			assert.NoError(t, err)
			_, err = postService.Remove(valuePost2.Id)
			assert.NoError(t, err)
		}()

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
		for _, post := range data {
			count += 1
			assert.NotEmpty(t, post.Id)
			assert.NotEmpty(t, post.Title)
			assert.NotEmpty(t, post.Slug)
			assert.NotEmpty(t, post.CreatedAt)
			assert.NotEmpty(t, post.UpdatedAt)
			assert.NotEmpty(t, post.IsDraft)
			for _, tag := range post.Tags {
				assert.NotEmpty(t, tag.Id)
				assert.NotEmpty(t, tag.Name)
			}
		}
		assert.Equal(t, count, 1)
	})

	t.Run("GetAllWithContent default success", func(t *testing.T) {
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
		defer func() {
			_, err = postService.Remove(valuePost1.Id)
			assert.NoError(t, err)
			_, err = postService.Remove(valuePost2.Id)
			assert.NoError(t, err)
		}()

		// Declare input
		search := ""
		tagsSearch := []models.BlogTag{}
		limit := 10
		page := 1

		// Get all database
		data, err := postService.GetAllWithContent(search, tagsSearch, limit, page)
		assert.NoError(t, err)

		assert.IsType(t, data[0], models.BlogPostContentWithTags{})
		count := 0
		for _, post := range data {
			count += 1
			assert.NotEmpty(t, post.Id)
			assert.NotEmpty(t, post.Title)
			assert.NotEmpty(t, post.Slug)
			assert.NotEmpty(t, post.Content)
			assert.NotEmpty(t, post.CreatedAt)
			assert.NotEmpty(t, post.UpdatedAt)
			assert.NotEmpty(t, post.IsDraft)
			for _, tag := range post.Tags {
				assert.NotEmpty(t, tag.Id)
				assert.NotEmpty(t, tag.Name)
			}
		}
		assert.Equal(t, count, 2)
	})

	t.Run("GetAll success with search", func(t *testing.T) {
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
		defer func() {
			_, err = postService.Remove(valuePost1.Id)
			assert.NoError(t, err)
			_, err = postService.Remove(valuePost2.Id)
			assert.NoError(t, err)
		}()

		// Declare input
		search := "TEST"
		tagsSearch := []models.BlogTag{}
		limit := 10
		page := 1

		// Get all database
		data, err := postService.GetAllWithContent(search, tagsSearch, limit, page)
		assert.NoError(t, err)

		assert.IsType(t, data[0], models.BlogPostContentWithTags{})
		count := 0
		for _, post := range data {
			count += 1
			assert.NotEmpty(t, post.Id)
			assert.NotEmpty(t, post.Title)
			assert.NotEmpty(t, post.Slug)
			assert.NotEmpty(t, post.Content)
			assert.NotEmpty(t, post.CreatedAt)
			assert.NotEmpty(t, post.UpdatedAt)
			assert.NotEmpty(t, post.IsDraft)
			for _, tag := range post.Tags {
				assert.NotEmpty(t, tag.Id)
				assert.NotEmpty(t, tag.Name)
			}
		}
		assert.Equal(t, count, 1)
	})

	t.Run("GetAll success with tags", func(t *testing.T) {
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
		defer func() {
			_, err = postService.Remove(valuePost1.Id)
			assert.NoError(t, err)
			_, err = postService.Remove(valuePost2.Id)
			assert.NoError(t, err)
		}()

		// Declare input
		search := ""
		tagsSearch := []models.BlogTag{
			tagValue1,
			tagValue2,
		}
		limit := 10
		page := 1

		// Get all database
		data, err := postService.GetAllWithContent(search, tagsSearch, limit, page)
		assert.NoError(t, err)

		assert.IsType(t, data[0], models.BlogPostContentWithTags{})
		count := 0
		for _, post := range data {
			count += 1
			assert.NotEmpty(t, post.Id)
			assert.NotEmpty(t, post.Title)
			assert.NotEmpty(t, post.Slug)
			assert.NotEmpty(t, post.Content)
			assert.NotEmpty(t, post.CreatedAt)
			assert.NotEmpty(t, post.UpdatedAt)
			assert.NotEmpty(t, post.IsDraft)
			for _, tag := range post.Tags {
				assert.NotEmpty(t, tag.Id)
				assert.NotEmpty(t, tag.Name)
			}
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
		defer func() {
			_, err = postService.Remove(valuePost1.Id)
			assert.NoError(t, err)
			_, err = postService.Remove(valuePost2.Id)
			assert.NoError(t, err)
		}()

		// Declare input
		search := ""
		tagsSearch := []models.BlogTag{}

		// Count database
		count, err := postService.Count(search, tagsSearch)
		assert.NoError(t, err)
		assert.Equal(t, count, 2)
	})

	t.Run("Count success with search", func(t *testing.T) {
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
		defer func() {
			_, err = postService.Remove(valuePost1.Id)
			assert.NoError(t, err)
			_, err = postService.Remove(valuePost2.Id)
			assert.NoError(t, err)
		}()

		// Declare input
		search := "TEST"
		tagsSearch := []models.BlogTag{}

		// Count database
		count, err := postService.Count(search, tagsSearch)
		assert.NoError(t, err)
		assert.Equal(t, count, 1)
	})

	t.Run("Count success with tags", func(t *testing.T) {
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
		defer func() {
			_, err = postService.Remove(valuePost1.Id)
			assert.NoError(t, err)
			_, err = postService.Remove(valuePost2.Id)
			assert.NoError(t, err)
		}()

		// Declare input
		search := ""
		tagsSearch := []models.BlogTag{
			tagValue1,
		}

		// Count database
		count, err := postService.Count(search, tagsSearch)
		assert.NoError(t, err)
		assert.Equal(t, count, 1)
	})
}
