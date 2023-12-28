package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type FileType int

const (
	Image FileType = iota
	Video
	Document
	Unknown
)

// CheckFileType checks the file extension and determines the type of file.
func CheckFileTypeExt(filename string) FileType {
	// Get the file extension
	extension := strings.ToLower(GetFileExtension(filename))

	// Check the file extension
	switch extension {
	case "jpg", "jpeg", "png", "gif":
		return Image
	case "mp4", "avi", "mkv":
		return Video
	case "pdf", "doc", "docx", "txt":
		return Document
	default:
		return Unknown
	}
}

// GetFileExtension returns the lowercase file extension of a filename.
func GetFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1]
	}
	return ""
}

// fileTypeToString converts FileType to a string representation.
func FileTypeToString(fileType FileType) string {
	switch fileType {
	case Image:
		return "Image"
	case Video:
		return "Video"
	case Document:
		return "Document"
	default:
		return "Unknown"
	}
}

func CheckFileType(fileName string) string {

	filetoCheck := CheckFileTypeExt(fileName)

	return FileTypeToString(filetoCheck)
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func GenerateUniqueFileName(ext string) string {

	timeStamp := time.Now().Format("20060102_150405")

	// Generate a random string
	randomString, err := generateRandomString(8)
	if err != nil {
		// Handle error (e.g., log it, return a default value, etc.)
		return "defaultFileName"
	}

	// Combine timestamp and random string to create a unique file name
	uniqueFileName := fmt.Sprintf("%s_%s", timeStamp, randomString)

	return uniqueFileName + "." + ext
}
