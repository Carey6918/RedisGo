package RedisGo

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

func isAllowed(key string, period time.Duration, max int64) (bool, error) {
	// 移除窗口外的数据，减少存储消耗
	GetClient().ZRemRangeByScore(key, "0", string(int64(time.Now().Nanosecond())-period.Nanoseconds()))
	// 获取剩余的数据数量（即原本窗口内的数据）
	num, err := GetClient().ZCard(key).Result()
	// 设置ZSet过期时间，及时清空未继续操作的用户
	GetClient().Expire(key, period+1)
	// btw，由于这几个操作都是用的同一个key，所以使用pipeline可以提升redis存取效率，pipeline的使用我写在了：
	// isAllowedByPipeline(key, period)
	if err != nil {
		return false, err
	}
	if num < max {
		return true, nil
	}
	return false, nil
}

func Action(action string, userID string, period time.Duration, max int64) (bool, error) {
	key := fmt.Sprintf("limiter:%s:%s", action, userID)
	if ok, err := isAllowed(key, period, max); !ok || err != nil {
		return false, err
	}
	score := time.Now().Nanosecond()
	member := redis.Z{
		Score:  float64(score),
		Member: score, // 保证唯一性即可
	}
	// 这个方法将所有的行为都记录了下来，如果max值比较大，会非常消耗存储空间，所以要谨慎使用
	err := GetClient().ZAdd(key, member).Err()
	return true, err
}

func isAllowedByPipeline(key string, period time.Duration) {
	pipe := GetClient().Pipeline()
	pipe.ZRemRangeByScore(key, "0", string(int64(time.Now().Nanosecond())-period.Nanoseconds()))
	pipe.ZCard(key).Result()
	pipe.Expire(key, period+1)
	pipe.Exec()
	pipe.Close()
}
