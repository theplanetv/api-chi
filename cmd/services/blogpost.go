package services

import (
	"api-chi/cmd/config"
	"api-chi/cmd/models"
	"fmt"
	"strings"

	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5"
)

type BlogPostService struct {
	Conn DatabaseService
}

func (s *BlogPostService) Open() error {
	return s.Conn.Open()
}

func (s *BlogPostService) Close() {
	s.Conn.Close()
}

func (s *BlogPostService) Count(search string, tags []models.BlogTag) (int, error) {
	// Base SQL query
	sql := "SELECT COUNT(blog_post.id) FROM blog_post "
	args := pgx.NamedArgs{
		"search": search,
	}

	// Add JOINs and tag filter only if tags are provided
	if len(tags) > 0 {
		sql += `
			INNER JOIN blog_post_tag ON blog_post_tag.post_id = blog_post.id
			INNER JOIN blog_tag ON blog_post_tag.tag_id = blog_tag.id
			WHERE blog_post.title ILIKE '%' || @search || '%'
		`

		// Add tag names to the query
		tagNames := make([]string, len(tags))
		for i, tag := range tags {
			paramName := fmt.Sprintf("tag_%d", i)
			args[paramName] = tag.Name
			tagNames[i] = fmt.Sprintf("@%s", paramName)
		}

		// Add the tag filter to the query
		sql += fmt.Sprintf(" AND blog_tag.name IN (%s)", strings.Join(tagNames, ", "))
	} else {
		sql += " WHERE blog_post.title ILIKE '%' || @search || '%'"
	}

	value := 0
	err := s.Conn.QueryRow(config.CTX, sql, args).Scan(&value)
	if err != nil {
		return value, err
	}

	// If success return nil
	return value, nil
}

func (s *BlogPostService) GetWithSlug(slug string) (models.BlogPostContentWithTags, error) {
	// post SQL query
	postSql := `
		SELECT
			id,
			title,
			slug,
			content,
			created_at,
			updated_at,
			is_draft
		FROM blog_post
		WHERE slug = @slug;
	`

	// Add tag filters if tags are provided
	args := pgx.NamedArgs{
		"slug": slug,
	}

	// Execute post sql
	value := models.BlogPostContentWithTags{}
	err := s.Conn.QueryRow(config.CTX, postSql, args).Scan(
		&value.Id,
		&value.Title,
		&value.Slug,
		&value.Content,
		&value.CreatedAt,
		&value.UpdatedAt,
		&value.IsDraft,
	)
	if err != nil {
		return value, err
	}

	// Get tags from sql using parameterized query
	tagSql := `
			SELECT blog_tag.id, blog_tag.name
			FROM blog_tag
			INNER JOIN blog_post_tag ON blog_post_tag.tag_id = blog_tag.id
			WHERE blog_post_tag.post_id = $1;`
	tagRows, err := s.Conn.Query(config.CTX, tagSql, value.Id)
	if err != nil {
		return value, err
	}

	for tagRows.Next() {
		tagItem := models.BlogTag{}
		if err := tagRows.Scan(
			&tagItem.Id,
			&tagItem.Name,
		); err != nil {
			return value, err
		}
		value.Tags = append(value.Tags, tagItem)
	}

	return value, nil
}

