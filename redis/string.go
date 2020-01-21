package redis

import (
	"github.com/Carey6918/RedisGo"
)

type StringCRUDServiceImpl struct{}

func (*StringCRUDServiceImpl) Create(key string, content Content) error {
	return RedisGo.GetClient().Set(key, content.Text, 0).Err()
}

func (*StringCRUDServiceImpl) Read(key string) ([]Content, error) {
	value, err := RedisGo.GetClient().Get(key).Result()
	if err != nil {
		return nil, err
	}
	return []Content{{
		Text: value,
		JSON: convertJSON(value),
	}}, nil
}

func (*StringCRUDServiceImpl) Update(key string, content Content) error {
	return RedisGo.GetClient().Set(key, content.JSON, 0).Err()
}

func (*StringCRUDServiceImpl) Delete(key string, content Content) error {
	return nil
}
