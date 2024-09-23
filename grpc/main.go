package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "grpc/exp"
)

type databaseService struct {
	pb.UnimplementedDatabaseServiceServer
}

func (s *databaseService) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	// implement your logic here
	return &pb.GetDataResponse{Data: "some data"}, nil
}

func main() {
	srv := grpc.NewServer()
	pb.RegisterDatabaseServiceServer(srv, &databaseService{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port 50051")
	srv.Serve(lis)
}
