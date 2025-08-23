package validation

import (
	"backend/internal/types"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validator *validator.Validate

func New() {
	Validator = validator.New()
	Validator.RegisterValidation("emoji_shortcode", validateEmojiShortcode)
}

func validateEmojiShortcode(fl validator.FieldLevel) bool {
	shortcode := fl.Field().String()

	// - Starts with [a-z]
	// - Ends with [a-z0-9]
	// - Middle can be [a-z0-9_] but no consecutive underscores
	// - Length 2-20
	pattern := `^[a-z]([a-z0-9]|_[a-z0-9])*[a-z0-9]$|^[a-z]$`

	return regexp.MustCompile(pattern).MatchString(shortcode)
}

func ParseAndValidate[T any](r *http.Request, body *T) *types.APIError {
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return types.NewAPIError(http.StatusBadRequest, "ERR_INVALID_BODY", "Invalid JSON body", err)
	}

	return Validate(body)
}

func formatValidationErrors(errs validator.ValidationErrors) error {
	var messages []string
	for _, err := range errs {
		messages = append(messages, fmt.Sprintf("field '%s' %s",
			strings.ToLower(err.Field()),
			getValidationMessage(err),
		))
	}
	return fmt.Errorf("validation failed: [%s]", strings.Join(messages, ", "))
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters", fe.Param())
	default:
		return fmt.Sprintf("failed validation '%s'", fe.Tag())
	}
}

type FileValidationConfig struct {
	MaxSize  int64
	MaxFiles int
}

var DefaultFileConfig = FileValidationConfig{
	MaxSize:  10 * 1024 * 1024, // 10MB
	MaxFiles: 10,
}

func ValidateFiles(fileHeaders []*multipart.FileHeader, config FileValidationConfig) *types.APIError {
	if len(fileHeaders) == 0 {
		return types.NewAPIError(http.StatusBadRequest, "ERR_NO_FILES", "No files provided", nil)
	}

	if len(fileHeaders) > config.MaxFiles {
		return types.NewAPIError(
			http.StatusBadRequest,
			"ERR_TOO_MANY_FILES",
			fmt.Sprintf("Too many files. Maximum allowed: %d", config.MaxFiles),
			nil,
		)
	}

	for i, fileHeader := range fileHeaders {
		if err := ValidateSingleFile(fileHeader, config); err != nil {
			return types.NewAPIError(
				err.Status,
				err.Code,
				fmt.Sprintf("File %d (%s): %s", i+1, fileHeader.Filename, err.Message),
				fmt.Errorf("file %d validation failed: %s", i+1, err.Cause),
			)
		}
	}

	return nil
}

func ValidateSingleFile(fileHeader *multipart.FileHeader, config FileValidationConfig) *types.APIError {
	if fileHeader.Size == 0 {
		return types.NewAPIError(
			http.StatusBadRequest,
			"ERR_EMPTY_FILE",
			"Empty file not allowed",
			nil,
		)
	}

	if fileHeader.Size > config.MaxSize {
		return types.NewAPIError(
			http.StatusBadRequest,
			"ERR_FILE_TOO_LARGE",
			fmt.Sprintf("File size %.2fMB exceeds maximum %.2fMB",
				float64(fileHeader.Size)/(1024*1024),
				float64(config.MaxSize)/(1024*1024)),
			nil,
		)
	}

	return nil
}

func Validate[T any](body *T) *types.APIError {
	err := Validator.Struct(body)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return types.NewAPIError(
				http.StatusBadRequest,
				"ERR_VALIDATION_FAILED",
				"Validation failed",
				formatValidationErrors(validationErrors),
			)
		}
		return types.NewAPIError(http.StatusBadRequest, "ERR_VALIDATION_FAILED", "Validation failed", err)
	}

	return nil
}

var sanitizeQueryRegex = regexp.MustCompile(`[%_ .!?|\\\/\x60~=+\-*&^$#@]+`)

func SanitizeQuery(query string) string {
	query = strings.TrimSpace(query)
	return sanitizeQueryRegex.ReplaceAllString(query, "")
}
