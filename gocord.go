package gocordstorage

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	referer   = "jscord-storage-0.0.10"
)

type response struct {
	Status int
	Data   []byte
}

func Upload(filename string, fileURLOrPath string) (*response, error) {
	isURL, err := url.ParseRequestURI(fileURLOrPath)
	if err != nil {
		// Upload from file
		return uploadFromFile(filename, fileURLOrPath)
	}

	if isURL.Scheme != "" && isURL.Host != "" {
		// Upload from URL
		return uploadFromURL(filename, fileURLOrPath)
	}

	return nil, fmt.Errorf("invalid file URL or path")
}

func uploadFromURL(filename string, fileURL string) (*response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	_ = writer.WriteField("filename", filename)
	_ = writer.WriteField("url", fileURL)

	req, err := http.NewRequest("POST", "https://discord-storage.animemoe.us/upload-from-url/", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", referer)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &response{
		Status: resp.StatusCode,
		Data:   data,
	}, nil
}

func uploadFromFile(filename string, filePath string) (*response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	_ = writer.WriteField("filename", filename)

	req, err := http.NewRequest("POST", "https://discord-storage.animemoe.us/upload-from-file/", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", referer)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &response{
		Status: resp.StatusCode,
		Data:   data,
	}, nil
}
