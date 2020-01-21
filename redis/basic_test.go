package redis

import (
	"encoding/json"
	"fmt"
	"github.com/Carey6918/RedisGo"
	"strconv"
	"testing"
	"time"
)

type TestTask struct {
	ID   string `json:"id"`
	Test string `json:"test"`
}

func TestGetKeysAndTypes(t *testing.T) {
	RedisGo.Init()
	task := TestTask{
		ID:   "56789",
		Test: "fefefe",
	}
	body, _ := json.Marshal(task)
	RedisGo.GetClient().Set("test_key", body, time.Hour)
	PrintPackages(ReadAll())
}

func PrintPackages(ps []Package) {
	for _, p := range ps {
		fmt.Println(p.Key + ": " + string(p.Type))
		for _, c := range p.Contents {
			fmt.Println()
			fmt.Println("	" + "Item: " + c.Item)
			fmt.Println("	" + "Score: " + strconv.FormatFloat(c.Score, 'f', -1, 64))
			fmt.Println("	" + "Value: " + c.Text)
			fmt.Println("	" + "JSON: " + c.JSON)
		}
		fmt.Println()
		fmt.Println()
	}
}
