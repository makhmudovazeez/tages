package main

import (
	"context"
	"github.com/makhmudovazeez/tages/proto/tages"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	f  string = "2.jpg"
	ID string
)

const fileChunk = 8192

func TestUploadFile(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:8181", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := tages.NewTagesClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	file, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	requests := []*tages.UploadFileRequest{
		{
			Data: &tages.UploadFileRequest_Mime{
				Mime: filepath.Ext(f),
			},
		},
	}

	for i := 0; i < len(file); i += fileChunk {
		if i+fileChunk > len(file) {
			requests = append(requests, &tages.UploadFileRequest{
				Data: &tages.UploadFileRequest_Chunk{
					Chunk: file[i:],
				},
			})
		} else {
			requests = append(requests, &tages.UploadFileRequest{
				Data: &tages.UploadFileRequest_Chunk{
					Chunk: file[i : i+fileChunk],
				},
			})
		}
	}

	stream, err := client.UploadFile(ctx)
	if err != nil {
		log.Fatalf("can not create stream: %v", err)
	}

	for _, req := range requests {
		if err := stream.Send(req); err != nil {
			t.Fatal(err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}

	ID = resp.Id
}

func TestGetFiles(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:8181", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := tages.NewTagesClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if _, err := client.GetFiles(ctx, &emptypb.Empty{}); err != nil {
		t.Fatal(err)
	}
}

func TestDownload(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:8181", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := tages.NewTagesClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	stream, err := client.Download(ctx, &tages.DownloadRequest{
		Id: ID,
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	_, err = stream.Recv()
	if err != nil {
		t.Fatal(err)
	}

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			log.Println("received all chunks")
			break
		}
		if err != nil {
			t.Fatal(err)
			break
		}
	}

	if err := stream.CloseSend(); err != nil {
		t.Fatal(err)
	}
}
