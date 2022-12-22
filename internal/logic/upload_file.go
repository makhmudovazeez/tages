package logic

import (
	"github.com/makhmudovazeez/tages/internal/svc"
	"github.com/makhmudovazeez/tages/proto/tages"
)

type UploadFileLogic struct {
	svcCtx svc.ServiceContext
}

func NewUploadFileLogic(svcCtx svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		svcCtx: svcCtx,
	}
}

func (l *UploadFileLogic) UploadFileLogic(stream tages.Tages_UploadFileServer) (*tages.UploadFileResponse, error) {

	return nil, nil
}
