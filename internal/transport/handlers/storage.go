package handlers

import (
	"encoding/base64"
	"net/http"
)

type FileUploadInfo struct {
	File_key string `json:"file_key"`
}

// @Summary Uploads files to blob storage
// @Description Creates new blob object in storage with file name
// @Tags file-upload
// @Accept multipart/form-data
// @Param file_name formData string true "new object name in storage"
// @Param file_object formData []byte true "new object base file"
// @Success 200 {object} FileUploadInfo
// @Failure 400
// @Failure 401
// @Failure 500
// @Produce json
// @Router /storage/upload [post]
func (s *Server) blobUpload() http.HandlerFunc {

	const (
		fileNameField = "file_name"
		fileBodyField = "file_object"

		MB = 1 << 20
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

		key, err := s.storageService.PutFile(r.Context(), fileName, body, fileHeader.Size)

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, M{
				"error": "can't put file to storage",
			})
			return
		}

		writeJSON(w, http.StatusOK, FileUploadInfo{
			File_key: base64.RawStdEncoding.EncodeToString([]byte(key)),
		})
	}
}
