package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thiagohmm/DesafioTecnicoRateLimiter/internal/limiter"
)

// RateLimiterMiddleware verifica se a requisição excede o limite baseado no IP ou token
func RateLimiterMiddleware(lim *limiter.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()              // Pega o IP do cliente
		token := c.GetHeader("API_KEY") // Pega o token do cabeçalho da requisição

		// Escolhe o limitador apropriado (por IP ou token) e obtém o limite
		limiterInstance, key, limit := lim.GetLimiter(ip, token)

		// Verifica se a requisição é permitida de acordo com o limite
		allowed, err := limiterInstance.AllowRequest(key, limit, lim.GetBlockTime())
		if err != nil {
			// Caso ocorra um erro interno, retorna uma resposta 500
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}

		if !allowed {
			// Se o limite for excedido, retorna uma resposta 429 (Too Many Requests)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "you have reached the maximum number of requests"})
			return
		}

		// Se tudo estiver ok, a requisição segue para o próximo middleware ou handler
		c.Next()
	}
}
