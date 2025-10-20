package main

import (
	proto "ChitChat/grpc"
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedChitChatServiceServer
}

func (s *Server) Join(in *proto.User, stream proto.ChitChatService_JoinServer) error {
	for {
		chat := proto.Chat{ Author: in.Name, Message: "Joined the server"}
		stream.Send(&chat)
		time.Sleep(1 * time.Second)
	}
}

func (s *Server) SendChat(ctx context.Context, in *proto.Chat) (*proto.Empty, error) {
	return &proto.Empty{}, nil;
}

func main() {
	server := &Server{}

	server.start_server()
}

func (s *Server) start_server() {
	grpcServer := grpc.NewServer()
	listener, err := net.Listen("tcp", ":5050")
	if err != nil {
		log.Fatalf("Did not work")
	}

	proto.RegisterChitChatServiceServer(grpcServer, s)
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("Did not work")
	}

}
