//go:build js && wasm

package pages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-app_files-mgmt/internal/common"
	"go-app_files-mgmt/internal/ui/comps"
	"go-app_files-mgmt/internal/ui/infra"
	"log/slog"
	"mime/multipart"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type FilesPage struct {
	app.Compo
	apiClient                  *infra.ApiClient
	DownloadableFilenames      []string `json:"files"`
	selectedFilenameToDownload string
}

func NewFilesPage(apiClient *infra.ApiClient) *FilesPage {
	return &FilesPage{
		apiClient:             apiClient,
		DownloadableFilenames: []string{},
	}
}

func (p *FilesPage) OnMount(ctx app.Context) {
	p.getFilesList()
}

func (p *FilesPage) Render() app.UI {

	return app.Div().Class("flex flex-col min-h-screen bg-gray-100").Body(
		&comps.Navbar{},
		app.Div().Class("flex flex-col min-h-screen justify-center items-center text-gray-800").
			Body(
				app.Div().
					Class("flex flex-col items-center bg-white p-6 rounded-lg drop-shadow-2xl min-w-[610px]").
					Body(
						app.H1().
							Class("text-3xl text-gray-400 mb-8").
							Text("File Upload"),
						app.Div().
							Text("Select a file to upload. After selecting one, it will be automatically read and uploaded."),
						app.Div().
							Text("Therefore, open the browser's Developer Tools' console and network to see the result."),
						app.Input().
							Class("border-0 mt-4 bg-slate-100 hover:bg-green-100").
							Type("file").
							Name("file-import-test.txt").
							Accept(".txt").
							OnInput(func(ctx app.Context, e app.Event) {
								p.handleTextFileUpload(e)
								p.getFilesList()
								ctx.Update()
							}),
					),
				app.Hr().Class("m-8"),
				app.Div().Class("flex flex-col items-center bg-white p-6 rounded-lg drop-shadow-2xl min-w-[610px]").
					Body(
						app.H1().Class("text-3xl text-gray-400 m-8").
							Text("File Download"),
						app.Range(p.DownloadableFilenames).Slice(func(i int) app.UI {
							return app.Div().
								Class("text-sm text-gray-600 py-1 px-4 hover:bg-green-100 rounded-lg transition duration-200 cursor-pointer").
								Text(fmt.Sprintf("File %d: %s", i, p.DownloadableFilenames[i])).
								OnClick(func(ctx app.Context, e app.Event) {
									p.selectedFilenameToDownload = p.DownloadableFilenames[i]
									p.handleDownload()
								})
						}),
					),
			),
	)
}

func (p *FilesPage) handleTextFileUpload(e app.Event) {

	// TODO: Currently, this does not handle multiple files.
	file := e.Get("target").Get("files").Index(0)

	// Read bytes from uploaded file.
	fileData, err := readFile(file)
	if err != nil {
		slog.Error("Failed to read uploaded file.", "error", err)
		app.Log(err)
		return
	}
	slog.Debug("Uploaded file",
		"name", file.Get("name").String(),
		"size", file.Get("size").Int(),
		"type", file.Type().String(),
		"data", fileData)

	// Reset file input for next upload.
	e.Get("target").Set("value", "")

	// Upload file to server.
	p.uploadFile(file.Get("name").String(), fileData)
}

func readFile(file app.Value) (data []byte, err error) {

	done := make(chan bool)
	// https://developer.mozilla.org/en-US/docs/Web/API/FileReader
	reader := app.Window().Get("FileReader").New()
	reader.Set("onloadend", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		done <- true
		return nil
	}))
	reader.Call("readAsArrayBuffer", file)
	<-done

	readerError := reader.Get("error")
	if !readerError.IsNull() {
		err = fmt.Errorf("file reader error : %s", readerError.Get("message").String())
	} else {
		uint8Array := app.Window().Get("Uint8Array").New(reader.Get("result"))
		data = make([]byte, uint8Array.Length())
		app.CopyBytesToGo(data, uint8Array)
	}

	return data, err
}

func (p *FilesPage) uploadFile(fileName string, fileData []byte) {

	buf := &bytes.Buffer{}
	mpw := multipart.NewWriter(buf)

	ffw1, err := mpw.CreateFormField("file")
	if err != nil {
		slog.Error("Failed to create form file.", "error", err)
		return
	}
	if _, err := ffw1.Write([]byte(fileName)); err != nil {
		slog.Error("Failed to write to form file.", "error", err)
		return
	}

	ffw2, err := mpw.CreateFormFile("filename", fileName)
	if err != nil {
		slog.Error("Failed to create form file.", "error", err)
		app.Log(err)
		return
	}
	if _, err := ffw2.Write(fileData); err != nil {
		slog.Error("Failed to write to form file.", "error", err)
		app.Log(err)
		return
	}
	if err := mpw.Close(); err != nil {
		slog.Error("Failed to close multipart writer.", "error", err)
		app.Log(err)
		return
	}
	// Close the multipart writer before creating the request.
	if err := mpw.Close(); err != nil {
		slog.Error("Failed to close multipart writer.", "error", err)
		app.Log(err)
		return
	}
	// Send the request.
	resp, err := p.apiClient.SendFile(common.FilesPath, mpw.FormDataContentType(), buf.Bytes())

	if err != nil {
		slog.Error("Failed to upload file.", "error", err)
		return
	}
	slog.Debug("File uploaded.", "response", string(resp))
}

func (p *FilesPage) getFilesList() {

	data, err := p.apiClient.Get(common.FilesPath)
	if err != nil {
		slog.Error("Failed to get files list.", "error", err)
	}
	if err := json.Unmarshal(data, &p); err != nil {
		slog.Error("Failed to unmarshal files list.", "error", err)
	}
}

func (p *FilesPage) handleDownload() {

	slog.Debug("Downloading ...", "filename", p.selectedFilenameToDownload)
	data, err := p.apiClient.Get(fmt.Sprintf("/files/%s", p.selectedFilenameToDownload))
	if err != nil {
		slog.Error("Failed to download file.", "error", err)
		return
	}
	uint8Array := app.Window().Get("Uint8Array").New(len(data))
	nb := app.CopyBytesToJS(uint8Array, data)
	if nb != len(data) {
		slog.Error("Unable to copy bytes to JS")
		return
	}
	blobConstructorParam := app.Window().Get("Array").New(uint8Array)
	blob := app.Window().Get("Blob").New(blobConstructorParam, map[string]interface{}{
		"type": "mime/type",
	})
	url := app.Window().Get("URL").JSValue().Call("createObjectURL", blob)
	a := app.Window().Get("document").Call("createElement", "a")
	a.Set("href", url)
	a.Set("download", p.selectedFilenameToDownload)
	a.Call("click")
	app.Window().Get("URL").JSValue().Call("revokeObjectURL", url)
}
