package services

import (
	"api-chi/cmd/config"
	"api-chi/cmd/models"

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

func (s *BlogPostService) Count(search string) (int, error) {
	// Execute SQL
	sql := "SELECT * FROM count_blog_post(@search);"
	args := pgx.NamedArgs{
		"search": search,
	}
	value := 0
	err := s.Conn.QueryRow(config.CTX, sql, args).Scan(&value)
	if err != nil {
		return value, err
	}

	// If success return nil
	return value, nil
}

func (s *BlogPostService) GetAll(search string, tags string, limit int, page int) ([]models.BlogPost, error) {
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

	// Execute SQL
	sql := "SELECT * FROM get_all_blog_post(@search, @tags, @limit, @page);"
	args := pgx.NamedArgs{
		"search": search,
		"tags":   tags,
		"limit":  limit,
		"page":   page,
	}
	value := []models.BlogPost{}
	rows, err := s.Conn.Query(config.CTX, sql, args)
	if err != nil {
		return value, err
	}
	for rows.Next() {
		item := models.BlogPost{}

		if err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.Slug,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.IsDraft,
		); err != nil {
			return nil, err
		}

		value = append(value, item)
	}

	// If success return nil
	return value, nil
}

func (s *BlogPostService) Create(input *models.BlogPost) (models.BlogPost, error) {
	// Execute SQL
	sql := "SELECT * FROM create_blog_post(@title, @slug, @content, @created_at, @updated_at, @is_draft);"
	args := pgx.NamedArgs{
		"title":      input.Title,
		"slug":       input.Slug,
		"content":    input.Content,
		"created_at": input.CreatedAt,
		"updated_at": input.UpdatedAt,
		"is_draft":   input.IsDraft,
	}
	value := models.BlogPost{}
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

	// If success return nil
	return value, nil
}

func (s *BlogPostService) Update(input *models.BlogPost) (models.BlogPost, error) {
	// Execute SQL
	sql := "SELECT * FROM update_blog_post(@id, @title, @slug, @content, @created_at, @updated_at, @is_draft);"
	args := pgx.NamedArgs{
		"id":         input.Id,
		"title":      input.Title,
		"slug":       input.Slug,
		"content":    input.Content,
		"created_at": input.CreatedAt,
		"updated_at": input.UpdatedAt,
		"is_draft":   input.IsDraft,
	}
	value := models.BlogPost{}
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

	// If success return nil
	return value, nil
}

func (s *BlogPostService) Remove(id string) (models.BlogPost, error) {
	// Execute SQL
	sql := "SELECT * FROM remove_blog_post(@id);"
	args := pgx.NamedArgs{
		"id": id,
	}
	value := models.BlogPost{}
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

	// If success return nil
	return value, nil
}
