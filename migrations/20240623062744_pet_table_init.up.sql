-- Add up migration script here
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE CHECK (name <> '')
);

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE CHECK (name <> '')
);

CREATE TABLE pets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (name <> ''),
    category_id INTEGER REFERENCES categories(id),
    status VARCHAR(255)
);

CREATE TABLE pet_tags (
    pet_id INTEGER REFERENCES pets(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (pet_id, tag_id)
);

CREATE TABLE pet_photos (
    pet_id INTEGER REFERENCES pets(id) ON DELETE CASCADE,
    url VARCHAR(255)
);