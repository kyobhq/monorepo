package files

import (
	"backend/internal/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/davidbyttow/govips/v2/vips"
	"github.com/nrednav/cuid2"
)

type Service interface {
	UploadFile(key string, mimeType string, fileData io.Reader, fileName string) error
	ProcessAndUploadFiles(files []*multipart.FileHeader, maxSize int64) ([]byte, *types.APIError)
	ProcessAndUploadImage(imageToUpload *multipart.FileHeader, crop types.Crop, maxSize int64) (*string, *types.APIError)
	ProcessAndUploadAvatar(userID, imageType string, avatarToUpload *multipart.FileHeader, crop types.Crop, maxSize int64) (*string, *types.APIError)
	DeleteFile(key string) error
}

type service struct {
	s3Client *s3.Client
	bucket   string
	cdnURL   string
}

type File struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Filename string `json:"file_name"`
	Filesize string `json:"file_size"`
	Type     string `json:"type"`
}

func New() Service {
	region := os.Getenv("AWS_REGION")
	keyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	creds := credentials.NewStaticCredentialsProvider(keyID, secretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(creds))
	if err != nil {
		panic(err)
	}

	return &service{
		s3Client: s3.NewFromConfig(cfg),
		bucket:   "nyo-files",
		cdnURL:   os.Getenv("CDN_URL"),
	}
}

func (s *service) ProcessAndUploadFiles(filesToUpload []*multipart.FileHeader, maxSize int64) ([]byte, *types.APIError) {
	var files []File

	for _, fileHeader := range filesToUpload {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, &types.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "ERR_OPEN_FILE",
				Cause:   err.Error(),
				Message: "Failed to open file.",
			}
		}
		defer file.Close()

		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, &types.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "ERR_READ_FILE",
				Cause:   err.Error(),
				Message: "Failed to read file.",
			}
		}

		mimeType := http.DetectContentType(buffer[:n])

		if seeker, ok := file.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}

		randomID := cuid2.Generate()
		var key string
		var fileData io.Reader = file

		if strings.Contains(mimeType, "image") {
			key = fmt.Sprintf("attachment-%s.webp", randomID)

			webpData, err := convertToWebp(file)
			if err != nil {
				return nil, &types.APIError{
					Status:  http.StatusInternalServerError,
					Code:    "ERR_CONVERT_IMAGE_TO_WEBP",
					Cause:   err.Error(),
					Message: "Failed to convert image to webp.",
				}
			}
			fileData = bytes.NewReader(webpData)
		} else {
			extension := getSecureExtension(mimeType)
			key = fmt.Sprintf("attachment-%s.%s", randomID, extension)
		}

		if err := s.UploadFile(key, mimeType, fileData, fileHeader.Filename); err != nil {
			return nil, &types.APIError{
				Status:  http.StatusInternalServerError,
				Code:    "ERR_UPLOAD_FILE",
				Cause:   err.Error(),
				Message: "Failed to upload file.",
			}
		}
		defer file.Close()

		fileURL := fmt.Sprintf("%s/%s", s.cdnURL, key)
		fileSize := bytesToHuman(fileHeader.Size)

		attachment := File{
			ID:       randomID,
			URL:      fileURL,
			Filename: sanitizeFilename(fileHeader.Filename),
			Filesize: fileSize,
			Type:     mimeType,
		}

		files = append(files, attachment)
	}

	res, err := json.Marshal(files)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_MARSHAL_FILES",
			Cause:   err.Error(),
			Message: "Failed to marshal files.",
		}
	}

	return res, nil
}

func (s *service) ProcessAndUploadImage(imageToUpload *multipart.FileHeader, crop types.Crop, maxSize int64) (*string, *types.APIError) {
	file, err := imageToUpload.Open()
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_OPEN_FILE",
			Cause:   err.Error(),
			Message: "Failed to open file.",
		}
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_READ_FILE",
			Cause:   err.Error(),
			Message: "Failed to read file.",
		}
	}

	mimeType := http.DetectContentType(buffer[:n])
	if !strings.Contains(mimeType, "image") {
		return nil, &types.APIError{
			Status:  http.StatusBadRequest,
			Code:    "ERR_INVALID_MIME_TYPE",
			Message: "Invalid mime type.",
		}
	}

	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	randomID := cuid2.Generate()
	var key string
	var fileData io.Reader = file

	key = fmt.Sprintf("%s.webp", randomID)

	imageData, err := cropImage(file, crop)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_CROP_IMAGE",
			Cause:   err.Error(),
			Message: "Failed to crop the image.",
		}
	}

	fileData = bytes.NewReader(imageData)
	if err := s.UploadFile(key, mimeType, fileData, imageToUpload.Filename); err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPLOAD_FILE",
			Cause:   err.Error(),
			Message: "Failed to upload file.",
		}
	}
	defer file.Close()

	fileURL := fmt.Sprintf("%s/%s", s.cdnURL, key)

	return &fileURL, nil
}

