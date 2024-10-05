package utils

import "os"

func IsDevEnv() bool {
	env := os.Getenv("ENV")
	return env == "" || env == "local" || env == "dev" || env == "development"
}
