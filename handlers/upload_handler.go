package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadMedia(c *gin.Context) {
	err := c.Request.ParseMultipartForm(20 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds limit"})
		return
	}

	fileHeader := c.Request.MultipartForm.File["file"]
	if len(fileHeader) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only one file allowed to upload"})
		return
	}

	file, err := fileHeader[0].Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving file"})
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading file"})
		return
	}
	if n == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is empty"})
		return
	}

	fileType := http.DetectContentType(buffer)
	fileSize := fileHeader[0].Size

	var maxFileSize int64
	if strings.HasPrefix(fileType, "video/") {
		maxFileSize = 20 << 20
	} else if strings.HasPrefix(fileType, "image/") {
		maxFileSize = 5 << 20
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
		return
	}
	if fileSize > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds limit"})
		return
	}

	current := time.Now()
	folderName := fmt.Sprintf("%d-%02d-%02d", current.Year(), current.Month(), current.Day())

	fileHash := generateFileHash(fileHeader[0], current.String())
	fileName := fileHeader[0].Filename
	fileExtension := strings.ToLower(fileName[strings.LastIndex(fileName, ".")+1:])
	fileNameWithoutExt := fileName[:strings.LastIndex(fileName, ".")]
	newFileName := fmt.Sprintf("%s-%s.%s", fileHash, fileNameWithoutExt, fileExtension)

	uploadDir := fmt.Sprintf("uploads/%s", folderName)
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory: " + err.Error()})
			return
		}
	}

	fileNameWithPath := fmt.Sprintf("%s/%s", uploadDir, newFileName)
	out, err := os.Create(fileNameWithPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create the file for writing."})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, io.MultiReader(bytes.NewReader(buffer[:n]), file))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error copying file."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "filename": fileNameWithPath})
}

func generateFileHash(fileHeader *multipart.FileHeader, additionalData string) string {
	file, err := fileHeader.Open()
	if err != nil {
		return ""
	}
	defer file.Close()

	hasher := sha256.New()
	io.Copy(hasher, file)
	io.WriteString(hasher, additionalData)
	return hex.EncodeToString(hasher.Sum(nil))
}
