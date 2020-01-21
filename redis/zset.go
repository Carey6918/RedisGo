package redis

import (
	"fmt"
	"github.com/Carey6918/RedisGo"
	"github.com/go-redis/redis"
)

type ZSetCRUDServiceImpl struct{}

func (*ZSetCRUDServiceImpl) Create(key string, content Content) error {
	return RedisGo.GetClient().ZAdd(key, redis.Z{
		Score:  content.Score,
		Member: content.Text,
	}).Err()
}

func (*ZSetCRUDServiceImpl) Read(key string) ([]Content, error) {
	values, err := RedisGo.GetClient().ZRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	contents := make([]Content, 0, len(values))
	for _, value := range values {
		score, _ := RedisGo.GetClient().ZScore(key, value).Result()
		contents = append(contents, Content{
			Item:  value,
			Score: score,
			Text:  value,
			JSON:  convertJSON(value),
		})
	}
	return contents, nil
}

func (*ZSetCRUDServiceImpl) Update(key string, content Content) error {
	oldScore, err := RedisGo.GetClient().ZScore(key, content.JSON).Result()
	if err != nil {
		return err
	}
	increment := content.Score - oldScore
	score, err := RedisGo.GetClient().ZIncrBy(key, increment, content.JSON).Result()
	if err != nil {
		return err
	}
	if score != content.Score {
		return fmt.Errorf("update score not equals setting")
	}
	return nil
}

func (*ZSetCRUDServiceImpl) Delete(key string, content Content) error {
	return RedisGo.GetClient().ZRem(key, content.JSON).Err()
}
