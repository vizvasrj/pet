package service

import (
	"context"
	"database/sql"
	"src/myerror"
	"src/proto_storage"
)

func (s *StorageService) createPetInTransaction(ctx context.Context, tx *sql.Tx, pet *proto_storage.NewPet) (*proto_storage.Pet, error) {
	// 1. Find or Create Category
	category, err := s.findOrCreateCategory(ctx, tx, pet.Category.Name)
	if err != nil {
		return nil, myerror.WrapError(s.Logger, err, "failed to find or create category")
	}

	// 2. Create Pet
	createdPet, err := s.createPetRecord(ctx, tx, pet.Name, category.Id, &pet.Status)
	if err != nil {
		return nil, myerror.WrapError(s.Logger, err, "failed to create pet")
	}

	// 3. Create Pet Tags (Find or Create Tags)
	if pet.Tags != nil {
		tags, err := s.createPetTags(ctx, tx, createdPet.Id, pet.Tags)
		if err != nil {
			return nil, myerror.WrapError(s.Logger, err, "failed to create pet tags")
		}
		pet.Tags = tags
	}

	// 4. Create Pet Photos
	if pet.PhotoUrls != nil {
		if err := s.createPetPhotos(ctx, tx, createdPet.Id, pet.PhotoUrls); err != nil {
			return nil, myerror.WrapError(s.Logger, err, "failed to create pet photos")
		}
	}

	return &proto_storage.Pet{
		Id:        createdPet.Id,
		Category:  category,
		PhotoUrls: pet.PhotoUrls,
		Name:      pet.Name,
		Status:    pet.Status,
		Tags:      pet.Tags,
	}, nil
}

// Helper functions for each step:

func (s *StorageService) findOrCreateCategory(ctx context.Context, tx *sql.Tx, categoryName string) (*proto_storage.Category, error) {
	category := &proto_storage.Category{}
	err := tx.QueryRowContext(ctx, "SELECT id, name FROM categories WHERE name = $1", categoryName).Scan(&category.Id, &category.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			// Category doesn't exist, create it
			err = tx.QueryRowContext(ctx, "INSERT INTO categories (name) VALUES ($1) RETURNING id, name", categoryName).Scan(&category.Id, &category.Name)
			if err != nil {
				return nil, myerror.WrapError(s.Logger, err, "failed to create category")
			}
		} else {
			return nil, myerror.WrapError(s.Logger, err, "failed to find category")
		}
	}
	return category, nil
}

func (s *StorageService) createPetRecord(ctx context.Context, tx *sql.Tx, name string, categoryID int64, status *string) (*proto_storage.Pet, error) {
	pet := proto_storage.Pet{}
	var returnId int64
	var returnName string
	var returnStatus string
	var returnCategoryId int64
	query := `
	INSERT INTO pets (name, category_id, status) VALUES ($1, $2, $3) RETURNING id, name, category_id, status
	`
	err := tx.QueryRowContext(ctx, query, name, categoryID, status).Scan(&returnId, &returnName, &returnCategoryId, &returnStatus)
	if err != nil {
		return nil, myerror.WrapError(s.Logger, err, "failed to create pet record")
	}
	pet.Id = returnId
	pet.Name = returnName
	pet.Status = returnStatus
	pet.Category = &proto_storage.Category{Id: returnCategoryId}

	return &pet, nil
}

func (s *StorageService) createPetTags(ctx context.Context, tx *sql.Tx, petID int64, tags []*proto_storage.Tag) ([]*proto_storage.Tag, error) {
	newTags := []*proto_storage.Tag{}
	if tags == nil {
		return nil, nil
	}
	for _, tag := range tags {
		newTag := &proto_storage.Tag{}
		err := tx.QueryRowContext(ctx, "SELECT id, name FROM tags WHERE name = $1", tag.Name).Scan(&newTag.Id, &newTag.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				// Tag doesn't exist, create it
				err = tx.QueryRowContext(ctx, "INSERT INTO tags (name) VALUES ($1) RETURNING id, name", tag.Name).Scan(&newTag.Id, &newTag.Name)
				if err != nil {
					return nil, myerror.WrapError(s.Logger, err, "failed to create tag")
				}
			} else {
				return nil, myerror.WrapError(s.Logger, err, "failed to find tag")
			}
		}
		newTags = append(newTags, newTag)

		// Insert into pet_tags
		_, err = tx.ExecContext(ctx, "INSERT INTO pet_tags (pet_id, tag_id) VALUES ($1, $2)", petID, newTag.Id)
		if err != nil {
			return nil, myerror.WrapError(s.Logger, err, "failed to insert into pet_tags")
		}
	}
	return newTags, nil
}

func (s *StorageService) createPetPhotos(ctx context.Context, tx *sql.Tx, petID int64, photoUrls []string) error {
	for _, url := range photoUrls {
		// Insert into pet_photos
		_, err := tx.ExecContext(ctx, "INSERT INTO pet_photos (pet_id, url) VALUES ($1, $2)", petID, url)
		if err != nil {
			return myerror.WrapError(s.Logger, err, "failed to insert into pet_photos")
		}
	}
	return nil
}
