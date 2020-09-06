package download

import (
	"context"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/amit/file-download-manager/internal/user"
	"github.com/amit/file-download-manager/internal/util"
	"github.com/amit/file-download-manager/pkg/db"
	"github.com/amit/file-download-manager/pkg/logger"
)

func incrementUniqueId(fileUniqueId string) string {
	id, _ := strconv.Atoi(fileUniqueId)
	id += 1
	fileUniqueId = strconv.Itoa(id)
	return fileUniqueId
}

func CreateDownloadRequest(ctx context.Context, downloadType string, files []string) (string, map[string]string, error) {
	uniqueId := util.GenerateUniqueId()
	downloadId := DOWNLOAD_REQUEST_ID_PREFIX + uniqueId
	downloadRequestDbObject := &DownloadRequest{
		Id:              downloadId,
		RequestedUserId: ctx.Value(user.USER_ID_CONTEXT_KEY).(string),
		DownloadType:    downloadType,
		DownloadStatus:  DOWNLOAD_STATUS_NEW,
		RequestedAt:     time.Now().Unix(),
	}
	dbCtx, beginErr := db.RepoClient.Begin(ctx)
	if beginErr != nil {
		return "", nil, beginErr
	}
	createErr := db.RepoClient.Create(dbCtx, DOWNLOAD_REQUEST_TABLE_NAME, downloadRequestDbObject)
	if createErr != nil {
		db.RepoClient.Rollback(dbCtx)
		return "", nil, createErr
	}
	fileIdVsFileName := make(map[string]string)
	fileUniqueId := util.GenerateUniqueId()
	for _, file := range files {
		fileUniqueId = incrementUniqueId(fileUniqueId)
		fileId := DOWNLOAD_FILES_ID_PREFIX + fileUniqueId
		fileDetailsDbObject := &FileDetails{
			Id:                fileId,
			DownloadRequestId: downloadId,
			FileName:          file,
			Status:            FILE_DOWNLOAD_STATUS_NEW,
			CreatedAt:         time.Now().Unix(),
		}
		fileIdVsFileName[fileId] = file
		createErr = db.RepoClient.Create(dbCtx, DOWNLOAD_FILES_TABLE_NAME, fileDetailsDbObject)
		if createErr != nil {
			db.RepoClient.Rollback(dbCtx)
			return "", nil, createErr
		}
	}
	commitErr := db.RepoClient.Commit(dbCtx)
	if commitErr != nil {
		db.RepoClient.Rollback(dbCtx)
		return "", nil, commitErr
	}
	return downloadId, fileIdVsFileName, nil
}

func DownloadFiles(ctx context.Context, downloadId string, fileIdVsFileUrl map[string]string) error {
	downloadRequestDbObj := DownloadRequest{}
	findErr := db.RepoClient.Find(ctx, &downloadRequestDbObj, downloadId, DOWNLOAD_REQUEST_TABLE_NAME)
	if findErr != nil {
		return findErr
	}
	updateAttributes := map[string]interface{}{
		"status": DOWNLOAD_STATUS_IN_PROGRESS,
	}
	updateErr := db.RepoClient.Update(ctx, &downloadRequestDbObj, updateAttributes)
	if updateErr != nil {
		return updateErr
	}
	downloadFilesErr := DownloadAndZipFiles(ctx, downloadId, fileIdVsFileUrl)
	if downloadFilesErr != nil {
		updateAttributes := map[string]interface{}{
			"status": DOWNLOAD_STATUS_FAILED,
		}
		updateErr := db.RepoClient.Update(ctx, &downloadRequestDbObj, updateAttributes)
		if updateErr != nil {
			return updateErr
		}
		return downloadFilesErr
	}
	updateAttributes = map[string]interface{}{
		"status": DOWNLOAD_STATUS_SUCCESS,
	}
	updateErr = db.RepoClient.Update(ctx, &downloadRequestDbObj, updateAttributes)
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func DownloadAndZipFiles(ctx context.Context, downloadId string, fileIdVsFileUrl map[string]string) error {
	downloadsDirectory := util.GetDownloadsDirectory()
	workingDirectory := path.Join(downloadsDirectory, downloadId)
	createDirErr := os.Mkdir(workingDirectory, 0755)
	if createDirErr != nil {
		logger.LogStatement("Create directory ERROR :: ", createDirErr)
		return createDirErr
	}
	for fileId, url := range fileIdVsFileUrl {
		fileDownloadErr := DownloadFile(ctx, workingDirectory, fileId, url)
		logger.LogStatement("Create directory ERROR :: ", fileDownloadErr)
	}
	return nil
}

func DownloadFile(ctx context.Context, downloadingDirectory string, fileId string, fileUrl string) error {
	fileDbObj := FileDetails{}
	findErr := db.RepoClient.Find(ctx, &fileDbObj, fileId, DOWNLOAD_FILES_TABLE_NAME)
	if findErr != nil {
		return findErr
	}
	updateAttributes := map[string]interface{}{
		"status": FILE_DOWNLOAD_STATUS_IN_PROGRESS,
	}
	updateErr := db.RepoClient.Update(ctx, &fileDbObj, updateAttributes)
	if updateErr != nil {
		return updateErr
	}
	file, httpGetErr := http.Get(fileUrl)
	if httpGetErr != nil {

	}
	defer file.Body.Close()
	return nil
}
