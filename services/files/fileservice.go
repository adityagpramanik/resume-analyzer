package files

import (
	"time"

	databasemodels "resume-scanner.com/resume-scanner/models/database"
	commonservices "resume-scanner.com/resume-scanner/services/common"
)

func GetFileByFileId(file_id string) (databasemodels.File, error) {
	session := commonservices.GetSession()
	query := session.Query("select file_id, bucket, createdat, file_name, updatedat, url from resume_analyzer.files where file_id = ?", file_id)
	file := databasemodels.File{
		FileId:    "",
		Url:       "",
		Bucket:    "",
		FileName:  "",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	err := query.Scan(&file.FileId, &file.Bucket, &file.CreatedAt, &file.FileName, &file.UpdatedAt, &file.Url)
	if err != nil {
		return file, err
	}
	return file, nil
}
