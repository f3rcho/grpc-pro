package server

import (
	"context"
	"io"

	"github.com/f3rcho/grpc-pro/models"
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