func (s *BlogPostService) GetAll(search string, tags []models.BlogTag, limit int, page int) ([]models.BlogPostWithTags, error) {
	// Set default range for limit
	if limit < 10 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	// Set default range for page
	if page < 1 {
		page = 0
	} else {
		page -= 1
	}

	// post SQL query
	postSql := `
		SELECT
			blog_post.id,
			blog_post.title,
			blog_post.slug,
			blog_post.created_at,
			blog_post.updated_at,
			blog_post.is_draft
		FROM blog_post
	`

	// Add tag filters if tags are provided
	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"page":   page * limit,
	}

	if len(tags) > 0 {
		postSql += `
			INNER JOIN blog_post_tag ON blog_post_tag.post_id = blog_post.id
			INNER JOIN blog_tag ON blog_post_tag.tag_id = blog_tag.id
			WHERE blog_post.title ILIKE '%' || @search || '%'
		`

		// Add tag names to the query
		tagNames := make([]string, len(tags))
		for i, tag := range tags {
			paramName := fmt.Sprintf("tag_%d", i)
			args[paramName] = tag.Name
			tagNames[i] = fmt.Sprintf("@%s", paramName)
		}

		// Add the tag filter to the query
		postSql += fmt.Sprintf(" AND blog_tag.name IN (%s)", strings.Join(tagNames, ", "))

		// Group by blog post ID
		postSql += " GROUP BY blog_post.id"
		
		// Add HAVING clause to ensure all specified tags are matched
		// This ensures the blog post has ALL of the requested tags, not just any of them
		postSql += fmt.Sprintf(" HAVING COUNT(DISTINCT blog_tag.name) >= %d", len(tags))
	} else {
		postSql += "WHERE blog_post.title ILIKE '%' || @search || '%'"
	}

	// Add pagination
	postSql += " LIMIT @limit OFFSET @page;"

	// Execute post sql
	value := []models.BlogPostWithTags{}
	rows, err := s.Conn.Query(config.CTX, postSql, args)
	if err != nil {
		return value, err
	}

	for rows.Next() {
		postItem := models.BlogPostWithTags{}

		// Scan post
		if err := rows.Scan(
			&postItem.Id,
			&postItem.Title,
			&postItem.Slug,
			&postItem.CreatedAt,
			&postItem.UpdatedAt,
			&postItem.IsDraft,
		); err != nil {
			return value, err
		}

		// Get tags from sql using parameterized query
		tagSql := `
			SELECT blog_tag.id, blog_tag.name
			FROM blog_tag
			INNER JOIN blog_post_tag ON blog_post_tag.tag_id = blog_tag.id
			WHERE blog_post_tag.post_id = $1;`
		tagRows, err := s.Conn.Query(config.CTX, tagSql, postItem.Id)
		if err != nil {
			return value, err
		}

		for tagRows.Next() {
			tagItem := models.BlogTag{}
			if err := tagRows.Scan(
				&tagItem.Id,
				&tagItem.Name,
			); err != nil {
				return value, err
			}
			postItem.Tags = append(postItem.Tags, tagItem)
		}

		value = append(value, postItem)
	}

	return value, nil
}

func (s *BlogPostService) Create(input *models.BlogPostCreated) (models.BlogPostContentWithTags, error) {
	// Get slug string
	slugString := slug.Make(input.Title)

	// Create post
	postSql := `
		INSERT INTO blog_post (title, slug, content, created_at, updated_at, is_draft)
	 	VALUES (@title, @slug, @content, @created_at, @updated_at, @is_draft)
		RETURNING *;
	`
	postArgs := pgx.NamedArgs{
		"title":      input.Title,
		"slug":       slugString,
		"content":    input.Content,
		"created_at": input.CreatedAt,
		"updated_at": input.UpdatedAt,
		"is_draft":   input.IsDraft,
	}
	value := models.BlogPostContentWithTags{}
	err := s.Conn.QueryRow(config.CTX, postSql, postArgs).Scan(
		&value.Id,
		&value.Title,
		&value.Slug,
		&value.Content,
		&value.CreatedAt,
		&value.UpdatedAt,
		&value.IsDraft,
	)
	if err != nil {
		return value, err
	}

	// Delete tags and create new tags for post
	dropPostTagSql := "DELETE FROM blog_post_tag WHERE post_id = @post_id;"
	_, err = s.Conn.Exec(config.CTX, dropPostTagSql, pgx.NamedArgs{"post_id": value.Id})
	if err != nil {
		return value, err
	}

	for _, item := range input.Tags {
		postTagSql := "INSERT INTO blog_post_tag (tag_id, post_id) VALUES (@tag_id, @post_id);"
		postTagArgs := pgx.NamedArgs{
			"tag_id":  item.Id,
			"post_id": value.Id,
		}
		_, err := s.Conn.Exec(config.CTX, postTagSql, postTagArgs)
		if err != nil {
			return value, err
		}
	}

	// Query tags data and append to value.Tags
	tagSql := `
		SELECT blog_tag.id, blog_tag.name
		FROM blog_tag
		INNER JOIN blog_post_tag ON blog_post_tag.tag_id = blog_tag.id
		WHERE blog_post_tag.post_id = @post_id;
	`
	tagRows, err := s.Conn.Query(config.CTX, tagSql, pgx.NamedArgs{"post_id": value.Id})
	if err != nil {
		return value, err
	}

	for tagRows.Next() {
		tagItem := models.BlogTag{}
		if err := tagRows.Scan(&tagItem.Id, &tagItem.Name); err != nil {
			return value, err
		}
		value.Tags = append(value.Tags, tagItem)
	}

	// If success return nil
	return value, nil
}

