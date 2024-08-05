package server

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/f3rcho/grpc-pro/models"
	studentpb "github.com/f3rcho/grpc-pro/proto/student"
	testpb "github.com/f3rcho/grpc-pro/proto/test"
	"github.com/f3rcho/grpc-pro/repository"
)

type TestServer struct {
	repo repository.Repository
	testpb.UnimplementedTestServiceServer
}

func NewTestServer(repo repository.Repository) *TestServer {
	return &TestServer{
		repo: repo,
	}
}

func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repo.GetTest(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &testpb.Test{
		Id:   test.ID,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {
	test := &models.Test{
		ID:   req.GetId(),
		Name: req.GetName(),
	}
	err := s.repo.SetTest(ctx, test)
	if err != nil {
		return nil, err
	}

	return &testpb.SetTestResponse{
		Id:   test.ID,
		Name: test.Name,
	}, nil
}
func (s *TestServer) SetQuestions(stream testpb.TestService_SetQuestionsServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: true})
		}
		if err != nil {
			return err
		}
		question := &models.Question{
			ID:       msg.GetId(),
			Question: msg.GetQuestion(),
			Answer:   msg.GetAnswer(),
			TestID:   msg.GetTestId(),
		}
		err = s.repo.SetQuestions(context.Background(), question)

		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: false})
		}
	}
}

func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: true})
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
			return err
		}
		enrollment := &models.Enrollment{
			StudentID: msg.GetStudentId(),
			TestID:    msg.GetTestId(),
		}
		err = s.repo.SetEnrollment(context.Background(), enrollment)

		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{Ok: false})
		}
	}
}

func (s *TestServer) GetStudentsPerTest(req *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(context.Background(), req.GetTestId())
	if err != nil {
		return err
	}
	for _, student := range students {
		student := &studentpb.Student{
			Id:   student.ID,
			Name: student.Name,
			Age:  student.Age,
		}
		err := stream.Send(student)
		time.Sleep(2 * time.Second) // just to test and see the performance
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TestServer) TakeTest(stream testpb.TestService_TakeTestServer) error {
	for {
		msg, err := stream.Recv()

		if err != io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
		questions, err := s.repo.GetQuestionsPerTest(context.Background(), msg.GetTestId())

		if err != nil {
			return err
		}
		var currentQuestion = &models.Question{}
		i := 0
		lenQuestions := len(questions)
		lenQuestions32 := int32(lenQuestions)
		for {
			if i < lenQuestions {
				currentQuestion = questions[i]
				questionToSend := &testpb.QuestionPerTest{
					Id:       currentQuestion.ID,
					Question: currentQuestion.Question,
					Ok:       false,
					Current:  int32(i + 1),
					Total:    lenQuestions32,
				}

				err := stream.Send(questionToSend)
				if err != nil {
					log.Printf("Error sending question: %v", err)
					return err
				}
				i++

				answer, err := stream.Recv()
				if err != io.EOF {
					return nil
				}

				if err != nil {
					log.Printf("Error receiving answer: %v", err)
					return err
				}

				log.Println("Answer: ", answer.GetAnswer())
				answerModel := &models.Answer{
					TestId:     msg.GetTestId(),
					QuestionId: currentQuestion.ID,
					StudentId:  msg.GetStudentId(),
					Answer:     answer.Answer,
					Correct:    (answer.Answer == currentQuestion.Answer),
				}

				err = s.repo.SetAnswer(context.Background(), answerModel)

				if err != nil {
					return err
				}
			} else {
				questionToSend := &testpb.QuestionPerTest{
					Id:       "",
					Question: "",
					Ok:       true,
					Current:  int32(0),
					Total:    int32(0),
				}

				err := stream.Send(questionToSend)
				if err != io.EOF {
					return nil
				}

				if err != nil {
					return err
				}

				break
			}
		}
	}
}

func (s *TestServer) GetTestScore(ctx context.Context, req *testpb.GetTestScoreRequest) (*testpb.TestScore, error) {
	testScore, err := s.repo.GetTestScore(ctx, req.TestId, req.StudentId)

	if err != nil {
		return nil, err
	}
	return &testpb.TestScore{
		TestId:    testScore.TestID,
		StudentId: testScore.StudentID,
		Ok:        testScore.Ok,
		Ko:        testScore.Ko,
		Total:     testScore.Total,
		Score:     testScore.Score,
	}, nil
}
