package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thiagohmm/DesafioTecnicoRateLimiter/internal/limiter"
)

func RateLimiterMiddleware(lim limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		allowed, err := lim.AllowRequest(ip, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "you have reached the maximum number of requests"})
			return
		}

		c.Next()
	}
}
