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
	questions, err := s.repo.GetQuestionsPerTest(context.Background(), "t1")
	if err != nil {
		return err
	}

	i := 0
	var currentQuestion = &models.Question{}
	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}

		if i <= len(questions) {
			questionToSend := &testpb.Question{
				Id:       currentQuestion.ID,
				Question: currentQuestion.Question,
			}
			err := stream.Send(questionToSend)
			if err != nil {
				return err
			}
			i++
		}
		answer, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Error sending question: %v", err)
			return nil
		}
		if err != nil {
			log.Printf("Error receiving answer: %v", err)
			return err
		}
		log.Println("Answer for question:", currentQuestion.Question, "is", answer.GetAnswer())
	}
}
