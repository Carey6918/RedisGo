package redis

import (
	"bytes"
	"encoding/json"
	"github.com/Carey6918/RedisGo"
)

type ValueType string

const (
	String ValueType = "string"
	Hash   ValueType = "hash"
	List   ValueType = "list"
	Set    ValueType = "set"
	ZSet   ValueType = "zset"
)

type Package struct {
	Key      string
	Type     ValueType
	Contents []Content // Set: n; not Set: 1
}

type Content struct {
	Item  string
	Score float64
	//Bytes []byte
	Text string
	JSON string
}

// 列出所有的key和数据类型
func ReadAll() []Package {
	cli := RedisGo.GetClient()
	keys, _, err := cli.Scan(0, "*", 50000).Result()
	if err != nil {
		return nil
	}
	packages := make([]Package, 0, len(keys))
	for _, key := range keys {
		vType, err := cli.Type(key).Result()
		if err != nil {
			continue
		}
		pack := Package{
			Key:      key,
			Type:     ValueType(vType),
			Contents: nil,
		}
		pack.Contents, err = GetService(pack.Type).Read(key)
		packages = append(packages, pack)
	}
	return packages
}

func convertJSON(str string) string {
	var out bytes.Buffer
	json.Indent(&out, []byte(str), "", "\t")
	return string(out.Bytes())
}

// 删除某个元素
func DeleteElement(key string, content Content) {
	vType, _ := RedisGo.GetClient().Type(key).Result()
	GetService(ValueType(vType)).Delete(key, content)
}

// 删除某些key
func DeleteKeys(keys []string) {
	RedisGo.GetClient().Del(keys...)
}

func UpdateContent(key string, content Content) {
	vType, _ := RedisGo.GetClient().Type(key).Result()
	GetService(ValueType(vType)).Update(key, content)
}
