package common

type UploadedFile struct {
	FileName    string `json:"filename"`
	ContentType string `json:"content_type,omitempty"`
	Content     []byte `json:"content,omitempty"`
	Size        int64  `json:"size"`
}

type UploadedFilesList struct {
	Files []UploadedFile `json:"files"`
}

func NewUploadedFilesList() *UploadedFilesList {
	return &UploadedFilesList{
		Files: []UploadedFile{},
	}
}
