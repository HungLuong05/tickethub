package grpc

import (
	"context"
	"log"
	"net"
	"sync"

	"tickethub.com/auth/proto"
	// "tickethub.com/auth/config"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedEventPermServer
}

func (s *server) AddEventPerm(ctx context.Context, req *proto.AddEventPermRequest) (*proto.AddEventPermResponse, error) {
	log.Printf("Receive event id: %v\n", req.GetEventId())
	log.Printf("Receive user id: %v\n", req.GetUserId())

	// query := "INSERT INTO event_perms(event_id, user_id) VALUES($1, $2)"
	// err := pg.DB.Exec(query, req.GetEventId(), req.GetUserId())
	// if err != nil {
	// 	log.Printf("Failed to execute query: %v", err)
	// 	return nil, err
	// }

	return &proto.AddEventPermResponse{Message: "Successful"}, nil
}

func (s *server) DeleteEventPerm(ctx context.Context, req *proto.DeleteEventPermRequest) (*proto.DeleteEventPermResponse, error) {
	log.Printf("Receive event id: %v\n", req.GetEventId())
	log.Printf("Receive user id: %v\n", req.GetUserId())

	// query := "DELETE FROM event_perms WHERE event_id = $1 AND user_id = $2"
	// err := pg.DB.Exec(query, req.GetEventId(), req.GetUserId())
	// if err != nil {
	// 	log.Printf("Failed to execute query: %v", err)
	// 	return nil, err
	// }

	return &proto.DeleteEventPermResponse{Message: "Successful"}, nil
}

func StartGrpcServer(wg *sync.WaitGroup) {
	defer wg.Done()
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Printf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterEventPermServer(grpcServer, &server{})
	log.Println("Server is running on port 50051...", lis.Addr().String())
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("Failed to serve: %v", err)
	}
}
