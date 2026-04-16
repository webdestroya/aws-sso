package testutil

import (
	"os"
	"testing"
	"testing/fstest"
	"time"

	"github.com/webdestroya/aws-sso/internal/utils/fsutils"
)

func NewFakeFS(t *testing.T) *FakeFSUtil {
	fs := &FakeFSUtil{
		Files: make(fstest.MapFS),
	}

	oldval := fsutils.Global
	t.Cleanup(func() {
		fsutils.Global = oldval
	})
	fsutils.Global = fs

	return fs
}

type FakeFSUtil struct {
	Files fstest.MapFS

	AtomicWriteFileFn func(string, []byte, os.FileMode) error
	WriteFileFn       func(string, []byte, os.FileMode) error
	ReadFileFn        func(string) ([]byte, error)
	StatFn            func(string) (os.FileInfo, error)
}

func (fs *FakeFSUtil) AtomicWriteFile(filename string, data []byte, fileMode os.FileMode) error {
	if fs.AtomicWriteFileFn != nil {
		return fs.AtomicWriteFileFn(filename, data, fileMode)
	}

	return fs.WriteFile(filename, data, fileMode)
}

func (fs *FakeFSUtil) WriteFile(filename string, data []byte, fileMode os.FileMode) error {
	if fs.WriteFileFn != nil {
		return fs.WriteFileFn(filename, data, fileMode)
	}

	fs.Files[filename] = &fstest.MapFile{
		Data:    data,
		Mode:    fileMode,
		ModTime: time.Now(),
		Sys:     nil,
	}

	return nil
}

func (fs *FakeFSUtil) ReadFile(filename string) ([]byte, error) {
	if fs.ReadFileFn != nil {
		return fs.ReadFileFn(filename)
	}

	return fs.Files.ReadFile(filename)
}

func (fs *FakeFSUtil) Stat(filename string) (os.FileInfo, error) {
	if fs.StatFn != nil {
		return fs.StatFn(filename)
	}

	return fs.Files.Stat(filename)
}
