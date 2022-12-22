package server

import (
	"github.com/makhmudovazeez/tages/internal/logic"
	"github.com/makhmudovazeez/tages/internal/svc"
	"github.com/makhmudovazeez/tages/proto/tages"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type TagesServer struct {
	tages.UnimplementedTagesServer
	svcCtx svc.ServiceContext
}

func NewTagesServer(svcCtx *svc.ServiceContext) *TagesServer {
	return &TagesServer{
		svcCtx: *svcCtx,
	}
}

func (t *TagesServer) UploadFile(stream tages.Tages_UploadFileServer) error {
	l := logic.NewUploadFileLogic(t.svcCtx)
	res, err := l.UploadFileLogic(stream)
	if err != nil {
		return err
	}

	if err = stream.SendAndClose(res); err != nil {
		log.Printf("cannot send response: %v", err)
		return status.Errorf(codes.Unknown, "cannot send response: %w", err)
	}
	return nil
}
