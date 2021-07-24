package xfile

import "os"

func FileExists(filename string) bool {
	if info, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	} else {
		return !info.IsDir()
	}
}

func DirExists(dir string) bool {
	if info, err := os.Stat(dir); os.IsNotExist(err) {
		return false
	} else {
		return info.IsDir()
	}
}
