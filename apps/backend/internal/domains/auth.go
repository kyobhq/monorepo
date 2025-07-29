package domains

import (
	"backend/internal/broker"
	"backend/internal/crypto"
	"backend/internal/database"
	"backend/internal/types"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	SignIn(ctx *gin.Context, user *types.SignInParams) (*string, *types.APIError)
	SignUp(ctx *gin.Context, user *types.SignUpParams) (*string, *types.APIError)
	Logout(ctx *gin.Context) *types.APIError
}

type authService struct {
	db     database.Service
	broker broker.Service
}

func NewAuthService(db database.Service, broker broker.Service) *authService {
	return &authService{
		db:     db,
		broker: broker,
	}
}

func (s *authService) SignIn(ctx *gin.Context, user *types.SignInParams) (*string, *types.APIError) {
	if user.Email == "admin" {
		return nil, &types.APIError{
			Status:  http.StatusForbidden,
			Code:    "ERR_ADMIN",
			Message: "no",
		}
	}

	dbUser, err := s.db.GetUser(ctx, user.Email)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusNotFound,
			Code:    "ERR_USER_NOT_FOUND",
			Cause:   err.Error(),
			Message: "User not found.",
		}
	}

	match, err := crypto.VerifyPassword(user.Password, dbUser.Password)
	if err != nil || !match {
		return nil, &types.APIError{
			Status:  http.StatusUnauthorized,
			Code:    "ERR_INVALID_CREDENTIALS",
			Cause:   err.Error(),
			Message: "Given credentials are invalid.",
		}
	}

	token, err := crypto.GenerateRandomBytes(64)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_TOKEN_GENERATION",
			Cause:   err.Error(),
			Message: "Failed to generate auth token.",
		}
	}

	b64Token := base64.RawStdEncoding.EncodeToString(token)
	err = s.broker.CacheUser(ctx, b64Token, dbUser)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_CACHING_USER",
			Cause:   err.Error(),
			Message: "Failed to cache the user in memdb.",
		}
	}

	return &b64Token, nil
}

func (s *authService) SignUp(ctx *gin.Context, user *types.SignUpParams) (*string, *types.APIError) {
	hashedPassword, err := crypto.HashPassword(user.Password)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_PASSWORD_HASHING_FAILED",
			Cause:   err.Error(),
			Message: "Failed to hash the password.",
		}
	}

	user.Password = hashedPassword

	dbUser, err := s.db.CreateUser(ctx, user)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_FAILED_USER_CREATION",
			Cause:   err.Error(),
			Message: "Failed to create the user account.",
		}
	}

	token, err := crypto.GenerateRandomBytes(64)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_TOKEN_GENERATION",
			Cause:   err.Error(),
			Message: "Failed to generate auth token.",
		}
	}

	b64Token := base64.RawStdEncoding.EncodeToString(token)
	err = s.broker.CacheUser(ctx, b64Token, dbUser)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_CACHING_USER",
			Cause:   err.Error(),
			Message: "Failed to cache the user in memdb.",
		}
	}

	return &b64Token, nil
}

func (s *authService) Logout(ctx *gin.Context) *types.APIError {
	token, err := ctx.Cookie("token")
	if err != nil {
		return &types.APIError{
			Status:  http.StatusNotFound,
			Code:    "ERR_MISSING_TOKEN",
			Cause:   err.Error(),
			Message: "No token found.",
		}
	}

	err = s.broker.RemoveCachedUser(ctx, token)
	if err != nil {
		return &types.APIError{
			Status:  http.StatusInternalServerError,
			Code:    "ERR_REMOVE_CACHED_USER",
			Cause:   err.Error(),
			Message: "Failed to disconnect user.",
		}
	}

	return nil
}
