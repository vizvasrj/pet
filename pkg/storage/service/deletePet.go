package service

import (
	"context"
	"database/sql"
	"src/myerror"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *StorageService) deletePetInTransaction(ctx context.Context, tx *sql.Tx, id int64) error {
	// Implement the method
	// 1. Delete from pet_tags
	_, err := tx.ExecContext(ctx, "DELETE FROM pet_tags WHERE pet_id = $1", id)
	if err != nil {
		return myerror.WrapError(s.Logger, err, "failed to delete pet tags")
	}

	// 2. Delete from pet_photospetstore.ErrPetNotFound
	_, err = tx.ExecContext(ctx, "DELETE FROM pet_photos WHERE pet_id = $1", id)
	if err != nil {
		return myerror.WrapError(s.Logger, err, "failed to delete pet photos")
	}

	// 3. Delete the pet from the pets table
	result, err := tx.ExecContext(ctx, "DELETE FROM pets WHERE id = $1", id)
	if err != nil {
		return myerror.WrapError(s.Logger, err, "failed to delete pet")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return myerror.WrapError(s.Logger, err, "failed to get rows affected")
	}
	if rowsAffected == 0 {
		return status.Errorf(codes.NotFound, "pet not found") // Return an error if no pet was deleted
	}

	return nil
}
