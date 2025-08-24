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
		return nil, types.NewAPIError(http.StatusForbidden, "ERR_ADMIN", "no", nil)
	}

	dbUser, err := s.db.GetUser(ctx, user.Email)
	if err != nil {
		return nil, types.NewAPIError(http.StatusNotFound, "ERR_USER_NOT_FOUND", "User not found.", err)
	}

	match, err := crypto.VerifyPassword(user.Password, dbUser.Password)
	if err != nil || !match {
		return nil, types.NewAPIError(http.StatusUnauthorized, "ERR_INVALID_CREDENTIALS", "Given credentials are invalid.", err)
	}

	token, err := crypto.GenerateRandomBytes(64)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_TOKEN_GENERATION", "Failed to generate auth token.", err)
	}

	b64Token := base64.RawStdEncoding.EncodeToString(token)
	err = s.broker.CacheUser(ctx, b64Token, dbUser)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_CACHING_USER", "Failed to cache the user in memdb.", err)
	}

	return &b64Token, nil
}

func (s *authService) SignUp(ctx *gin.Context, user *types.SignUpParams) (*string, *types.APIError) {
	hashedPassword, err := crypto.HashPassword(user.Password)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_PASSWORD_HASHING_FAILED", "Failed to hash the password.", err)
	}

	user.Password = hashedPassword

	dbUser, err := s.db.CreateUser(ctx, user)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_FAILED_USER_CREATION", "Failed to create the user account.", err)
	}

	token, err := crypto.GenerateRandomBytes(64)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_TOKEN_GENERATION", "Failed to generate auth token.", err)
	}

	b64Token := base64.RawStdEncoding.EncodeToString(token)
	err = s.broker.CacheUser(ctx, b64Token, dbUser)
	if err != nil {
		return nil, types.NewAPIError(http.StatusInternalServerError, "ERR_CACHING_USER", "Failed to cache the user in memdb.", err)
	}

	return &b64Token, nil
}

func (s *authService) Logout(ctx *gin.Context) *types.APIError {
	token, err := ctx.Cookie("token")
	if err != nil {
		return types.NewAPIError(http.StatusNotFound, "ERR_MISSING_TOKEN", "No token found.", err)
	}

	err = s.broker.RemoveCachedUser(ctx, token)
	if err != nil {
		return types.NewAPIError(http.StatusInternalServerError, "ERR_REMOVE_CACHED_USER", "Failed to disconnect user.", err)
	}

	return nil
}
