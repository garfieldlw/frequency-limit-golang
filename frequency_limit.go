package frequency_limit_golang

import (
	"context"
	"fmt"
	"github.com/garfieldlw/frequency-limit-golang/pkg/cache/redis"
	"github.com/garfieldlw/frequency-limit-golang/pkg/utils"
	"github.com/spf13/cast"
	"time"
)

// 5 times in 24h by user
const (
	limit     = 5
	frequency = 24 * 60 * 60
)

func IncrAndCheck(ctx context.Context, userId int64) (bool, error) {
	now := time.Now().Second()

	li, err := Check(ctx, userId, now)
	if err != nil {
		return false, err
	}

	if !li {
		return li, nil
	}

	client, err := redis.GetRedis()
	if err != nil {
		return false, err
	}

	key := fmt.Sprintf("%d", userId)
	err = client.HSet(ctx, key, utils.UUID(), fmt.Sprintf("%d", now), time.Second*frequency)
	if err != nil {
		return false, nil
	}

	return li, nil
}

func Check(ctx context.Context, userId int64, now int) (bool, error) {
	if now == 0 {
		now = time.Now().Second()
	}

	client, err := redis.GetRedis()
	if err != nil {
		return false, err
	}

	key := fmt.Sprintf("%d", userId)
	ttl, err := client.TTL(ctx, key)
	if err != nil {
		return false, err
	}

	if ttl < time.Second && ttl != -1 {
		return true, nil
	}

	items, err := client.HGetAll(ctx, key)
	if err != nil {
		return false, err
	}

	count := 0
	var delKeys []string
	for k, v := range items {
		keySecond := cast.ToInt(v)

		if keySecond+frequency < now {
			//remove key
			delKeys = append(delKeys, k)
		} else {
			count += 1
		}
	}

	if delKeys != nil || len(delKeys) > 0 {
		err = client.HDel(ctx, key, delKeys...)
		if err != nil {
			return false, err
		}
	}

	if count >= limit {
		return false, nil
	} else {
		return true, nil
	}
}
