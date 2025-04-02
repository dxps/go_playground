package repos

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log/slog"
	"os"
)

const appWASMFilepath = "./web/app.wasm"

// GetGzCompressedAppWASM returns the contents of the `app.wasm` file, compressed using gzip.
func GetGzCompressedAppWASM() (compressedAppWASM []byte, appWASMSize int, err error) {

	bs, err := os.ReadFile(appWASMFilepath)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to read file '%s': %w", appWASMFilepath, err)
	}
	appWASMSize = len(bs)
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err = gz.Write(bs)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to write to gzip writer: %w", err)
	}
	if err := gz.Flush(); err != nil {
		return nil, 0, fmt.Errorf("Failed to flush the gzip writer: %w", err)
	}
	if err := gz.Close(); err != nil {
		return nil, 0, fmt.Errorf("Failed to close the gzip writer: %w", err)
	}
	compressedAppWASM = b.Bytes()
	slog.Debug(fmt.Sprintf("Compressed app.wasm size=%d KB", len(compressedAppWASM)/1024))
	return
}
