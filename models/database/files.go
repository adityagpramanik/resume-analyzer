package databasemodels

import (
	"time"
)

type File struct {
	FileId string
	Url string
	Bucket string
	FileName string
	CreatedAt time.Time
	UpdatedAt time.Time
}