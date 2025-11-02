package store

import (
	"database/sql"

	"github.com/stevmwhitfield/recipe-api/internal/model"
)

type SQLiteTagStore struct {
	db *sql.DB
}

func NewSQLiteTagStore(db *sql.DB) *SQLiteTagStore {
	return &SQLiteTagStore{db: db}
}

type TagStore interface {
	ListTags() ([]model.Tag, error)
	CreateTag(*model.Tag) (*model.Tag, error)
}

func (s *SQLiteTagStore) ListTags() ([]model.Tag, error) {
	q := `
		SELECT id, name
		FROM tags
		ORDER BY name ASC;
	`

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []model.Tag{}
	for rows.Next() {
		var t model.Tag
		err = rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (s *SQLiteTagStore) CreateTag(t *model.Tag) (*model.Tag, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	q := `
		INSERT INTO tags (id, name)
		VALUES (?, ?);
	`

	_, err = tx.Exec(q, t.ID, t.Name)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return t, nil
}
