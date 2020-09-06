package apiservice

import (
	"context"
	"github.com/amit/file-download-manager/internal/download"
	"github.com/amit/file-download-manager/internal/user"
	"github.com/amit/file-download-manager/pkg/logger"
	rpcPkg "github.com/amit/file-download-manager/rpc"
	"github.com/twitchtv/twirp"
)

type Server struct{}

func (s *Server) Sample(ctx context.Context, request *rpcPkg.SampleRequest) (*rpcPkg.SampleResponse, error) {
	resp := &rpcPkg.SampleResponse{
		IsSuccess: true,
		Message:   "working",
	}
	return resp, nil
}

func (s *Server) AddUser(ctx context.Context, request *rpcPkg.AddUserRequest) (*rpcPkg.AddUserResponse, error) {
	logger.LogStatement("ADD USER REQUEST :: ", request)
	// TODO : Input Validation
	newUser := &user.User{
		Username:  			request.Username,
		Email:     			request.Email,
	}
	addErr := user.AddUser(ctx, newUser)
	if addErr != nil {
		if addErr.Error() == "username_already_exists" || addErr.Error() == "multiple_users_found_for_username" {
			return &rpcPkg.AddUserResponse{}, twirp.NewError(twirp.Internal, "duplicate username")
		} else {
			return &rpcPkg.AddUserResponse{}, twirp.NewError(twirp.Internal, "internal server error")
		}
	}
	response := &rpcPkg.AddUserResponse{
		Username: newUser.Username,
		SecretKey: newUser.SecretKey,
	}
	logger.LogStatement("ADD USER RESPONSE :: ", response)
	return response, nil
}

func (s *Server) Download(ctx context.Context, request *rpcPkg.DownloadRequest) (*rpcPkg.DownloadResponse, error) {
	logger.LogStatement("DOWNLOAD REQUEST :: ", request)
	downloadId, fileIdVsFileUrl, createRequestErr := download.CreateDownloadRequest(ctx, request.DownloadType.String(), request.Files)
	if createRequestErr != nil {
		return nil, createRequestErr
	}
	if request.DownloadType == rpcPkg.DownloadType_SYNC {
		downloadErr := download.DownloadFiles(ctx, downloadId, fileIdVsFileUrl)
		if downloadErr != nil {
			return nil, downloadErr
		}
	} else {
		go func() {
			downloadErr := download.DownloadFiles(ctx, downloadId, fileIdVsFileUrl)
			if downloadErr != nil {
				errorMsg := "Error downloading file for download Id :: " + downloadId
				logger.LogStatement(errorMsg, downloadErr)
			}
		}()
	}
	response := &rpcPkg.DownloadResponse {
		Id:     downloadId,
		Files:  request.Files,
	}
	logger.LogStatement("DOWNLOAD RESPONSE :: ", response)
	return response, nil
}