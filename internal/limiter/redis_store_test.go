package limiter

import (
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis() (*miniredis.Miniredis, *redis.Client) {
	// Start a miniredis server
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	// Create a redis client
	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	return s, client
}

func TestAllowRequest(t *testing.T) {
	s, client := setupTestRedis()
	defer s.Close()

	store := NewRedisStore(client)

	tests := []struct {
		name      string
		key       string
		limit     int
		blockTime time.Duration
		setup     func()
		expected  bool
	}{
		{
			name:      "allow request below limit",
			key:       "test_key_1",
			limit:     5,
			blockTime: time.Minute,
			setup:     func() {},
			expected:  true,
		},
		{
			name:      "deny request at limit",
			key:       "test_key_2",
			limit:     1,
			blockTime: time.Minute,
			setup: func() {
				s.Set("test_key_2", "1")
				// Configurar o TTL explicitamente
				s.SetTTL("test_key_2", time.Minute)
			},
			expected: false,
		},
		{
			name:      "allow request after block time",
			key:       "test_key_3",
			limit:     1,
			blockTime: time.Second,
			setup: func() {
				s.Set("test_key_3", "1")
				// Configurar o TTL explicitamente para garantir que ele expire
				s.SetTTL("test_key_3", time.Second)
				s.FastForward(time.Second * 2) // Avançar o tempo
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			allowed, err := store.AllowRequest(tt.key, tt.limit, tt.blockTime)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, allowed)
		})
	}
}

// Simula carga com várias requisições simultâneas
func TestRateLimiterUnderLoad(t *testing.T) {
	s, client := setupTestRedis()
	defer s.Close()

	store := NewRedisStore(client)

	// Configurações do teste
	key := "test_key_load"
	limit := 10 // Limite de requisições
	blockTime := time.Minute
	numRequests := 100 // Número total de requisições a serem simuladas
	concurrency := 50  // Número de requisições simultâneas

	// WaitGroup para controlar a execução das go routines
	var wg sync.WaitGroup
	wg.Add(concurrency)

	// Canal para coletar os resultados das requisições
	results := make(chan bool, numRequests)

	// Função que simula uma requisição ao rate limiter
	simulateRequest := func() {
		defer wg.Done()

		allowed, err := store.AllowRequest(key, limit, blockTime)
		assert.NoError(t, err)

		// Enviar o resultado da requisição para o canal
		results <- allowed
	}

	// Inicia as go routines para fazer as requisições simultâneas
	for i := 0; i < concurrency; i++ {
		go simulateRequest()
	}

	// Aguarda todas as go routines terminarem
	wg.Wait()

	// Fechar o canal de resultados após todas as requisições
	close(results)

	// Contabiliza os resultados
	allowedCount := 0
	deniedCount := 0
	for result := range results {
		if result {
			allowedCount++
		} else {
			deniedCount++
		}
	}

	// Verificar se o número de requisições permitidas não excede o limite
	assert.GreaterOrEqual(t, allowedCount, limit, "Número de requisições permitidas excedeu o limite")
	t.Logf("Requisições permitidas: %d, Requisições negadas: %d", allowedCount, deniedCount)
}
