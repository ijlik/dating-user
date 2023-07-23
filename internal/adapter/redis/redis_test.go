package redis

import (
	"context"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockRedisClient struct {
	Mock mock.Mock
}
type RedisService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) bool
	Del(ctx context.Context, keys ...string) int
}

func NewRedisMockRepository(redisService RedisService) *MockRedisService {
	return &MockRedisService{
		redisService: redisService,
	}
}

type MockRedisService struct {
	redisService RedisService
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	return "test value", nil
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) bool {
	return true
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) int {
	return 1
}

func TestRedisRepository(t *testing.T) {
	// Create a mock Redis client
	mockClient := &MockRedisClient{Mock: mock.Mock{}}

	// Create the Redis repository using the mock client
	repo := NewRedisMockRepository(mockClient)

	// Test Get method
	ctx := context.Background()
	value, err := repo.redisService.Get(ctx, "test_key")
	assert.NoError(t, err)
	assert.Equal(t, "test value", value)

	// Test Set method
	result := repo.redisService.Set(ctx, "test_value", "test_key", 10)
	assert.True(t, result)

	// Test Del method
	intResult := repo.redisService.Del(ctx, "test_key")
	assert.Equal(t, intResult, 1)
}
