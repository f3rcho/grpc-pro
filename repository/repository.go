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
	SetQuestions(ctx context.Context, question *models.Question) error
	SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	GetStudentsPerTest(ctx context.Context, testID string) ([]*models.Student, error)
	GetQuestionsPerTest(ctx context.Context, testID string) ([]*models.Question, error)
	GetTestScore(ctx context.Context, testID, studentID string) (*models.TestScore, error)
	SetAnswer(ctx context.Context, answer *models.Answer) error
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
func SetQuestions(ctx context.Context, question *models.Question) error {
	return implementation.SetQuestions(ctx, question)
}
func SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	return implementation.SetEnrollment(ctx, enrollment)
}

func GetStudentsPerTest(ctx context.Context, testID string) ([]*models.Student, error) {
	return implementation.GetStudentsPerTest(ctx, testID)
}

func GetQuestionsPerTest(ctx context.Context, testID string) ([]*models.Question, error) {
	return implementation.GetQuestionsPerTest(ctx, testID)
}

func GetTestScore(ctx context.Context, testID, studentID string) (*models.TestScore, error) {
	return implementation.GetTestScore(ctx, testID, studentID)
}

func SetAnswer(ctx context.Context, answer *models.Answer) error {
	return implementation.SetAnswer(ctx, answer)
}
