package common

type UploadedFile struct {
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Filename    string `json:"filename"`
	FileContent []byte `json:"file_content"`
}

type UploadedFilesList struct {
	Files []string `json:"files"`
}

func NewUploadedFilesList() *UploadedFilesList {
	return &UploadedFilesList{
		Files: []string{},
	}
}
