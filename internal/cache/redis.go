package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"strings"
	"time"
)

var (
	RedisClient *redis.Client
	TTL         = 12 * time.Hour
	ErrNilCache = errors.New("nil cache content")
)

// CachedContent represents the structure of cached data
type CachedContent struct {
	Content     string    `json:"content"`
	LastUpdated time.Time `json:"last_updated"`
}

// InitRedis initializes the Redis connection
// In your InitRedis function:
func InitRedis() error {
	redisConfig := getRedisConfig()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.host, redisConfig.port),
		Password: redisConfig.password,
		DB:       0,
	})

	// Test connection first
	if err := testConnection(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	// After successful connection, set memory configs
	ctx := context.Background()

	// Add Memory Policy: Set memory policy to LRU
	if err := RedisClient.ConfigSet(ctx, "maxmemory-policy", "allkeys-lru").Err(); err != nil {
		return fmt.Errorf("setting maxmemory-policy: %w", err)
	}

	// Verify the setting was applied
	policy, err := RedisClient.ConfigGet(ctx, "maxmemory-policy").Result()
	if err != nil {
		return fmt.Errorf("getting maxmemory-policy: %w", err)
	}
	fmt.Printf("Memory Policy set to: %v\n", policy[1])

	fmt.Printf("Successfully connected to Redis at %s:%s", redisConfig.host, redisConfig.port)

	// Start cleanup routine in background
	go startExpiryCleanup()

	return nil
}

type redisConfig struct {
	host     string
	port     string
	password string
}

func getRedisConfig() redisConfig {
	config := redisConfig{
		host:     os.Getenv("REDIS_HOST"),
		port:     os.Getenv("REDIS_PORT"),
		password: os.Getenv("REDIS_PASSWORD"),
	}

	// Set defaults if not provided
	if config.host == "" {
		config.host = "localhost"
	}
	if config.port == "" {
		config.port = "6379"
	}

	return config
}

func testConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}

// startExpiryCleanup periodically scans for expired keys
func startExpiryCleanup() {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := cleanupExpiredKeys(); err != nil {
			fmt.Printf("Error during cleanup: %v\n", err)
		}
	}
}

func cleanupExpiredKeys() error {
	ctx := context.Background()
	var cursor uint64

	for {
		keys, nextCursor, err := RedisClient.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			return fmt.Errorf("scanning keys: %w", err)
		}

		for _, key := range keys {
			if err := handleExpiredKey(ctx, key); err != nil {
				fmt.Printf("Error handling key %s: %v\n", key, err)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

func handleExpiredKey(ctx context.Context, key string) error {
	ttl, err := RedisClient.TTL(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("checking TTL: %w", err)
	}

	if ttl < 0 {
		if err := RedisClient.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("deleting expired key: %w", err)
		}
		fmt.Printf("Deleted expired key: %s\n", key)
	}
	return nil
}

// GetCachedContent retrieves content from Redis
func GetCachedContent(ctx context.Context, key string) (*CachedContent, error) {
	data, err := RedisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNilCache
	}
	if err != nil {
		return nil, fmt.Errorf("reading from Redis: %w", err)
	}

	var content CachedContent
	if err := json.Unmarshal([]byte(data), &content); err != nil {
		return nil, fmt.Errorf("unmarshaling cached content: %w", err)
	}

	return &content, nil
}

// SetCachedContent stores content in Redis
func SetCachedContent(ctx context.Context, key string, content *CachedContent) error {
	if content == nil {
		return errors.New("nil content provided")
	}

	data, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("marshaling content: %w", err)
	}

	if err := RedisClient.Set(ctx, key, data, TTL).Err(); err != nil {
		return fmt.Errorf("writing to Redis: %w", err)
	}

	return nil
}

// GetCacheStats returns basic statistics about the Redis cache
func GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	memoryStats := make(map[string]string)

	// Get Redis info
	info := RedisClient.Info(ctx, "memory").Val()

	// Parse the memory info into a structured format
	for _, line := range strings.Split(info, "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			memoryStats[parts[0]] = strings.TrimSpace(parts[1])
		}
	}

	// Get total keys
	size, err := RedisClient.DBSize(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("getting DB size: %w", err)
	}

	stats["total_keys"] = size
	stats["memory_info"] = info         // Keep original for backwards compatibility
	stats["memory_stats"] = memoryStats // Add structured format

	return stats, nil
}
