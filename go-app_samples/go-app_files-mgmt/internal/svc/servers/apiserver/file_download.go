package apiserver

import (
	"encoding/json"
	"errors"
	"go-app_files-mgmt/internal/common"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *ApiServer) handleFilesList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rsp := common.NewUploadedFilesList()
	for _, file := range s.uploadedFiles {
		rsp.Files = append(rsp.Files, common.UploadedFile{
			FileName:    file.FileName,
			ContentType: file.ContentType,
			Size:        file.Size,
		})
	}
	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		slog.Error("Failed to encode files list.", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ApiServer) handleFileDownload(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		log.Fatalln(errors.New("Invalid method. It must be GET."))
		return
	}

	fileName := chi.URLParam(r, "filename")
	if fileName == "." || fileName == ".." {
		http.Error(w, "Invalid file name.", http.StatusBadRequest)
		return
	}

	if file, ok := s.uploadedFiles[fileName]; ok {
		w.Header().Set("Content-Type", file.ContentType)
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		if _, err := w.Write(file.Content); err != nil {
			slog.Error("Failed to serve file.", "filename", fileName, "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "File not found.", http.StatusNotFound)
	}

}
