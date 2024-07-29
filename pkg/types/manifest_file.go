package types

import (
	"archive/zip"
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
		return fmt.Errorf("failed to seek to the start of the file: %v", err)
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	if stat.IsDir() {
		return nil
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

	return nil
}

func (m *ManifestFile) ReadArchiveFileContent(file *zip.File) error {
	f, err := file.Open()
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if len(content) == 0 {
		// return fmt.Errorf("Content is empty")
		return nil
	}

	m.Content = content

	return nil
}
