package fsutils

import "os"

type FSUtil interface {
	AtomicWriteFile(string, []byte, os.FileMode) error
	WriteFile(string, []byte, os.FileMode) error
	ReadFile(string) ([]byte, error)
	Stat(string) (os.FileInfo, error)
}

var Global FSUtil

func AtomicWriteFile(v string, data []byte, perm os.FileMode) error {
	return Global.AtomicWriteFile(v, data, perm)
}

func WriteFile(v string, data []byte, perm os.FileMode) error {
	return Global.WriteFile(v, data, perm)
}

func ReadFile(v string) ([]byte, error) {
	return Global.ReadFile(v)
}

func Stat(v string) (os.FileInfo, error) {
	return Global.Stat(v)
}
