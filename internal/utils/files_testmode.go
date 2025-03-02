//go:build testmode

package utils

import "os"

func AtomicWriteFile(filename string, data []byte, fileMode os.FileMode) error {
	return nil
}

func WriteFile(filename string, data []byte, fileMode os.FileMode) error {
	return nil
}
