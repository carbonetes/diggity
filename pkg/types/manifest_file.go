package types

import (
	"fmt"
	"io"
	"os"
)

type ManifestFile struct {
	Content []byte
	Size    int64
	Path    string
	Layer   string
}

func (m *ManifestFile) ReadContent(file *os.File) error {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("Failed to seek to the start of the file: %v", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	m.Size = stat.Size()

	content, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	if len(content) == 0 {
		// return fmt.Errorf("Content is empty")
		return nil
	}

	m.Content = content

	err = file.Close()
	if err != nil {
		return err
	}

	err = os.Remove(file.Name())
	if err != nil {
		return err
	}

	return nil
}