package models

import "mime/multipart"

type UploadImageInput struct {
	ProductID int
	Header    *multipart.FileHeader
	IsPrimary bool
}

type UploadMultipleInput struct {
	ProductID int
	Files     []*multipart.FileHeader
}
