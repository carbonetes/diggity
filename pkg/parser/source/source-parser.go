package source

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/carbonetes/diggity/pkg/model"
	"github.com/carbonetes/diggity/pkg/parser/bom"
)

var (
	// SourceInfo directory image information
	SourceInfo *model.SourceInfo
)

// ParseSourceProperties appends directory metadata to parser.Result
func ParseSourceProperties() {
	if *bom.Arguments.Dir == "" {
		bom.WG.Done()
		return
	}

	pathHash, err := getPathHash(*bom.Arguments.Dir)
	if err != nil {
		err = errors.New("source-parser: " + err.Error())
		bom.Errors = append(bom.Errors, &err)
	}

	SourceInfo = &model.SourceInfo{
		ID:   pathHash,
		Path: *bom.Arguments.Dir,
	}

	defer bom.WG.Done()
}

// Get the hash of the directory
func getPathHash(path string) (string, error) {
	// Create a new hash.Hash object to calculate the SHA-256 hash
	hash := sha256.New()

	// Walk the directory and add each file's contents to the hash
	err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(hash, file); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	// Get the final hash value as a byte slice
	hashBytes := hash.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashString := fmt.Sprintf("%x", hashBytes)

	return hashString, nil
}
