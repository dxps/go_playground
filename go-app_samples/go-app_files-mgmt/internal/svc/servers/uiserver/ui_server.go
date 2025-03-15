package uiserver

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/klauspost/compress/gzhttp"
	"github.com/klauspost/compress/gzip"
)

func InitAndStartWebUiServerSide(uiPort, apiPort int) (*http.Server, error) {

	initAppHomeRoutesServerSide()

	handler := newCustomAppHandler()
	gzipWrapper, err := gzhttp.NewWrapper(gzhttp.MinSize(1000), gzhttp.CompressionLevel(gzip.BestSpeed))
	if err != nil {
		slog.Error("Failed to create gzip wrapper for app handler.", "error", err)
		return nil, err
	}
	uiSrv := http.Server{
		Addr:    fmt.Sprintf(":%d", uiPort),
		Handler: gzipWrapper(handler),
	}

	go func() {
		if err := uiSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	return &uiSrv, nil
}
