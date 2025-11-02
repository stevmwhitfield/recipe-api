package store

import (
	"database/sql"
	"time"
)

type Recipe struct {
	ID              string        `json:"id"`
	Slug            string        `json:"slug"`
	Name            string        `json:"name"`
	Servings        int           `json:"servings"`
	PrepTimeSeconds int           `json:"prepTimeSeconds"`
	CookTimeSeconds int           `json:"cookTimeSeconds"`
	Ingredients     []Ingredient  `json:"ingredients"`
	Instructions    []Instruction `json:"instructions"`
	Tags            []Tag         `json:"tags"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

type Ingredient struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Note     string  `json:"note"`
}

type Instruction struct {
	ID          string `json:"id"`
	StepNumber  int    `json:"stepNumber"`
	Description string `json:"description"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SQLiteRecipeStore struct {
	db *sql.DB
}

func NewSQLiteRecipeStore(db *sql.DB) *SQLiteRecipeStore {
	return &SQLiteRecipeStore{db: db}
}

type RecipeStore interface {
	ListRecipes() ([]Recipe, error)
	CreateRecipe(*Recipe) (*Recipe, error)
	GetRecipeByID(id string) (*Recipe, error)
	UpdateRecipe(*Recipe) (*Recipe, error)
	DeleteRecipe(id string) error
}

func (s *SQLiteRecipeStore) ListRecipes() ([]Recipe, error) {
	query := `
		SELECT id, slug, name, servings, prep_time_seconds, cook_time_seconds, created_at, updated_at 
		FROM recipes
		ORDER BY name ASC;
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		err = rows.Scan(&r.ID, &r.Slug, &r.Name, &r.Servings, &r.PrepTimeSeconds, &r.CookTimeSeconds, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if r.Ingredients, err = s.getIngredientsForRecipe(r.ID); err != nil {
			return nil, err
		}
		if r.Instructions, err = s.getInstructionsForRecipe(r.ID); err != nil {
			return nil, err
		}
		if r.Tags, err = s.getTagsForRecipe(r.ID); err != nil {
			return nil, err
		}

		recipes = append(recipes, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return recipes, nil
}

func (s *SQLiteRecipeStore) CreateRecipe(recipe *Recipe) (*Recipe, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO recipes (id, slug, name, servings, prep_time_seconds, cook_time_seconds)
		VALUES (?, ?, ?, ?, ?, ?);
	`

	_, err = tx.Exec(query, recipe.ID, recipe.Slug, recipe.Name, recipe.Servings, recipe.PrepTimeSeconds, recipe.CookTimeSeconds)
	if err != nil {
		return nil, err
	}

	for _, i := range recipe.Ingredients {
		query := `
			INSERT INTO recipe_ingredient (recipe_id, ingredient_id, quantity, unit, note)
			VALUES (?, ?, ?, ?, ?);
		`

		_, err = tx.Exec(query, recipe.ID, i.ID, i.Quantity, i.Unit, i.Note)
		if err != nil {
			return nil, err
		}
	}

	for _, i := range recipe.Instructions {
		query := `
			INSERT INTO instructions (id, recipe_id, step_number, description)
			VALUES (?, ?, ?, ?);
		`

		_, err = tx.Exec(query, i.ID, recipe.ID, i.StepNumber, i.Description)
		if err != nil {
			return nil, err
		}
	}

	for _, t := range recipe.Tags {
		query := `
			INSERT INTO recipe_tag (recipe_id, tag_id)
			VALUES (?, ?);
		`

		_, err = tx.Exec(query, recipe.ID, t.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func (s *SQLiteRecipeStore) GetRecipeByID(id string) (*Recipe, error) {
	r := &Recipe{}
	query := `
		SELECT id, slug, name, servings, prep_time_seconds, cook_time_seconds, created_at, updated_at 
		FROM recipes
		WHERE id = ?;
	`

	err := s.db.QueryRow(query, id).Scan(&r.ID, &r.Slug, &r.Name, &r.Servings, &r.PrepTimeSeconds, &r.CookTimeSeconds, &r.CreatedAt, &r.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if r.Ingredients, err = s.getIngredientsForRecipe(r.ID); err != nil {
		return nil, err
	}
	if r.Instructions, err = s.getInstructionsForRecipe(r.ID); err != nil {
		return nil, err
	}
	if r.Tags, err = s.getTagsForRecipe(r.ID); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *SQLiteRecipeStore) UpdateRecipe(recipe *Recipe) (*Recipe, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		UPDATE recipes
		SET name = ?, servings = ?, prep_time_seconds = ?, cook_time_seconds = ?
		WHERE id = ?;
	`

	result, err := tx.Exec(query, recipe.Name, recipe.Servings, recipe.PrepTimeSeconds, recipe.CookTimeSeconds, recipe.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	_, err = tx.Exec(`DELETE FROM recipe_ingredient WHERE recipe_id = ?`, recipe.ID)
	if err != nil {
		return nil, err
	}

	for _, i := range recipe.Ingredients {
		query := `
			INSERT INTO recipe_ingredient (recipe_id, ingredient_id, quantity, unit, note)
			VALUES (?, ?, ?, ?, ?);
		`

		_, err = tx.Exec(query, recipe.ID, i.ID, i.Quantity, i.Unit, i.Note)
		if err != nil {
			return nil, err
		}
	}

	_, err = tx.Exec(`DELETE FROM instructions WHERE recipe_id = ?`, recipe.ID)
	if err != nil {
		return nil, err
	}

	for _, i := range recipe.Instructions {
		query := `
			INSERT INTO instructions (id, recipe_id, step_number, description)
			VALUES (?, ?, ?, ?);
		`

		_, err = tx.Exec(query, i.ID, recipe.ID, i.StepNumber, i.Description)
		if err != nil {
			return nil, err
		}
	}

	_, err = tx.Exec(`DELETE FROM recipe_tag WHERE recipe_id = ?`, recipe.ID)
	if err != nil {
		return nil, err
	}

	for _, t := range recipe.Tags {
		query := `
			INSERT INTO recipe_tag (recipe_id, tag_id)
			VALUES (?, ?);
		`

		_, err = tx.Exec(query, recipe.ID, t.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func (s *SQLiteRecipeStore) DeleteRecipe(id string) error {
	query := `
		DELETE FROM recipes
		WHERE id = ?;
	`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *SQLiteRecipeStore) getIngredientsForRecipe(recipeID string) ([]Ingredient, error) {
	query := `
		SELECT i.id, i.name, ri.quantity, ri.unit, ri.note
		FROM recipe_ingredient ri
		JOIN ingredients i ON i.id = ri.ingredient_id
		WHERE ri.recipe_id = ?;
	`

	rows, err := s.db.Query(query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ingredients := []Ingredient{}
	for rows.Next() {
		var i Ingredient
		err = rows.Scan(&i.ID, &i.Name, &i.Quantity, &i.Unit, &i.Note)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, i)
	}
	return ingredients, rows.Err()
}

func (s *SQLiteRecipeStore) getInstructionsForRecipe(recipeID string) ([]Instruction, error) {
	query := `
		SELECT id, step_number, description
		FROM instructions
		WHERE recipe_id = ?
		ORDER BY step_number ASC;
	`

	rows, err := s.db.Query(query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	instructions := []Instruction{}
	for rows.Next() {
		var i Instruction
		err = rows.Scan(&i.ID, &i.StepNumber, &i.Description)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions, i)
	}
	return instructions, rows.Err()
}

func (s *SQLiteRecipeStore) getTagsForRecipe(recipeID string) ([]Tag, error) {
	query := `
		SELECT t.id, t.name
		FROM recipe_tag rt
		JOIN tags t ON t.id = rt.tag_id
		WHERE rt.recipe_id = ?;
	`

	rows, err := s.db.Query(query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []Tag{}
	for rows.Next() {
		var t Tag
		err = rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}
