package models

type BlogTag struct {
	Id   string `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}
