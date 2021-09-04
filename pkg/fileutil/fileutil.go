package fileutil

import (
	"os"
)

func Ls(dirname string) ([]string, []string, error) {
	var files, dirs []string
	if flag, err := IsDir(dirname); err != nil {
		return files, dirs, err
	} else {
		if !flag {
			return append(files, dirname), dirs, nil
		}
	}

	if dirEntries, err := os.ReadDir(dirname); err != nil {
		return files, dirs, err
	} else {
		for _, dirEntry := range dirEntries {
			name := dirname + "/" + dirEntry.Name()
			if dirEntry.IsDir() {
				dirs = append(dirs, name)
			} else {
				files = append(files, name)
			}
		}
	}
	return files, dirs, nil
}

func ReadFile(path string) ([]byte, error) {
	if bytes, err := os.ReadFile(path); err != nil {
		return bytes, err
	} else {
		return bytes, err
	}
}

func IsDir(path string) (bool, error) {
	if fileinfo, err := os.Stat(path); err != nil {
		return false, err
	} else {
		return fileinfo.IsDir(), nil
	}
}
