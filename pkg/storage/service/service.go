package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"src/myerror"
	"src/proto_storage"

	"go.uber.org/zap"
)

type StorageService struct {
	proto_storage.UnimplementedStorageServiceServer
	Db     *sql.DB
	Logger *zap.Logger
}

func (s *StorageService) CreatePet(ctx context.Context, pet *proto_storage.NewPet) (*proto_storage.Pet, error) {
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, myerror.WrapError(s.Logger, err, "failed to begin transaction")
	}
	defer tx.Rollback() // Rollback if any step fails

	createdPet, err := s.createPetInTransaction(ctx, tx, pet)
	if err != nil {
		return nil, myerror.WrapError(s.Logger, err, "") // Error will be specific to a step
	}

	if err := tx.Commit(); err != nil {
		return nil, myerror.WrapError(s.Logger, err, "failed to commit transaction")
	}

	return createdPet, nil
}

type Photo struct {
	Url string `json:"url"`
}

func (s *StorageService) FindPets(ctx context.Context, req *proto_storage.FindPetsRequest) (*proto_storage.FindPetsResponse, error) {
	// 1. Build the base query
	query := `
	SELECT 
		p.id, 
		p.name, 
		p.category_id, 
		c.name as category_name, 
		p.status,
		(
			SELECT json_agg(json_build_object('id', t.id, 'name', t.name))
			FROM pet_tags pt 
			JOIN tags t ON pt.tag_id = t.id
			WHERE pt.pet_id = p.id
		) as tags,
		(
			SELECT json_agg(json_build_object('url', pp.url))
			FROM pet_photos pp
			WHERE pp.pet_id = p.id
		) as photos
	FROM pets p
	LEFT JOIN categories c ON p.category_id = c.id

	`
	// 2. Add tag filtering if tags are provided
	var args []interface{}
	if len(req.Tags) > 0 {
		query += " WHERE t.name IN ("
		for i := range req.Tags {
			query += fmt.Sprintf("$%d,", i+1) // Add placeholders for tags
			args = append(args, req.Tags[i])  // Append tag values to arguments
		}
		query = query[:len(query)-1] + ")" // Remove trailing comma and close IN clause
	}

	// 3. Add grouping and limit
	query += " GROUP BY p.id, p.name, p.category_id, c.name, p.status "
	if req.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	// 4. Execute the query
	rows, err := s.Db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, myerror.WrapError(s.Logger, err, "failed to find pets")
	}
	defer rows.Close()

	// 5. Process the result set
	pets := []*proto_storage.Pet{}
	for rows.Next() {
		var (
			petId        int64
			petName      string
			categoryId   int64
			categoryName string
			status       sql.NullString
			tagsString   sql.NullString
			photosString sql.NullString
		)

		if err := rows.Scan(&petId, &petName, &categoryId, &categoryName, &status, &tagsString, &photosString); err != nil {
			return nil, myerror.WrapError(s.Logger, err, "failed to scan pet row")
		}
		var photos []Photo
		if photosString.Valid {
			err := json.Unmarshal([]byte(photosString.String), &photos)
			if err != nil {
				return nil, myerror.WrapError(s.Logger, err, "failed to unmarshal photos")
			}
		}
		photoStrings := []string{}
		for _, p := range photos {
			photoStrings = append(photoStrings, p.Url)
		}

		// Unmarshal the tags
		var tags []*proto_storage.Tag
		if tagsString.Valid {
			err := json.Unmarshal([]byte(tagsString.String), &tags)
			if err != nil {
				return nil, myerror.WrapError(s.Logger, err, "failed to unmarshal tags")
			}
		}

		// photoUrls := make([]string, 0, len(photosString))
		// for _, photo := range photosString {
		// 	if photo.Valid { // Check if the URL is not NULL
		// 		photoUrls = append(photoUrls, photo.String)
		// 	}
		// }
		// Create pet object

		pet := &proto_storage.Pet{
			Id:   petId,
			Name: petName,
			Category: &proto_storage.Category{
				Id:   categoryId,
				Name: categoryName,
			},
			Status:    "", // Initialize status as empty string
			Tags:      tags,
			PhotoUrls: photoStrings,
		}

		// Set status only if it's not NULL in the database
		if status.Valid {
			pet.Status = status.String
		}
		pets = append(pets, pet)
	}
	resp := proto_storage.FindPetsResponse{
		Pets: pets,
	}

	return &resp, nil
}

func (s *StorageService) FindPetById(ctx context.Context, req *proto_storage.PetID) (*proto_storage.Pet, error) {
	// Implement the method
	pet := &proto_storage.Pet{Id: req.Id}
	query := `
	SELECT 
		p.id, 
		p.name, 
		p.category_id, 
		c.name as category_name, 
		p.status,
		(
			SELECT json_agg(json_build_object('id', t.id, 'name', t.name))
			FROM pet_tags pt 
			JOIN tags t ON pt.tag_id = t.id
			WHERE pt.pet_id = p.id
		) as tags,
		(
			SELECT json_agg(json_build_object('url', pp.url))
			FROM pet_photos pp
			WHERE pp.pet_id = p.id
		) as photos
	FROM pets p
	LEFT JOIN categories c ON p.category_id = c.id
	WHERE p.id = $1
	`
	var (
		// petId        int64
		petName      string
		categoryId   int64
		categoryName string
		status       sql.NullString
		tagsString   sql.NullString
		photosString sql.NullString
	)
	err := s.Db.QueryRowContext(ctx, query, req.Id).Scan(&pet.Id, &petName, &categoryId, &categoryName, &status, &tagsString, &photosString)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, myerror.WrapError(s.Logger, err, "failed to find pet")
	}
	pet.Name = petName
	pet.Category = &proto_storage.Category{
		Id:   categoryId,
		Name: categoryName,
	}
	if status.Valid {
		pet.Status = status.String
	}
	// Unmarshal the tags
	var tags []*proto_storage.Tag
	if tagsString.Valid {
		err = json.Unmarshal([]byte(tagsString.String), &tags)
		if err != nil {
			return nil, myerror.WrapError(s.Logger, err, "failed to unmarshal tags")
		}
	}
	pet.Tags = tags

	if photosString.Valid {
		var photos []Photo
		err := json.Unmarshal([]byte(photosString.String), &photos)
		if err != nil {
			return nil, myerror.WrapError(s.Logger, err, "failed to unmarshal photos")
		}
		photoStrings := []string{}
		for _, p := range photos {
			photoStrings = append(photoStrings, p.Url)
		}
		pet.PhotoUrls = photoStrings
	}
	// pet.PhotoUrls = photosString
	return pet, nil
}

func (s *StorageService) DeletePet(ctx context.Context, req *proto_storage.PetID) (*proto_storage.Empty, error) {
	// Implement the method
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Rollback if any step fails
	if err := s.deletePetInTransaction(ctx, tx, req.Id); err != nil {
		return nil, myerror.WrapError(s.Logger, err, "") // Error will be specific to a step
	}

	if err := tx.Commit(); err != nil {
		return nil, myerror.WrapError(s.Logger, err, "failed to commit transaction")
	}

	return nil, nil
}
