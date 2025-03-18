package services

import (
	"api-chi/cmd/config"
	"api-chi/cmd/models"

	"github.com/jackc/pgx/v5"
)

type BlogTagService struct {
	Conn DatabaseService
}

func (s *BlogTagService) Open() error {
	return s.Conn.Open()
}

func (s *BlogTagService) Close() {
	s.Conn.Close()
}

func (s *BlogTagService) Count(search string) (int, error) {
	// Execute SQL
	sql := "SELECT COUNT(id) FROM blog_tag WHERE name ILIKE '%' || @search || '%';"
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

func (s *BlogTagService) GetAll(search string, limit int, page int) ([]models.BlogTag, error) {
	// Set default range for page
	if page < 1 {
		page = 0
	} else {
		page -= 1
	}

	// Execute SQL
	sql := "SELECT id, name FROM blog_tag WHERE name ILIKE '%' || @search || '%' LIMIT @limit OFFSET @page"
	args := pgx.NamedArgs{
		"search": search,
		"limit":  limit,
		"page":   page*limit,
	}
	value := []models.BlogTag{}
	rows, err := s.Conn.Query(config.CTX, sql, args)
	if err != nil {
		return value, err
	}
	for rows.Next() {
		item := models.BlogTag{}

		if err := rows.Scan(&item.Id, &item.Name); err != nil {
			return nil, err
		}

		value = append(value, item)
	}

	// If success return nil
	return value, nil
}

func (s *BlogTagService) Create(input *models.BlogTag) (models.BlogTag, error) {
	// Execute SQL
	sql := "INSERT INTO blog_tag (name) VALUES (@name) RETURNING id, name;"
	args := pgx.NamedArgs{
		"name": input.Name,
	}
	value := models.BlogTag{}
	err := s.Conn.QueryRow(config.CTX, sql, args).Scan(&value.Id, &value.Name)
	if err != nil {
		return value, err
	}

	// If success return nil
	return value, nil
}

func (s *BlogTagService) Update(input *models.BlogTag) (models.BlogTag, error) {
	// Execute SQL
	sql := "UPDATE blog_tag SET name=@name WHERE id=@id RETURNING id, name;"
	args := pgx.NamedArgs{
		"id":   input.Id,
		"name": input.Name,
	}
	value := models.BlogTag{}
	err := s.Conn.QueryRow(config.CTX, sql, args).Scan(&value.Id, &value.Name)
	if err != nil {
		return value, err
	}

	// If success return nil
	return value, nil
}

func (s *BlogTagService) Remove(id string) (string, error) {
	// Execute SQL
	sql := "DELETE FROM blog_tag WHERE id = @id RETURNING id;"
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
