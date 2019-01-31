package request

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type UploadFile struct {
	FieldName string
	FileName  string
}

type UploadFiles []*UploadFile

type UploadRequest struct {
	ContentType string
	RequestBody string
}

func (u UploadFiles) ToRequestBody() (UploadRequest, error) {
	var writer *multipart.Writer

	body := &bytes.Buffer{}
	writer = multipart.NewWriter(body)

	for _, file := range u {
		fh, err := os.Open(file.FileName)
		if err != nil {
			return UploadRequest{}, err
		}
		defer fh.Close()

		part, err := writer.CreateFormFile(file.FieldName, filepath.Base(file.FileName))
		if err != nil {
			return UploadRequest{}, err
		}
		_, err = io.Copy(part, fh)

		err = writer.Close()
		if err != nil {
			return UploadRequest{}, err
		}
	}

	return UploadRequest{
		ContentType: writer.FormDataContentType(),
		RequestBody: body.String(),
	}, nil
}
