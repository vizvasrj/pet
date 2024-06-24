package grpcstorageservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"src/myerror"
	protostorageservice "src/protoStorageService"

	"github.com/fatih/color"
)

type StorageService struct {
	protostorageservice.UnimplementedStorageServiceServer
	Db *sql.DB
}

func (s *StorageService) CreatePet(ctx context.Context, pet *protostorageservice.NewPet) (*protostorageservice.Pet, error) {
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, myerror.WrapError(err, "failed to begin transaction")
	}
	defer tx.Rollback() // Rollback if any step fails

	createdPet, err := s.createPetInTransaction(ctx, tx, pet)
	if err != nil {
		return nil, myerror.WrapError(err, "") // Error will be specific to a step
	}

	if err := tx.Commit(); err != nil {
		return nil, myerror.WrapError(err, "failed to commit transaction")
	}

	return createdPet, nil
}

func (s *StorageService) FindPets(ctx context.Context, req *protostorageservice.FindPetsRequest) (*protostorageservice.FindPetsResponse, error) {
	color.Red("hi there")

	// 1. Build the base query
	query := `
	SELECT 
		p.id, 
		p.name, 
		p.category_id, 
		c.name as category_name, 
		p.status,
		json_agg(json_build_object('id', t.id, 'name', t.name)) as tags
	FROM pets p
	JOIN categories c ON p.category_id = c.id
	LEFT JOIN pet_tags pt ON p.id = pt.pet_id
	LEFT JOIN tags t ON pt.tag_id = t.id
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
		return nil, myerror.WrapError(err, "failed to find pets")
	}
	defer rows.Close()

	// 5. Process the result set
	pets := []*protostorageservice.Pet{}
	for rows.Next() {
		var (
			petId        int64
			petName      string
			categoryId   int64
			categoryName string
			status       sql.NullString
			tagsString   string
		)

		if err := rows.Scan(&petId, &petName, &categoryId, &categoryName, &status, &tagsString); err != nil {
			return nil, myerror.WrapError(err, "failed to scan pet row")
		}

		// Unmarshal the tags
		var tags []*protostorageservice.Tag
		err = json.Unmarshal([]byte(tagsString), &tags)
		if err != nil {
			return nil, myerror.WrapError(err, "failed to unmarshal tags")
		}
		// Create pet object
		pet := &protostorageservice.Pet{
			Id:   petId,
			Name: petName,
			Category: &protostorageservice.Category{
				Id:   categoryId,
				Name: categoryName,
			},
			Status: "", // Initialize status as empty string
			Tags:   tags,
		}

		// Set status only if it's not NULL in the database
		if status.Valid {
			pet.Status = status.String
		}
		pets = append(pets, pet)
	}
	resp := protostorageservice.FindPetsResponse{
		Pets: pets,
	}

	return &resp, nil
}

func (s *StorageService) FindPetById(ctx context.Context, req *protostorageservice.PetID) (*protostorageservice.Pet, error) {
	color.Red("hi there")
	// Implement the method
	pet := &protostorageservice.Pet{Id: req.Id}
	query := `
	SELECT 
		p.id, 
		p.name, 
		p.category_id, 
		c.name as category_name, 
		p.status,
		json_agg(json_build_object('id', t.id, 'name', t.name)) as tags
	FROM pets p
	JOIN categories c ON p.category_id = c.id
	JOIN pet_tags pt ON p.id = pt.pet_id
	JOIN tags t ON pt.tag_id = t.id
	WHERE p.id = $1
	GROUP BY p.id, p.name, p.category_id, c.name, p.status`
	var (
		// petId        int64
		petName      string
		categoryId   int64
		categoryName string
		status       sql.NullString
		tagsString   string
	)
	err := s.Db.QueryRowContext(ctx, query, req.Id).Scan(&pet.Id, &petName, &categoryId, &categoryName, &status, &tagsString)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, myerror.WrapError(err, "failed to find pet")
	}
	pet.Name = petName
	pet.Category = &protostorageservice.Category{
		Id:   categoryId,
		Name: categoryName,
	}
	if status.Valid {
		pet.Status = status.String
	}
	// Unmarshal the tags
	var tags []*protostorageservice.Tag
	err = json.Unmarshal([]byte(tagsString), &tags)
	if err != nil {
		return nil, myerror.WrapError(err, "failed to unmarshal tags")
	}
	pet.Tags = tags

	return pet, nil
}

func (s *StorageService) DeletePet(ctx context.Context, req *protostorageservice.PetID) (*protostorageservice.Empty, error) {
	// Implement the method
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Rollback if any step fails
	if err := s.deletePetInTransaction(ctx, tx, req.Id); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, myerror.WrapError(err, "failed to commit transaction")
	}

	return nil, nil
}
