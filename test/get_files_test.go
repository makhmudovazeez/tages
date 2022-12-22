package test

import (
	"context"
	"github.com/makhmudovazeez/tages/proto/tages"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
	"time"
)

func TestGetFiles(t *testing.T) {
	t.Parallel()

	conn, err := grpc.Dial("127.0.0.1:8181", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client := tages.NewTagesClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if _, err := client.GetFiles(ctx, &emptypb.Empty{}); err != nil {
		t.Error(err)
	}
}
