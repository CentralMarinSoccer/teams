package mocks

import (
	"io"
	"os"
	"github.com/centralmarinsoccer/teams/filesystem"
)

// FileSystem provides an interface to set and get data for mocked file system calls
type FileSystem struct {
	OpenCall struct {
		Receives struct {
			Filename string
		}
		Returns struct {
			Data  io.Reader
			Error error
		}
	}
	StatCall struct {
		Receives struct {
			Filename string
		}
		Returns struct {
			FileInfo os.FileInfo
			Error    error
		}
	}
	ReadFileCall struct {
		Receives struct {
			Filename string
		}
		Returns struct {
			Data  []byte
			Error error
		}
	}
	WriteCall struct {
		Receives struct {
			Data []byte
		}
		Returns struct {
			Count int
			Error error
		}
	}
}

// Open provides a mocked file system Open call
func (f *FileSystem) Open(name string) (filesystem.FileInterface, error) {
	f.OpenCall.Receives.Filename = name

	return f, nil
}

// Stat provides a mocked file system Stat call
func (f *FileSystem) Stat(name string) (os.FileInfo, error) {
	f.StatCall.Receives.Filename = name

	return f.StatCall.Returns.FileInfo, f.StatCall.Returns.Error
}

// ReadFile provides a mocked file system ReadFile call
func (f *FileSystem) ReadFile(filename string) ([]byte, error) {
	f.ReadFileCall.Receives.Filename = filename

	return f.ReadFileCall.Returns.Data, f.ReadFileCall.Returns.Error
}

// Close provides a mocked Close call
func (f *FileSystem) Close() error { return nil }

// Read provides a mocked Read call
func (f *FileSystem) Read(p []byte) (int, error) { return f.OpenCall.Returns.Data.Read(p) }

// Write provides a mocked Write call
func (f *FileSystem) Write(p []byte) (int, error) {
	f.WriteCall.Receives.Data = p
	return f.WriteCall.Returns.Count, f.WriteCall.Returns.Error
}
