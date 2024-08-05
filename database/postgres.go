package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/f3rcho/grpc-pro/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) GetStudent(ctx context.Context, ID string) (*models.Student, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id = $1", ID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var student = models.Student{}
	for rows.Next() {
		if err = rows.Scan(&student.ID, &student.Name, &student.Age); err == nil {
			return &student, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &student, nil
}

func (repo *PostgresRepository) SetStudent(ctx context.Context, student *models.Student) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO students(id, name, age) VALUES($1, $2, $3)", student.ID, student.Name, student.Age)
	return err
}

func (repo *PostgresRepository) GetTest(ctx context.Context, ID string) (*models.Test, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name FROM tests WHERE id = $1", ID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var test = models.Test{}
	for rows.Next() {
		if err = rows.Scan(&test.ID, &test.Name); err == nil {
			return &test, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &test, nil
}

func (repo *PostgresRepository) SetTest(ctx context.Context, test *models.Test) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO tests(id, name) VALUES($1, $2)", test.ID, test.Name)
	return err
}
func (repo *PostgresRepository) SetQuestions(ctx context.Context, question *models.Question) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO questions(id, answer, question, test_id) VALUES($1, $2, $3, $4)", question.ID, question.Answer, question.Question, question.TestID)
	return err
}
func (repo *PostgresRepository) SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO enrollments(student_id, test_id) VALUES($1, $2)", enrollment.StudentID, enrollment.TestID)
	return err
}

func (repo *PostgresRepository) GetStudentsPerTest(ctx context.Context, testID string) ([]*models.Student, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id IN (SELECT student_id FROM enrollments WHERE test_id = $1)", testID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var students []*models.Student

	for rows.Next() {
		var student = models.Student{}
		if err = rows.Scan(&student.ID, &student.Name, &student.Age); err == nil {
			students = append(students, &student)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

func (repo *PostgresRepository) GetQuestionsPerTest(ctx context.Context, testID string) ([]*models.Question, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, question FROM questions WHERE test_id = $1", testID)

	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var questions []*models.Question

	for rows.Next() {
		var question = models.Question{}
		if err = rows.Scan(&question.ID, &question.Question); err == nil {
			questions = append(questions, &question)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return questions, nil
}

func (repo *PostgresRepository) SetAnswer(ctx context.Context, answer *models.Answer) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO answers(student_id, test_id, question_id, answer, correct) VALUES($1, $2, $3, $4, $5)", answer.StudentId, answer.TestId, answer.QuestionId, answer.Answer, answer.Correct)
	return err
}

func (repo *PostgresRepository) GetTestScore(ctx context.Context, testID, studentID string) (*models.TestScore, error) {
	rows, err := repo.db.QueryContext(ctx, "select correct FROM answers WHERE test_id = $1 AND student_id = $2", testID, studentID)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	var answer = models.Answer{}
	var testScore = models.TestScore{
		TestID:    testID,
		StudentID: studentID,
	}

	for rows.Next() {
		err := rows.Scan(&answer.Correct)

		if err != nil {
			return nil, err
		}
		testScore.Total += 1

		if answer.Correct {
			testScore.Ok += 1
		} else {
			testScore.Ko += 1
		}

	}
	testScore.Score = testScore.Ok * 10 / testScore.Total

	return &testScore, nil
}
