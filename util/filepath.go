package util

import (
	"os"
	"path/filepath"
)

func GetExcellist(path string) ([]string, error) {
	var fileLists []string
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() && filepath.Ext(path) == ".xlsx" {
			fileLists = append(fileLists, path)
		}
		return nil
	})

	return fileLists, err
}
