package file

import (
	"io"
	"os"
)

// OpenWriter opens file for write and truncate old data.
// File opens on call resulting function.
func OpenWriter(filename string) WriterFunc {
	return func() (io.WriteCloser, error) {
		return os.Create(filename)
	}
}

// ExistsWriter just return writer.
func ExistsWriter(w io.WriteCloser) WriterFunc {
	return func() (io.WriteCloser, error) {
		return w, nil
	}
}

// OpenReaderFile opens file for read.
//
// If file not exists will return nil values.
// Storage supported this case (see Storage.Init).
//
// Also function returns io.ReadCloser for closing file in user code.
func OpenReaderFile(path string) (io.ReadCloser, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil

		}
		return nil, err
	}
	return file, nil
}
