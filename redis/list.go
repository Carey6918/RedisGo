package redis

import "github.com/Carey6918/RedisGo"

type ListCRUDServiceImpl struct{}

func (*ListCRUDServiceImpl) Create(key string, content Content) error {
	return RedisGo.GetClient().LPush(key, content.Text).Err()
}

func (*ListCRUDServiceImpl) Read(key string) ([]Content, error) {
	values, err := RedisGo.GetClient().LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	contents := make([]Content, 0, len(values))
	for index, value := range values {
		contents = append(contents, Content{
			Item:  value,
			Score: float64(index),
			Text:  value,
			JSON:  convertJSON(value),
		})
	}
	return contents, nil
}

func (*ListCRUDServiceImpl) Update(key string, content Content) error {
	return RedisGo.GetClient().LSet(key, int64(content.Score), content.JSON).Err()
}

func (*ListCRUDServiceImpl) Delete(key string, content Content) error {
	return RedisGo.GetClient().LRem(key, 0, content.JSON).Err()
}
