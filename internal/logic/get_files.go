package logic

import (
	"context"
	"github.com/makhmudovazeez/tages/internal/svc"
	"github.com/makhmudovazeez/tages/proto/tages"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type GetFilesLogic struct {
	svcCtx svc.ServiceContext
}

func NewGetFilesLogic(svcCtx svc.ServiceContext) *GetFilesLogic {
	return &GetFilesLogic{
		svcCtx: svcCtx,
	}
}

func (l *GetFilesLogic) GetFilesLogic(ctx context.Context) (*tages.GetFileResponse, error) {
	files, err := l.svcCtx.FileModel.GetAll(ctx)
	if err != nil {
		log.Println(err)
		return nil, status.Error(codes.Unimplemented, err.Error())
	}

	var resp tages.GetFileResponse
	for _, f := range files {
		resp.Files = append(resp.Files, &tages.File{
			Id:        f.Id,
			Name:      f.Name,
			CreatedAt: f.CreatedAt.String(),
			UpdatedAt: f.UpdatedAt.String(),
		})
	}

	return &resp, nil
}
