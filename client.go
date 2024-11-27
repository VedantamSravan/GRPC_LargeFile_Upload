package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc_ex1/pb"
)

const (
	//it connects server ip address
	serverAddr = "10.91.170.144:50051"
)

func main() {
	// Load the client certificate
	certFile := "server.crt"
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("failed to load client certificate: %v", err)
	}

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileTransferServiceClient(conn)

	// Open the file
	file, err := os.Open("examplefile.txt")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a stream
	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatalf("could not upload file: %v", err)
	}

	// Send file chunks
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to read file: %v", err)
		}

		// Send the chunk
		if err := stream.Send(&pb.FileChunk{
			Filename: "examplefile.txt",
			Content:  buffer[:n],
		}); err != nil {
			log.Fatalf("failed to send chunk: %v", err)
		}
	}

	// Receive the status
	status, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive status: %v", err)
	}
	fmt.Printf("Upload status: %s\n", status.Message)
}

