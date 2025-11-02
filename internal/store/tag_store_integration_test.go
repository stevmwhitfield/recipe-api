package store_test

import (
	"database/sql"
	"testing"

	"github.com/stevmwhitfield/recipe-api/internal/data/migrations"
	"github.com/stevmwhitfield/recipe-api/internal/model"
	"github.com/stevmwhitfield/recipe-api/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	err = store.MigrateFS(db, migrations.FS, ".")
	require.NoError(t, err)

	return db
}

func TestListTags_Integration(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	q := `
		INSERT INTO tags (id, name)
		VALUES ("1", "t1"), ("2", "t2")
	`

	_, err := db.Exec(q)
	require.NoError(t, err)

	tagStore := store.NewSQLiteTagStore(db)

	tags, err := tagStore.ListTags()

	assert.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, "t1", tags[0].Name)
	assert.Equal(t, "t2", tags[1].Name)
}

func TestCreateTag_Integration(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	tagStore := store.NewSQLiteTagStore(db)

	newTag := &model.Tag{
		ID:   "1",
		Name: "t1",
	}
	createdTag, err := tagStore.CreateTag(newTag)

	assert.NoError(t, err)
	assert.Equal(t, "t1", createdTag.Name)
}
