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
	ProcessAndUploadFiles(files []*multipart.FileHeader) ([]byte, *types.APIError)
	ProcessAndUploadEmojis(files []*multipart.FileHeader) ([]string, *types.APIError)
	ProcessAndUploadAvatar(entityID, imageType string, avatarToUpload *multipart.FileHeader, crop types.Crop) (*string, *types.APIError)
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

func (s *service) ProcessAndUploadFiles(filesToUpload []*multipart.FileHeader) ([]byte, *types.APIError) {
	var files []File

	for _, fileHeader := range filesToUpload {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_OPEN_FILE", "Failed to open file.", err)
		}
		defer file.Close()

		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_READ_FILE", "Failed to read file.", err)
		}

		mimeType := http.DetectContentType(buffer[:n])

		if seeker, ok := file.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}

		randomID := cuid2.Generate()
		var key string
		var fileData io.Reader = file

		if strings.Contains(mimeType, "image") {
			processedImg, err := processImageVersions(file, nil, mimeType)
			if err != nil {
				return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_PROCESS_IMAGE", "Failed to process image.", err)
			}

			key = fmt.Sprintf("attachment-%s.webp", randomID)
			fileData = bytes.NewReader(processedImg.StaticData)
			if err := s.UploadFile(key, mimeType, fileData, fileHeader.Filename); err != nil {
				return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_UPLOAD_FILE", "Failed to upload file.", err)
			}

			if processedImg.IsGIF && processedImg.AnimatedData != nil {
				key = fmt.Sprintf("attachment-%s-animated.webp", randomID)
				fileData = bytes.NewReader(processedImg.AnimatedData)
				if err := s.UploadFile(key, mimeType, fileData, fileHeader.Filename); err != nil {
					return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_UPLOAD_FILE", "Failed to upload animated file.", err)
				}
			}
		} else {
			extension := getSecureExtension(mimeType)
			key = fmt.Sprintf("attachment-%s.%s", randomID, extension)

			if err := s.UploadFile(key, mimeType, fileData, fileHeader.Filename); err != nil {
				return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_UPLOAD_FILE", "Failed to upload file.", err)
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
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_MARSHAL_FILES", "Failed to marshal files.", err)
	}

	return res, nil
}

func (s *service) ProcessAndUploadEmojis(emojisToUpload []*multipart.FileHeader) ([]string, *types.APIError) {
	var emojis []string

	for _, fileHeader := range emojisToUpload {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_OPEN_FILE", "Failed to open file.", err)
		}
		defer file.Close()

		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_READ_FILE", "Failed to read file.", err)
		}

		mimeType := http.DetectContentType(buffer[:n])
		if !strings.Contains(mimeType, "image") {
			return nil, types.NewAPIError(http.StatusBadRequest, "ERR_INVALID_MIME_TYPE", "Invalid mime type.", nil)
		}

		if seeker, ok := file.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}

		var key string
		var fileData io.Reader = file

		key = fmt.Sprintf("emoji-%s.webp", cuid2.Generate())

		emojiData, err := processEmoji(file)
		if err != nil {
			return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_PROCESS_EMOJI", "Failed to process the emoji.", err)
		}

		fileData = bytes.NewReader(emojiData)
		if err := s.UploadFile(key, mimeType, fileData, fileHeader.Filename); err != nil {
			return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_UPLOAD_EMOJI", "Failed to upload emoji.", err)
		}
		defer file.Close()

		fileURL := fmt.Sprintf("%s/%s", s.cdnURL, key)

		emojis = append(emojis, fileURL)
	}

	return emojis, nil
}

func (s *service) ProcessAndUploadAvatar(entityID, imageType string, avatarToUpload *multipart.FileHeader, crop types.Crop) (*string, *types.APIError) {
	file, err := avatarToUpload.Open()
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_OPEN_FILE", "Failed to open file.", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_READ_FILE", "Failed to read file.", err)
	}

	mimeType := http.DetectContentType(buffer[:n])
	if !strings.Contains(mimeType, "image") {
		return nil, types.NewAPIError(http.StatusBadRequest, "ERR_INVALID_MIME_TYPE", "Invalid mime type.", nil)
	}

	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	staticID := cuid2.Generate()
	var key string
	var fileData io.Reader

	processedImg, err := processImageVersions(file, &crop, mimeType)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_PROCESS_IMAGE", "Failed to process avatar image.", err)
	}

	key = fmt.Sprintf("%s-%s-%s.webp", entityID, imageType, staticID)
	fileData = bytes.NewReader(processedImg.StaticData)
	if err := s.UploadFile(key, mimeType, fileData, avatarToUpload.Filename); err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_UPLOAD_FILE", "Failed to upload avatar.", err)
	}

	if processedImg.IsGIF && processedImg.AnimatedData != nil {
		key = fmt.Sprintf("%s-%s-%s-animated.webp", entityID, imageType, staticID)
		fileData = bytes.NewReader(processedImg.AnimatedData)
		if err := s.UploadFile(key, mimeType, fileData, avatarToUpload.Filename); err != nil {
			return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_UPLOAD_FILE", "Failed to upload animated avatar.", err)
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

type ProcessedImage struct {
	StaticData   []byte
	AnimatedData []byte // nil if not a GIF
	IsGIF        bool
}

func processImageVersions(file multipart.File, crop *types.Crop, mimeType string) (*ProcessedImage, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	isGIF := mimeType == "image/gif"
	result := &ProcessedImage{IsGIF: isGIF}

	staticData, err := convertToWebpFromBytes(fileBytes, crop, false)
	if err != nil {
		return nil, fmt.Errorf("failed to crop static image: %w", err)
	}
	result.StaticData = staticData

	if isGIF {
		animatedData, err := convertToWebpFromBytes(fileBytes, crop, true)
		if err != nil {
			return nil, fmt.Errorf("failed to crop animated image: %w", err)
		}
		result.AnimatedData = animatedData
	}

	return result, nil
}

func convertToWebpFromBytes(fileBytes []byte, crop *types.Crop, isAnimated bool) ([]byte, error) {
	var intSet vips.IntParameter
	if isAnimated {
		intSet = vips.IntParameter{}
		intSet.Set(-1)
	} else {
		intSet = vips.IntParameter{}
		intSet.Set(1)
	}

	params := vips.NewImportParams()
	params.NumPages = intSet

	image, err := vips.LoadImageFromBuffer(fileBytes, params)
	if err != nil {
		return nil, err
	}
	defer image.Close()

	if crop != nil {
		err = image.ExtractArea(crop.X, crop.Y, crop.Width, crop.Height)
		if err != nil {
			return nil, err
		}
	}

	buf, _, err := image.ExportWebp(getWebpDefaultConfig())
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func processEmoji(file multipart.File) ([]byte, error) {
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

	buf, _, err := image.ExportWebp(getWebpDefaultConfig())
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

func getWebpDefaultConfig() *vips.WebpExportParams {
	return &vips.WebpExportParams{
		Lossless:      false,
		NearLossless:  false,
		Quality:       85,
		StripMetadata: true,
	}
}
