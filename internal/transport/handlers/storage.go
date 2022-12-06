package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"github.com/abdukhashimov/student_aggregator/pkg/slugify"
)

type FileUploadInfo struct {
	FileKey  string `json:"file_key"`
	FileURl  string `json:"file_url"`
	FileName string `json:"file_name"`
}

// @Summary Uploads files to blob storage
// @Description Creates new blob object in storage with file name
// @Security UsersAuth
// @Tags file-upload
// @Accept multipart/form-data
// @Param file_name formData string true "new object name in storage"
// @Param file formData file true "file"
// @Success 200 {object} FileUploadInfo
// @Failure 400
// @Failure 401
// @Failure 500
// @Produce json
// @Router /storage/upload [post]
func (s *Server) blobUpload() http.HandlerFunc {

	const (
		fileNameField = "file_name"
		fileBodyField = "file"

		MB = 1 << 20
	)

	var (
		contentType = "application/octet-stream"
	)

	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(int64(s.config.Project.FileUploadMaxMegabytes) * MB)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, M{
				"error": "can't parse request form",
			})
			return
		}

		fileName := r.FormValue(fileNameField)
		if len(fileName) == 0 {
			writeJSON(w, http.StatusBadRequest, M{
				"error": "file name not provided",
			})
			return
		}

		fileName = slugify.GenSlugWithID(fileName)

		body, fileHeader, err := r.FormFile(fileBodyField)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, M{
				"error": "can't parse file content",
			})
			return
		}

		if fileHeader.Size == 0 {
			writeJSON(w, http.StatusBadRequest, M{
				"error": "file size is 0",
			})
			return
		}

		if types, ok := fileHeader.Header["Content-Type"]; ok && len(types) > 0 {
			contentType = types[0]
		}

		key, err := s.storageService.PutFile(r.Context(), domain.PutFileOptions{
			ObjectName: fileName, Body: body, Size: fileHeader.Size, ContentType: contentType})
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, M{
				"error": "can't put file to storage",
			})
			return
		}

		writeJSON(w, http.StatusOK, FileUploadInfo{
			FileKey:  base64.RawStdEncoding.EncodeToString([]byte(key)),
			FileURl:  fmt.Sprintf("%s/%s/%s", s.config.Storage.URI, s.config.Storage.BucketName, fileName),
			FileName: fileName,
		})
	}
}
