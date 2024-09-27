package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "grpc/exp2" // assuming your proto file is in the same directory
)

type wordService struct {
	pb.UnimplementedWordServiceServer // Add this line to embed the interface
}

func (s *wordService) AddWord(ctx context.Context, req *pb.AddWordRequest) (*pb.AddWordResponse, error) {
	// implement your logic to add a word here
	// for demonstration purposes, I'll just return a success response
	return &pb.AddWordResponse{Result: "Word added successfully"}, nil
}

func main() {
	fmt.Println("Starting server...")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterWordServiceServer(srv, &wordService{})

	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
