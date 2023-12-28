package main_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileConversion(t *testing.T) {
	// Create a temporary test directory
	testDir := t.TempDir()

	// Create a temporary test file with sample content
	content := "This is a test file content."
	testFilePath := filepath.Join(testDir, "test.txt")
	err := ioutil.WriteFile(testFilePath, []byte(content), 0644)
	assert.NoError(t, err)

	// Create a mock HTTP request with the test file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(testFilePath))
	assert.NoError(t, err)
	file, err := os.Open(testFilePath)
	assert.NoError(t, err)
	_, err = io.Copy(part, file)
	assert.NoError(t, err)
	writer.Close()

	// Create a mock HTTP POST request to the /convert endpoint
	req, err := http.NewRequest("POST", "/convert", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a mock HTTP response recorder
	recorder := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(recorder, req)

	// Check the HTTP response status code
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Check if the converted SRT file exists
	srtFilePath := filepath.Join(testDir, "outputs", "test.srt")
	assert.FileExists(t, srtFilePath)

	// Read the content of the converted SRT file
	srtContent, err := ioutil.ReadFile(srtFilePath)
	assert.NoError(t, err)

	// Define the expected content of the SRT file based on the test content
	expectedSRTContent := "1\n00:00:00,000 --> 00:00:04,000\nThis is a test file content.\n\n"
	assert.Equal(t, expectedSRTContent, string(srtContent))
}




