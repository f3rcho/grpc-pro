package repository

import (
	"context"

	"github.com/f3rcho/grpc-pro/models"
)

type Repository interface {
	GetStudent(ctx context.Context, ID string) (*models.Student, error)
	SetStudent(ctx context.Context, student *models.Student) error
	GetTest(ctx context.Context, ID string) (*models.Test, error)
	SetTest(ctx context.Context, test *models.Test) error
}

var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

func GetStudent(ctx context.Context, ID string) (*models.Student, error) {
	return implementation.GetStudent(ctx, ID)
}

func SetStudent(ctx context.Context, student *models.Student) error {
	return implementation.SetStudent(ctx, student)
}

func GetTest(ctx context.Context, ID string) (*models.Test, error) {
	return implementation.GetTest(ctx, ID)
}

func SetTest(ctx context.Context, test *models.Test) error {
	return implementation.SetTest(ctx, test)
}
