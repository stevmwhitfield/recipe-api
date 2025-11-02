package model

import "time"

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
