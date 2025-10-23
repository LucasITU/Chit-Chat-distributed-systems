package main

import (
	proto "ChitChat/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

var streams map[string]chan *proto.Chat = make(map[string]chan *proto.Chat)
var timestamp int64

type Server struct {
	proto.UnimplementedChitChatServiceServer
}

func (s *Server) Join(in *proto.User, stream proto.ChitChatService_JoinServer) error {
	channel := make(chan *proto.Chat, 10)
	streams[in.Name] = channel
	timestamp++
	log.Printf("Logical time: %d Component: [Server] Event type: [Join] Identifier: %s", timestamp, in.Name)

	s.SendChat(context.Background(), &proto.Chat{
		Timestamp: timestamp,
		Author:    "Server",
		Message:   fmt.Sprintf("Participant %s joined Chit Chat at logical time %d", in.Name, timestamp),
	})

	for {
		select {
		case chat := <-channel:
			log.Printf("Logical time: %d Component: [Client] Event type: [Delivery] Identifier: %s", timestamp, in.Name)
			stream.Send(chat)
		case <-time.After(1 * time.Second):
			if streams[in.Name] != channel {
				timestamp++
				s.SendChat(context.Background(), &proto.Chat{
					Timestamp: timestamp,
					Author:    "Server",
					Message:   fmt.Sprintf("Participant %s left Chit Chat at logical time %d", in.Name, timestamp),
				})
				return nil
			}
		}
	}
}

func (s *Server) Leave(ctx context.Context, in *proto.User) (*proto.Empty, error) {
	delete(streams, in.Name)
	log.Println("[Server]", in.Name+" has left") //maybe more info
	return &proto.Empty{}, nil
}

func (s *Server) SendChat(ctx context.Context, in *proto.Chat) (*proto.Empty, error) {
	timestamp = max(timestamp, in.Timestamp)
	log.Println("[Client]", in)

	for _, v := range streams {
		v <- in
	}
	return &proto.Empty{}, nil
}

func main() {
	server := &Server{}

	server.start_server()
	log.Println("Server has started")
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

//log messages
//concurrency bugs??? channels?
//is logical time correct?
//is anything correct?
//Should client log?
