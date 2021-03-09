package utils

import "os"

func IsExsit(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}