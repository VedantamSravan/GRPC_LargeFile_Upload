package main

import (
//	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc_ex1/pb"
)

const (
	port = "0.0.0.0:50051"
)

type server struct{
	pb.UnimplementedFileTransferServiceServer
}

func (s *server) Upload(stream pb.FileTransferService_UploadServer) error {
	var filename string
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.UploadStatus{
				Message: "Upload completed",
				Code:    0,
			})
		}
		if err != nil {
			return err
		}

		if filename == "" {
			filename = chunk.Filename
			outFile, err := os.Create(filepath.Join("./uploads", filename))
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = outFile.Write(chunk.Content)
			if err != nil {
				return err
			}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Load the server's certificate and private key
	certFile := "server.crt"
	keyFile := "server.key"
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to generate credentials: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterFileTransferServiceServer(grpcServer, &server{})

	log.Printf("Server listening on %v", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

