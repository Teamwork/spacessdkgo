package util

import "github.com/joho/godotenv"

// LoadEnv loads variables from .env files into the process environment.
// Files are loaded in order; later files override earlier ones.
func LoadEnv(files ...string) error {
	return godotenv.Load(files...)
}

// MustLoadEnv loads variables from .env files and panics if loading fails.
func MustLoadEnv(files ...string) {
	if err := godotenv.Load(files...); err != nil {
		panic(err)
	}
}
