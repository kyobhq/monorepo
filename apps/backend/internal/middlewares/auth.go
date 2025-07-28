package middlewares

import (
	"backend/internal/broker"
	"backend/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(broker broker.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			error := types.NewAPIError(http.StatusUnauthorized, "ERR_INVALID_TOKEN", "The auth token is invalid.", err)
			error.Respond(c)
			return
		}

		user, err := broker.GetCachedUser(c, token)
		if err != nil {
			error := types.NewAPIError(http.StatusUnauthorized, "ERR_MISSING_CACHED_USER", "The auth token is invalid.", err)
			error.Respond(c)
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
