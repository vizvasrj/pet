package storageservice

import (
	"context"
	"database/sql"
	"src/env"
	"src/petstore"

	_ "github.com/lib/pq"
)

func GetConnection(e *env.Env) *sql.DB {
	connStr := "user=" + e.PostgreUser + " password=" + e.PostgrePassword + " dbname=" + e.PostgreDB + " sslmode=" + e.PostgreSSLMode
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
	return &petstore.Pet{}, nil
}

func (s *StorageService) DeletePet(ctx context.Context, id int64) error {
	// Implement the method
	return nil
}

func (s *StorageService) FindPetByID(ctx context.Context, id int64) (*petstore.Pet, error) {
	// Implement the method
	return &petstore.Pet{}, nil
}

func (s *StorageService) FindPets(ctx context.Context, limit int64, tags []string) ([]*petstore.Pet, error) {
	// Implement the method
	return []*petstore.Pet{}, nil
}
