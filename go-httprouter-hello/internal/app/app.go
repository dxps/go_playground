package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/dxps/go_playground/go-httprouter-hello/internal/api"
)

type App struct {
	cfg *LocalConfig
	api *api.API
	wg  *sync.WaitGroup
}

func NewApp() (*App, error) {

	a := App{}
	if err := a.initConfig(); err != nil {
		return nil, err
	}
	return &a, nil
}

func (a *App) initConfig() error {

	var err error
	if a.cfg, err = NewConfig(); err != nil {
		return fmt.Errorf("Local config init failed: %s", err)
	}
	return nil
}

func (a *App) Start(wg *sync.WaitGroup) error {

	a.wg = wg
	var err error
	a.api, err = api.NewAPI(a.cfg.API_host, a.cfg.API_port)
	if err != nil {
		return fmt.Errorf("Failed to start due to API error: %s", err)
	}

	a.startAPIJob()

	return nil
}

func (a *App) Stop(stopCtx context.Context) {

	if err := a.api.Shutdown(stopCtx); err != nil {
		log.Println("API shutdown issue: %w", err)
	}
	a.wg.Done()
}

func (a *App) startAPIJob() {

	go func() {
		if err := a.api.Serve(); err != http.ErrServerClosed {
			// TODO: Error handling needs improvement in this case.
			log.Fatalf("Error on HTTP ListenAndServe: %s", err)
		}
	}()
}
