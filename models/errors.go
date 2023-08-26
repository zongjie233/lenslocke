package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrEmailTaken = errors.New("models: email address is already in use")

	ErrNotFound = errors.New("models: resource could not be found")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invaild file:%v", fe.Issue)
}

func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	// Create a byte array of 512 bytes
	testBytes := make([]byte, 512)
	// Read the file into the byte array
	_, err := r.Read(testBytes)
	// If there is an error, return an error
	if err != nil {
		return fmt.Errorf("models: error reading file: %v", err)
	}
	// Seek to the beginning of the file
	_, err = r.Seek(0, 0)
	// If there is an error, return an error
	if err != nil {
		return fmt.Errorf("models: error seeking file: %v", err)
	}
	// Get the content type of the file
	contentType := http.DetectContentType(testBytes)
	// Iterate through the allowed types
	for _, allowedType := range allowedTypes {
		// If the content type is allowed, return nil
		if contentType == allowedType {
			return nil
		}
	}
	// If the content type is not allowed, return an error
	return FileError{Issue: fmt.Sprintf("models: file content type is not allowed:%v", contentType)}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if !hasExtension(filename, allowedExtensions) {
		return FileError{Issue: fmt.Sprintf("models: file extension is not allowed:%v", filepath.Ext(filename))}
	}
	return nil
}
