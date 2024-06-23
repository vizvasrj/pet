package storageservice

import (
	"context"
	"database/sql"
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
	// pet := &petstore.Pet{Id: id}
	// query := `
	// SELECT p.name, c.name, p.status, c.id AS category_id, c.name AS category_name
	// `
	// err := s.Db.QueryRowContext(ctx, `

	// `)
	return &petstore.Pet{}, nil
}

func (s *StorageService) FindPets(ctx context.Context, limit int64, tags []string) ([]*petstore.Pet, error) {
	// Implement the method
	return []*petstore.Pet{}, nil
}
