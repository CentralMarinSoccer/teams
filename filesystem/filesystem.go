package filesystem

import (
	"io"
	"io/ioutil"
	"os"
)

// LocalDiskInterface creates an interface for the filesystem so we can mock it for testing
type LocalDiskInterface interface {
	Open(name string) (FileInterface, error)
	Stat(name string) (os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

// FileInterface creates an interface for a File so we can mock it for testing
type FileInterface interface {
	io.Closer
	io.Reader
	io.Writer
}

// OSFS implements fileSystem using the local disk.
type OSFS struct{}

// Open provides the default File System open
func (OSFS) Open(name string) (FileInterface, error) { return os.Open(name) }

// Stat provices the default File System stat
func (OSFS) Stat(name string) (os.FileInfo, error) { return os.Stat(name) }

func (OSFS) ReadFile(filename string) ([]byte, error) { return ioutil.ReadFile(filename) }
