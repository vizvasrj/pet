package storageservice

import (
	"context"
	"database/sql"
	"src/myerror"
	"src/petstore"

	"github.com/fatih/color"
)

func (s *StorageService) createPetInTransaction(ctx context.Context, tx *sql.Tx, pet *petstore.NewPet) (*petstore.Pet, error) {
	// 1. Find or Create Category
	categoryID, err := s.findOrCreateCategory(ctx, tx, *pet.Category.Name)
	if err != nil {
		return nil, myerror.WrapError(err, "failed to find or create category")
	}

	// 2. Create Pet
	createdPetID, err := s.createPetRecord(ctx, tx, pet.Name, categoryID, pet.Status)
	if err != nil {
		return nil, myerror.WrapError(err, "failed to create pet")
	}

	// 3. Create Pet Tags (Find or Create Tags)
	if pet.Tags != nil {
		if err := s.createPetTags(ctx, tx, createdPetID, *pet.Tags); err != nil {
			return nil, myerror.WrapError(err, "failed to create pet tags")
		}
	}

	// 4. Create Pet Photos
	if pet.PhotoUrls != nil {
		if err := s.createPetPhotos(ctx, tx, createdPetID, *pet.PhotoUrls); err != nil {
			return nil, myerror.WrapError(err, "failed to create pet photos")
		}
	}

	return &petstore.Pet{
		Id:       createdPetID,
		Category: pet.Category,
	}, nil
}

// Helper functions for each step:

func (s *StorageService) findOrCreateCategory(ctx context.Context, tx *sql.Tx, categoryName string) (int64, error) {
	var categoryID int64
	color.Red("%#v, ", categoryName)
	err := tx.QueryRowContext(ctx, "SELECT id FROM categories WHERE name = $1", categoryName).Scan(&categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Category doesn't exist, create it
			err = tx.QueryRowContext(ctx, "INSERT INTO categories (name) VALUES ($1) RETURNING id", categoryName).Scan(&categoryID)
			if err != nil {
				return 0, myerror.WrapError(err, "failed to create category")
			}
		} else {
			return 0, myerror.WrapError(err, "failed to find category")
		}
	}
	return categoryID, nil
}

func (s *StorageService) createPetRecord(ctx context.Context, tx *sql.Tx, name string, categoryID int64, status *string) (int64, error) {
	var petID int64
	err := tx.QueryRowContext(ctx, "INSERT INTO pets (name, category_id, status) VALUES ($1, $2, $3) RETURNING id", name, categoryID, status).Scan(&petID)
	if err != nil {
		return 0, myerror.WrapError(err, "failed to create pet record")
	}
	return petID, nil
}

func (s *StorageService) createPetTags(ctx context.Context, tx *sql.Tx, petID int64, tags []petstore.Tag) error {
	if tags == nil {
		return nil
	}
	for _, tag := range tags {
		var tagID int64
		err := tx.QueryRowContext(ctx, "SELECT id FROM tags WHERE name = $1", tag.Name).Scan(&tagID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Tag doesn't exist, create it
				err = tx.QueryRowContext(ctx, "INSERT INTO tags (name) VALUES ($1) RETURNING id", tag.Name).Scan(&tagID)
				if err != nil {
					return myerror.WrapError(err, "failed to create tag")
				}
			} else {
				return myerror.WrapError(err, "failed to find tag")
			}
		}

		// Insert into pet_tags
		_, err = tx.ExecContext(ctx, "INSERT INTO pet_tags (pet_id, tag_id) VALUES ($1, $2)", petID, tagID)
		if err != nil {
			return myerror.WrapError(err, "failed to insert into pet_tags")
		}
	}
	return nil
}

func (s *StorageService) createPetPhotos(ctx context.Context, tx *sql.Tx, petID int64, photoUrls []string) error {
	for _, url := range photoUrls {
		// Insert into pet_photos
		_, err := tx.ExecContext(ctx, "INSERT INTO pet_photos (pet_id, url) VALUES ($1, $2)", petID, url)
		if err != nil {
			return myerror.WrapError(err, "failed to insert into pet_photos")
		}
	}
	return nil
}
