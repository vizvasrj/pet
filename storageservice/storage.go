package storageservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"src/env"
	"src/myerror"
	"src/petstore"

	"github.com/fatih/color"
	_ "github.com/lib/pq"
)

func GetConnection(e *env.Env) *sql.DB {
	connStr := "user=" + e.PostgreUser + " password=" + e.PostgrePassword + " dbname=" + e.PostgreDB + " sslmode=" + e.PostgreSSLMode
	color.Green("Connection, %s", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

type StorageService struct {
	Db *sql.DB
}

func (s *StorageService) CreatePet(ctx context.Context, pet *petstore.NewPet) (*petstore.Pet, error) {
	// Implement the method
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

func (s *StorageService) DeletePet(ctx context.Context, id int64) error {
	// Implement the method
	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return myerror.WrapError(err, "failed to begin transaction")
	}
	defer tx.Rollback() // Rollback if any step fails

	if err := s.deletePetInTransaction(ctx, tx, id); err != nil {
		return myerror.WrapError(err, "failed to delete pet")
	}

	if err := tx.Commit(); err != nil {
		return myerror.WrapError(err, "failed to commit transaction")
	}

	return nil
}

func (s *StorageService) FindPetByID(ctx context.Context, id int64) (*petstore.Pet, error) {
	// Implement the method
	pet := &petstore.Pet{Id: id}
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
	err := s.Db.QueryRowContext(ctx, query, id).Scan(&pet.Id, &petName, &categoryId, &categoryName, &status, &tagsString)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, myerror.WrapError(err, "failed to find pet")
	}
	pet.Name = petName
	pet.Category = &petstore.Category{
		Id:   &categoryId,
		Name: &categoryName,
	}
	if status.Valid {
		pet.Status = &status.String
	}
	// Unmarshal the tags
	var tags []petstore.Tag
	err = json.Unmarshal([]byte(tagsString), &tags)
	if err != nil {
		return nil, myerror.WrapError(err, "failed to unmarshal tags")
	}
	pet.Tags = &tags

	return pet, nil
}

func (s *StorageService) FindPets(ctx context.Context, limit int64, tags []string) ([]*petstore.Pet, error) {
	// Implement the method
	return []*petstore.Pet{}, nil
}
