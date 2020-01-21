package redis

import (
	"fmt"
	"github.com/Carey6918/RedisGo"
)

type HashCRUDServiceImpl struct{}

func (*HashCRUDServiceImpl) Create(key string, content Content) error {
	exists, err := RedisGo.GetClient().HExists(key, content.Item).Result()
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("key already exists")
	}
	if err = RedisGo.GetClient().HSetNX(key, content.Item, content.Text).Err(); err != nil {
		return err
	}
	return nil
}

func (*HashCRUDServiceImpl) Read(key string) ([]Content, error) {
	values, err := RedisGo.GetClient().HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	contents := make([]Content, 0, len(values))
	for field, value := range values {
		contents = append(contents, Content{
			Item:  field,
			Score: 0,
			Text:  value,
			JSON:  convertJSON(value),
		})
	}
	return contents, nil
}

func (*HashCRUDServiceImpl) Update(key string, content Content) error {
	return RedisGo.GetClient().HSet(key, content.Item, content.JSON).Err()
}

func (*HashCRUDServiceImpl) Delete(key string, content Content) error {
	return RedisGo.GetClient().HDel(key, content.JSON).Err()
}
