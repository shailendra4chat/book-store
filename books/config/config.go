package config

import (
	"os"
)

func Conf(key string) string {
	return os.Getenv(key)
}
