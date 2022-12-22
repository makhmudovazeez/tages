package logic

import (
	"bytes"
	"context"
	"fmt"
	"github.com/makhmudovazeez/tages/internal/services"
	"github.com/makhmudovazeez/tages/internal/svc"
	"github.com/makhmudovazeez/tages/proto/tages"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
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
	firstReq, err := stream.Recv()
	var errStr string
	if err != nil {
		errStr = fmt.Sprintf("cannot receive file info: %v", err)
		log.Println(errStr)
		return nil, status.Error(codes.Unknown, errStr)
	}

	fileData := bytes.Buffer{}

	for {
		otherReq, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			errStr = fmt.Sprintf("cannot receive chunk data 1: %v", err)
			log.Println(errStr)
			return nil, status.Error(codes.Unknown, errStr)
		}

		chunk := otherReq.GetChunk()

		if len(chunk) == 0 {
			errStr = fmt.Sprintf("cannot receive chunk data 2: %v", err)
			log.Println(errStr)
			return nil, status.Error(codes.Unknown, errStr)
		}

		if _, err := fileData.Write(chunk); err != nil {
			errStr = fmt.Sprintf("cannot write chund data: %v", err)
			log.Println(errStr)
			return nil, status.Errorf(codes.Internal, errStr)
		}
	}

	dfs := services.NewFileStore(l.svcCtx.Config.Storage)
	id, err := dfs.Save(firstReq.GetMime(), fileData)
	if err != nil {
		errStr := fmt.Sprintf("cannot save file to storage: %v", err)
		log.Println(errStr)
		return nil, status.Errorf(codes.Internal, errStr)
	}

	var ctx context.Context
	if err := l.svcCtx.FileModel.Save(ctx, id, firstReq.GetMime()); err != nil {
		errStr := fmt.Sprintf("cannot save file to db: %v", err)
		log.Println(errStr)
		return nil, status.Errorf(codes.Internal, errStr)
	}

	res := &tages.UploadFileResponse{
		Id: id,
	}

	return res, nil
}
