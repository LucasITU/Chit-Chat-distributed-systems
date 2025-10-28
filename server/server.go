package main

import (
	proto "ChitChat/grpc"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
)

var streams map[string]chan *proto.Chat = make(map[string]chan *proto.Chat)
var timestamp chan int64

type Server struct {
	proto.UnimplementedChitChatServiceServer
}

func (s *Server) Join(in *proto.User, stream proto.ChitChatService_JoinServer) error {
	channel := make(chan *proto.Chat, 10)
	streams[in.Name] = channel
	temp_timestamp := <-timestamp + 1
	timestamp <- temp_timestamp
	log_message(temp_timestamp, "Client", "Join", in.Name, nil)

	sendToClients(&proto.Chat{
		Timestamp: temp_timestamp,
		Author:    "Server",
		Message:   fmt.Sprintf("Participant %s joined Chit Chat at logical time %d", in.Name, temp_timestamp),
	})

	for {
		select {
		case chat := <-channel:
			temp_timestamp := <-timestamp
			timestamp <- temp_timestamp

			log_message(temp_timestamp, "Server", "Delivery", in.Name, nil)
			stream.Send(chat)
		case <-time.After(1 * time.Second):
			if streams[in.Name] != channel {
				temp_timestamp := <-timestamp + 1
				timestamp <- temp_timestamp

				log_message(temp_timestamp, "Client", "Leave", in.Name, nil)
				sendToClients(&proto.Chat{
					Timestamp: temp_timestamp,
					Author:    "Server",
					Message:   fmt.Sprintf("Participant %s left Chit Chat at logical time %d", in.Name, temp_timestamp),
				})
				return nil
			}
		}
	}
}

func (s *Server) Leave(ctx context.Context, in *proto.User) (*proto.Empty, error) {
	delete(streams, in.Name)

	temp_timestamp := <-timestamp
	timestamp <- temp_timestamp

	return &proto.Empty{}, nil
}

func (s *Server) SendChat(ctx context.Context, in *proto.Chat) (*proto.Empty, error) {
	temp_timestamp := max(<-timestamp, in.Timestamp) + 1
	timestamp <- temp_timestamp

	sendToClients(in)

	return &proto.Empty{}, nil
}

func sendToClients(in *proto.Chat) {
	temp_timestamp := <-timestamp
	timestamp <- temp_timestamp + 1

	log_message(temp_timestamp, "Client", "Recieve", in.Author, in)

	in.Timestamp = temp_timestamp + 1

	for _, v := range streams {
		v <- in
	}
}

func main() {
	file, err := os.Create("events.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	server := &Server{}

	timestamp = make(chan int64, 1)
	timestamp <- 0

	go server.start_server()
	log_message(0, "Server", "Start", "", nil)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		log_message(<-timestamp+1, "Server", "Stop", "", nil)
		os.Exit(0)
	}
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

func log_message(timestamp int64, component string, event_type string, identifier string, message any) {
	if message == nil {
		message = ""
	}

	log.Printf(
		"[Timestamp:%d] [Component:%s] [EventType:%s] [Identifier:%s] %s",
		timestamp,
		component,
		event_type,
		identifier,
		message,
	)
}
