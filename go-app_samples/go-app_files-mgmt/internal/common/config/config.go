package config

type Config struct {
	Servers struct {
		FrontendPort int `env:"FRONTEND_PORT, default=9001"`
		BackendPort  int `env:"BACKEND_PORT, default=9002"`
	}
}