func (s *BlogPostService) Update(input *models.BlogPostUpdated) (models.BlogPostContentWithTags, error) {
	// Get slug string
	slugString := slug.Make(input.Title)

	// Update post
	sql := `
		UPDATE blog_post SET
			title=@title,
			slug=@slug,
			content=@content,
			created_at=@created_at,
			updated_at=@updated_at,
			is_draft=@is_draft
		WHERE id=@id
		RETURNING *;
	`
	args := pgx.NamedArgs{
		"id":         input.Id,
		"title":      input.Title,
		"slug":       slugString,
		"content":    input.Content,
		"created_at": input.CreatedAt,
		"updated_at": input.UpdatedAt,
		"is_draft":   input.IsDraft,
	}
	value := models.BlogPostContentWithTags{}
	err := s.Conn.QueryRow(config.CTX, sql, args).Scan(
		&value.Id,
		&value.Title,
		&value.Slug,
		&value.Content,
		&value.CreatedAt,
		&value.UpdatedAt,
		&value.IsDraft,
	)
	if err != nil {
		return value, err
	}

	// Delete tags and create new tags for post
	dropPostTagSql := "DELETE FROM blog_post_tag WHERE post_id = @post_id;"
	_, err = s.Conn.Exec(config.CTX, dropPostTagSql, pgx.NamedArgs{"post_id": value.Id})
	if err != nil {
		return value, err
	}

	for _, item := range input.Tags {
		postTagSql := "INSERT INTO blog_post_tag (tag_id, post_id) VALUES (@tag_id, @post_id);"
		postTagArgs := pgx.NamedArgs{
			"tag_id":  item.Id,
			"post_id": value.Id,
		}
		_, err := s.Conn.Exec(config.CTX, postTagSql, postTagArgs)
		if err != nil {
			return value, err
		}
	}

	// Query tags data and append to value.Tags
	tagSql := `
		SELECT blog_tag.id, blog_tag.name
		FROM blog_tag
		INNER JOIN blog_post_tag ON blog_post_tag.tag_id = blog_tag.id
		WHERE blog_post_tag.post_id = @post_id;
	`
	tagRows, err := s.Conn.Query(config.CTX, tagSql, pgx.NamedArgs{"post_id": value.Id})
	if err != nil {
		return value, err
	}

	for tagRows.Next() {
		tagItem := models.BlogTag{}
		if err := tagRows.Scan(&tagItem.Id, &tagItem.Name); err != nil {
			return value, err
		}
		value.Tags = append(value.Tags, tagItem)
	}

	// If success return nil
	return value, nil
}

func (s *BlogPostService) Remove(id string) (string, error) {
	// Execute SQL
	sql := "DELETE FROM blog_post WHERE id=@id RETURNING id;"
	args := pgx.NamedArgs{
		"id": id,
	}
	value := ""
	err := s.Conn.QueryRow(config.CTX, sql, args).Scan(&value)
	if err != nil {
		return value, err
	}

	// If success return nil
	return value, nil
}
