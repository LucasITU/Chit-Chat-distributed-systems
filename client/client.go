package main

import (
	proto "ChitChat/grpc"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	author := os.Args[1]

	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewChitChatServiceClient(conn)

	user := proto.User{Name: author}

	stream, _ := client.Join(context.Background(), &user)
	go recieve(stream)

	select {}
}

func recieve(stream grpc.ServerStreamingClient[proto.Chat]) {
out:
	for {
		chat, err := stream.Recv()

		if err == io.EOF {
			break out
		}

		fmt.Println("[" + strconv.Itoa(int(chat.Timestamp)) + "] " + chat.Author + ": " + chat.Message)
	}
}
