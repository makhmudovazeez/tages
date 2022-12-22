package server

import (
	"context"
	"github.com/makhmudovazeez/tages/internal/logic"
	"github.com/makhmudovazeez/tages/internal/svc"
	"github.com/makhmudovazeez/tages/proto/tages"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type TagesServer struct {
	tages.UnimplementedTagesServer
	svcCtx svc.ServiceContext
}

var uploadFileLimit = make(chan int, 10)
var getFilesLimit = make(chan int, 100)
var downloadLimit = make(chan int, 10)

func NewTagesServer(svcCtx *svc.ServiceContext) *TagesServer {
	return &TagesServer{
		svcCtx: *svcCtx,
	}
}

func (t *TagesServer) UploadFile(stream tages.Tages_UploadFileServer) error {
	uploadFileLimit <- time.Now().Second()
	defer func() {
		<-uploadFileLimit
	}()

	errChan := make(chan error)

	go func() {
		l := logic.NewUploadFileLogic(t.svcCtx)
		res, err := l.UploadFileLogic(stream)
		if err != nil {
			errChan <- err
		}

		if err = stream.SendAndClose(res); err != nil {
			log.Printf("cannot send response: %v", err)
			errChan <- status.Errorf(codes.Unknown, "cannot send response: %w", err)
		}
		errChan <- nil
	}()

	return <-errChan
}

func (t *TagesServer) GetFiles(ctx context.Context, in *emptypb.Empty) (*tages.GetFileResponse, error) {
	getFilesLimit <- time.Now().Second()
	defer func() {
		<-getFilesLimit
	}()

	errChan := make(chan error)
	getFileChan := make(chan *tages.GetFileResponse)

	go func() {
		l := logic.NewGetFilesLogic(t.svcCtx)
		resp, err := l.GetFilesLogic(ctx)
		if err != nil {
			errChan <- err
		}
		getFileChan <- resp
	}()

	select {
	case err := <-errChan:
		return nil, err
	case fileResp := <-getFileChan:
		return fileResp, nil
	}
}

func (t *TagesServer) Download(in *tages.DownloadRequest, stream tages.Tages_DownloadServer) error {
	downloadLimit <- time.Now().Second()
	defer func() {
		<-downloadLimit
	}()

	errChan := make(chan error)

	go func() {
		l := logic.NewDownloadLogic(t.svcCtx)
		errChan <- l.DownloadLogic(in, stream)
	}()

	return <-errChan
}