func (s *service) ProcessAndUploadAvatar(userID, imageType string, avatarToUpload *multipart.FileHeader, crop types.Crop, maxSize int64) (*string, *types.APIError) {
	file, err := avatarToUpload.Open()
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_OPEN_FILE",
			Cause:   err.Error(),
			Message: "Failed to open file.",
		}
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_READ_FILE",
			Cause:   err.Error(),
			Message: "Failed to read file.",
		}
	}

	mimeType := http.DetectContentType(buffer[:n])
	if !strings.Contains(mimeType, "image") {
		return nil, &types.APIError{
			Status:  http.StatusBadRequest,
			Code:    "ERR_INVALID_MIME_TYPE",
			Message: "Invalid mime type.",
		}
	}

	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	var key string
	var fileData io.Reader = file

	key = fmt.Sprintf("%s-%s-%s.webp", userID, imageType, cuid2.Generate())

	imageData, err := cropImage(file, crop)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_CROP_IMAGE",
			Cause:   err.Error(),
			Message: "Failed to crop the image.",
		}
	}

	fileData = bytes.NewReader(imageData)
	if err := s.UploadFile(key, mimeType, fileData, avatarToUpload.Filename); err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_UPLOAD_FILE",
			Cause:   err.Error(),
			Message: "Failed to upload file.",
		}
	}
	defer file.Close()

	fileURL := fmt.Sprintf("%s/%s", s.cdnURL, key)

	return &fileURL, nil
}

func (s *service) UploadFile(key string, mimeType string, fileData io.Reader, fileName string) error {
	input := &s3.PutObjectInput{
		Key:    &key,
		Bucket: aws.String("nyo-files"),
		Body:   fileData,
	}

	if !strings.Contains(mimeType, "image") && !strings.Contains(mimeType, "video") {
		input.ContentDisposition = aws.String(fmt.Sprintf(`attachment; filename="%s"`,
			strings.ReplaceAll(fileName, `"`, `\"`)))
	}

	_, err := s.s3Client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (s *service) DeleteFile(key string) error {
	_, err := s.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Key:    &key,
		Bucket: aws.String("nyo-files"),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

func getSecureExtension(mimeType string) string {
	extensions := map[string]string{
		"application/pdf": "pdf",
		"text/plain":      "txt",
		"image/jpeg":      "jpg",
		"image/png":       "png",
	}

	if ext, ok := extensions[mimeType]; ok {
		return ext
	}

	if parts := strings.Split(mimeType, "/"); len(parts) == 2 {
		return parts[1]
	}

	return "bin"
}

func sanitizeFilename(filename string) string {
	filename = filepath.Base(filename)
	filename = strings.ReplaceAll(filename, "..", "")

	if len(filename) > 255 {
		filename = filename[:255]
	}

	return filename
}

func convertToWebp(file multipart.File) ([]byte, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	intSet := vips.IntParameter{}
	intSet.Set(-1)

	params := vips.NewImportParams()
	params.NumPages = intSet

	image, err := vips.LoadImageFromBuffer(fileBytes, params)
	if err != nil {
		return nil, err
	}
	defer image.Close()

	webp := vips.NewWebpExportParams()
	webp.Lossless = false
	webp.NearLossless = false
	webp.Quality = 85
	webp.StripMetadata = true

	buf, _, err := image.ExportWebp(webp)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func bytesToHuman(bytes int64) string {
	const unit = 1024
	units := []string{"B", "KB", "MB", "GB"}

	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	exp := int(math.Log(float64(bytes)) / math.Log(unit))
	value := float64(bytes) / math.Pow(unit, float64(exp))

	if value >= 100 {
		return fmt.Sprintf("%.0f %s", value, units[exp])
	} else if value >= 10 {
		return fmt.Sprintf("%.1f %s", value, units[exp])
	} else {
		return fmt.Sprintf("%.2f %s", value, units[exp])
	}
}

func cropImage(file multipart.File, crop types.Crop) ([]byte, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	intSet := vips.IntParameter{}
	intSet.Set(-1)

	params := vips.NewImportParams()
	params.NumPages = intSet

	image, err := vips.LoadImageFromBuffer(fileBytes, params)
	if err != nil {
		return nil, err
	}
	defer image.Close()

	err = image.ExtractArea(crop.X, crop.Y, crop.Width, crop.Height)
	if err != nil {
		return nil, err
	}

	webp := vips.NewWebpExportParams()
	webp.Lossless = false
	webp.NearLossless = false
	webp.Quality = 85
	webp.StripMetadata = true

	buf, _, err := image.ExportWebp(webp)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
