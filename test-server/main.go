package main

import (
	"log"
	"net"

	"github.com/f3rcho/grpc-pro/database"
	testpb "github.com/f3rcho/grpc-pro/proto/test"
	"github.com/f3rcho/grpc-pro/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	list, err := net.Listen("tcp", ":5070")

	if err != nil {
		log.Fatalf("Error listening: %s", err.Error())
	}

	repo, err := database.NewPostgresRepository("postgres://root:123456@localhost:5432/mydb?sslmode=disable")

	server := server.NewTestServer(repo)

	if err != nil {
		log.Fatalf("Error creating repository: %s", err.Error())
	}

	s := grpc.NewServer()
	testpb.RegisterTestServiceServer(s, server)

	reflection.Register(s)
	log.Println("Server running on :5070")

	if err := s.Serve(list); err != nil {
		log.Fatalf("Error serving: %s", err.Error())
	}
}
