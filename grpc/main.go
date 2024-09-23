package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "exp"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "ghbdtn"
	dbname   = "users"
)

func main() {
	// Connect to the database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new gRPC server
	srv := grpc.NewServer()
	pb.RegisterDatabaseServiceServer(srv, &databaseService{db: db})

	// Start the gRPC server
	log.Println("gRPC server listening on port 50051")
	if err := srv.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

type databaseService struct {
	db *sql.DB
}

func (s *databaseService) GetIdByName(ctx context.Context, req *pb.Name) (*pb.Id, error) {
	// Execute a query to retrieve the ID by name
	var id int32
	err := s.db.QueryRow("SELECT id FROM users WHERE name = $1", req.Value).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &pb.Id{Value: id}, nil
}
