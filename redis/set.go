package redis

import "github.com/Carey6918/RedisGo"

type SetCRUDServiceImpl struct{}

func (*SetCRUDServiceImpl) Create(key string, content Content) error {
	return RedisGo.GetClient().SAdd(key, content.Text).Err()
}

func (*SetCRUDServiceImpl) Read(key string) ([]Content, error) {
	values, err := RedisGo.GetClient().SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	contents := make([]Content, 0, len(values))
	for _, value := range values {
		contents = append(contents, Content{
			Item:  value,
			Score: 0,
			Text:  value,
			JSON:  convertJSON(value),
		})
	}
	return contents, nil
}

func (*SetCRUDServiceImpl) Update(key string, content Content) error {
	return nil
}

func (*SetCRUDServiceImpl) Delete(key string, content Content) error {
	return RedisGo.GetClient().SRem(key, content.JSON).Err()
}
