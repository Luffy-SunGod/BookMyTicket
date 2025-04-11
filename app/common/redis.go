package common

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func CreateRedisEntry(showID int, totalCapacity int) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer rdb.Close()

	// Context for the Redis operations.
	ctx := context.Background()

	// Update the value in Redis.
	err := rdb.Set(ctx, strconv.Itoa(showID), totalCapacity, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting value in Redis: %v", err)
	}

	return nil
}
