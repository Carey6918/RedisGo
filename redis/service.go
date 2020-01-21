package redis

type CRUDService interface {
	Create(key string, content Content) error
	Read(key string) ([]Content, error)
	Update(key string, content Content) error
	Delete(key string,content Content) error
}

func GetService(t ValueType) CRUDService {
	switch t {
	case String:
		return new(StringCRUDServiceImpl)
	case Hash:
		return new(HashCRUDServiceImpl)
	case List:
		return new(ListCRUDServiceImpl)
	case Set:
		return new(SetCRUDServiceImpl)
	case ZSet:
		return new(ZSetCRUDServiceImpl)
	}
	return nil
}
