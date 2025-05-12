package main

import (
	"embed"

	"github.com/dxps/go_playground/go_static_file_serve/internal/handlers"
)

//go:embed assets
var assetsFS embed.FS

const PORT = 9001

func main() {

	h := handlers.New(PORT, assetsFS)
	h.Start()
}
