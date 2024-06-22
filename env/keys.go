package env

import "os"

type Env struct {
	PostgreDB       string
	PostgreUser     string
	PostgrePassword string
	PostgreHost     string
	PostgrePort     string
	PostgreSSLMode  string
}

func GetEnvs() *Env {
	return &Env{
		PostgreDB:       os.Getenv("POSTGRES_DB"),
		PostgreUser:     os.Getenv("POSTGRES_USER"),
		PostgrePassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgreHost:     os.Getenv("POSTGRES_HOST"),
		PostgrePort:     os.Getenv("POSTGRES_PORT"),
		PostgreSSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}
