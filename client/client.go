package main

import (
	proto "ChitChat/grpc"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var timestamp chan int64 //use channel for this var

func main() {
	author := os.Args[1]

	conn, err := grpc.NewClient("localhost:5050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Not working")
	}

	client := proto.NewChitChatServiceClient(conn)

	user := proto.User{Name: author}

	timestamp = make(chan int64, 1)
	timestamp <- 0

	stream, _ := client.Join(context.Background(), &user)
	go recieve(stream)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\r> ")
		text, _ := reader.ReadString('\n')

		temp_timestamp := <-timestamp + 1
		timestamp <- temp_timestamp

		if strings.TrimSpace(text) == "/leave" {
			client.Leave(context.Background(), &user)
			return
		} else {
			chat := proto.Chat{Timestamp: temp_timestamp, Author: author, Message: text[:min(len(text)-1, 128)]}
			client.SendChat(context.Background(), &chat)
		}

	}
}

func recieve(stream grpc.ServerStreamingClient[proto.Chat]) {
	for {
		chat, err := stream.Recv()

		if err == io.EOF {
			fmt.Println("\rStream closed")
			os.Exit(0)
		}

		temp_timestamp := max(<-timestamp, chat.Timestamp) + 1
		timestamp <- temp_timestamp

		fmt.Println("\r[" + strconv.Itoa(int(chat.Timestamp)) + "] " + chat.Author + ": " + chat.Message)
		fmt.Print("\r> ")
	}
}
