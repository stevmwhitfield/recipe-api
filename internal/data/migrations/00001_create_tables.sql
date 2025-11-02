-- +goose Up

-- Reference tables
CREATE TABLE ingredients (
    id TEXT PRIMARY KEY, -- UUIDv7
    name TEXT NOT NULL,        
    category TEXT NOT NULL
);

CREATE TABLE tags (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Main recipe table
CREATE TABLE recipes (
    id TEXT PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    servings INTEGER,
    prep_time_seconds INTEGER,
    cook_time_seconds INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Related tables
CREATE TABLE instructions (
    id TEXT PRIMARY KEY,
    recipe_id TEXT NOT NULL,
    step_number INTEGER NOT NULL,
    description TEXT NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
);

CREATE TABLE recipe_ingredient (
    recipe_id TEXT NOT NULL,
    ingredient_id TEXT NOT NULL,
    quantity REAL NOT NULL,
    unit TEXT NOT NULL,
    note TEXT,
    PRIMARY KEY (recipe_id, ingredient_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id)
);

CREATE TABLE recipe_tag (
    recipe_id TEXT NOT NULL,
    tag_id TEXT NOT NULL,
    PRIMARY KEY (recipe_id, tag_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

CREATE TABLE stocked_ingredients (
    id TEXT PRIMARY KEY,
    ingredient_id TEXT NOT NULL,
    quantity REAL NOT NULL,
    unit TEXT NOT NULL,
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id)
);

-- Indexes for common queries
CREATE INDEX idx_recipe_name ON recipes(name);                  -- searching recipes by name (e.g. name ~ chicken)
CREATE INDEX idx_ingredient_name ON ingredients(name);          -- searching ingredients by name (e.g. name ~ flour)
CREATE INDEX idx_ingredient_category ON ingredients(category);  -- searching ingredients by category (e.g. category = dairy)
CREATE INDEX idx_instruction_recipe ON instructions(recipe_id); -- searching instructions for recipe

-- +goose Down

DROP TABLE stocked_ingredients;
DROP TABLE recipe_tag;
DROP TABLE recipe_ingredient;
DROP TABLE instructions;
DROP TABLE recipes;
DROP TABLE tags;
DROP TABLE ingredients;
