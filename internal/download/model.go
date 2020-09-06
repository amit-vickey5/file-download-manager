package download

const (
	DOWNLOAD_REQUEST_PUBLIC_PREFIX 		= "download_"
	DOWNLOAD_REQUEST_ID_PREFIX 			= "DOWNLOAD"
	DOWNLOAD_REQUEST_TABLE_NAME 		= "download_request"

	DOWNLOAD_FILES_TABLE_NAME 			= "download_files"
	DOWNLOAD_FILES_ID_PREFIX 			= "FILE"

	DOWNLOAD_STATUS_NEW 				= "new"
	DOWNLOAD_STATUS_IN_PROGRESS 		= "in_progress"
	DOWNLOAD_STATUS_FAILED 				= "failed"
	DOWNLOAD_STATUS_SUCCESS 			= "success"

	FILE_DOWNLOAD_STATUS_NEW 			= "new"
	FILE_DOWNLOAD_STATUS_IN_PROGRESS 	= "in_progress"
	FILE_DOWNLOAD_STATUS_FAILED 		= "failed"
	FILE_DOWNLOAD_STATUS_SUCCESS 		= "success"
)

type DownloadRequest struct {
	Id 					string	`gorm:"string(255);not null" json:"id" tag_name:"id"`
	RequestedUserId		string	`gorm:"string(255);not null" json:"requested_user_id" tag_name:"requested_user_id"`
	DownloadType		string	`gorm:"string(14);not null" json:"download_type" tag_name:"download_type"`
	DownloadStatus		string	`gorm:"string(14);not null" json:"download_status" tag_name:"download_status"`
	ZipFileName			string	`gorm:"string(255)" sql:"DEFAULT:null" json:"zip_file_name" tag_name:"zip_file_name"`
	FailureReason		string	`gorm:"string(255)" sql:"DEFAULT:null" json:"failure_reason" tag_name:"failure_reason"`
	RequestedAt			int64	`gorm:"integer(11);not null" sql:"DEFAULT:null" json:"requested_at" tag_name:"requested_at"`
	FinishedAt			int64	`gorm:"integer(11)" sql:"DEFAULT:0" json:"finished_at" tag_name:"finished_at"`
	FailedAt			int64	`gorm:"integer(11)" sql:"DEFAULT:0" json:"failed_at" tag_name:"failed_at"`
}

type FileDetails struct {
	Id 					string	`gorm:"string(255);not null" json:"id" tag_name:"id"`
	DownloadRequestId	string 	`gorm:"string(255);not null" json:"download_request_id" tag_name:"download_request_id"`
	FileName			string	`gorm:"string(255);not null" json:"file_name" tag_name:"file_name"`
	FileSize			int64	`gorm:"integer(11)" sql:"DEFAULT:null" json:"file_size" tag_name:"file_size"`
	FileType			string	`gorm:"string(255);not null" json:"file_type" tag_name:"file_type"`
	Status				string	`gorm:"string(255);not null" json:"status" tag_name:"status"`
	CreatedAt			int64	`gorm:"integer(11)" sql:"DEFAULT:null" json:"created_at" tag_name:"created_at"`
}