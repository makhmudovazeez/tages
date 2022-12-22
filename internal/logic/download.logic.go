package logic

import (
	"context"
	"fmt"
	"github.com/makhmudovazeez/tages/internal/svc"
	"github.com/makhmudovazeez/tages/proto/tages"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"path/filepath"
)

const fileChunk = 8192

type DownloadLogic struct {
	svcCtx svc.ServiceContext
}

func NewDownloadLogic(svcCtx svc.ServiceContext) *DownloadLogic {
	return &DownloadLogic{
		svcCtx: svcCtx,
	}
}

func (l *DownloadLogic) DownloadLogic(in *tages.DownloadRequest, stream tages.Tages_DownloadServer) error {
	var ctx context.Context
	f, err := l.svcCtx.FileModel.FindOne(ctx, in.Id)
	if err != nil {
		return status.Errorf(codes.Canceled, "no such file exists with id = %v")
	}

	file, err := os.ReadFile(fmt.Sprintf("%s/%s", l.svcCtx.Config.Storage, f.Name))
	if err != nil {
		return status.Error(codes.Unknown, "can not read the file")
	}

	if err := stream.Send(&tages.DownloadResponse{
		Data: &tages.DownloadResponse_Mime{
			Mime: filepath.Ext(f.Name),
		},
	}); err != nil {
		return status.Error(codes.Unknown, err.Error())
	}

	for i := 0; i < len(file); i += fileChunk {
		if i+fileChunk > len(file) {
			if err := stream.Send(&tages.DownloadResponse{
				Data: &tages.DownloadResponse_Chunk{
					Chunk: file[i:],
				},
			}); err != nil {
				return status.Error(codes.Unknown, err.Error())
			}
		} else {
			if err := stream.Send(&tages.DownloadResponse{
				Data: &tages.DownloadResponse_Chunk{
					Chunk: file[i:],
				},
			}); err != nil {
				return status.Error(codes.Unknown, err.Error())
			}
		}
	}

	return nil
}
