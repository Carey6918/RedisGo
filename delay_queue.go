package RedisGo

import (
	"github.com/go-redis/redis"
	"time"
)

const DelayTime = 60

func PushTask(key, task string) error {
	score := time.Now().Second() + DelayTime
	member := redis.Z{
		Score:  float64(score),
		Member: task,
	}
	return GetClient().ZAdd(key, member).Err()
}

func LoopExecTask(key string) {
	for {
		opt := redis.ZRangeBy{
			Min:    "0",
			Max:    string(time.Now().Second()),
			Offset: 0,
			Count:  1, // 最多取1条
		}
		task, err := GetClient().ZRangeByScore(key, opt).Result()
		if err != nil {
			continue
		}
		if len(task) <= 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		err = GetClient().ZRem(key, task[0]).Err()
		if err != nil { // 删除成功以后才执行，确保一个task只被执行一次
			continue
		}
		ExecTask(task[0])
	}
}

func ExecTask(task string) {
	// do something
}
