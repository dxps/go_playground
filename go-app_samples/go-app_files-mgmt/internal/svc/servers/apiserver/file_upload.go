package apiserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
)

func (s *ApiServer) handleFileUpload(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		log.Fatalln(errors.New("invalid method"))
		return
	}

	slog.Debug("Handing file upload ...")

	err := r.ParseMultipartForm(32 << 10) // 32 MB
	if err != nil {
		log.Fatal(err)
	}

	filename := r.FormValue("file")

	type uploadedFile struct {
		Size        int64  `json:"size"`
		ContentType string `json:"content_type"`
		Filename    string `json:"filename"`
		FileContent string `json:"file_content"`
	}

	var newFile uploadedFile
	for _, fheaders := range r.MultipartForm.File {
		for _, headers := range fheaders {
			file, err := headers.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			buff := make([]byte, 512)
			if _, err := file.Read(buff); err != nil {
				slog.Error("Failed to read file.", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := file.Seek(0, 0); err != nil {
				slog.Error("Failed to seek file.", "error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			contentType := http.DetectContentType(buff)
			newFile.ContentType = contentType

			var sizeBuff bytes.Buffer
			fileSize, err := sizeBuff.ReadFrom(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err := file.Seek(0, 0); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			newFile.Size = fileSize
			newFile.Filename = headers.Filename
			contentBuf := bytes.NewBuffer(nil)

			if _, err := io.Copy(contentBuf, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			newFile.FileContent = contentBuf.String()
		}
	}

	data := make(map[string]interface{})
	data["form_field_file"] = filename
	data["status"] = 200
	data["file_stats"] = newFile

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if err = json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
