package RedisGo

import (
	"encoding/json"
	"time"
)

type funnel struct {
	capacity int64     // 最大容量，op
	rate     float64   // 最大频率，op/s
	quota    int64     // 剩余容量，op
	lastTime time.Time // 上次漏水时间（这里的漏水是惰性漏水，也可以理解为注水时间，只有注水时才会漏）
}

// 漏水
func (f *funnel) leak() {
	// 计算截止到上次漏水时间，漏斗漏了多少水
	leakedWater := int64(time.Now().Sub(f.lastTime).Seconds() * f.rate)
	// 如果漏的水不足1，则不记录本次漏水时间，也不做任何操作
	if leakedWater < 1 {
		return
	}
	// 更新漏水时间
	f.lastTime = time.Now()
	// 漏水，并保证剩余容量不会大于最大容量
	f.quota = f.quota + leakedWater
	if f.quota > f.capacity {
		f.quota = f.capacity
	}
}

// 注水
func (f *funnel) infuse(water int64) time.Duration {
	f.leak()
	// 如果漏斗是空的，代表可以直接执行，无需等待
	if f.quota == f.capacity {
		return 0
	}
	// 如果漏斗容量不足，则拒绝服务
	if f.quota < water {
		return -1
	}
	// 如果漏斗中容量足够，但是水没漏完，说明需要等待，计算需要等待的时间
	interval := float64(f.capacity-f.quota) / f.rate
	f.quota = f.quota - water
	return time.Duration(interval)
}

func getFunnel(key string) (*funnel, error) {
	v, err := GetClient().Get(key).Bytes()
	if err != nil {
		return nil, err
	}
	f := new(funnel)
	err = json.Unmarshal(v, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func setFunnel(key string, f *funnel) error {
	v, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return GetClient().Set(key, v, 60*time.Second).Err()
}

func isAllowedByFunnel(key string, rate float64, capacity int64) (time.Duration, error) {
	// 以下行为必须保证是原子操作，所以需要加锁
	if !Lock(key) {
		return -1, nil
	}
	defer Unlock(key)
	f, err := getFunnel(key)
	if err != nil {
		return -1, err
	}
	if f == nil {
		f = &funnel{capacity, rate, capacity, time.Now()}
	}
	// 检查的同时进行行为注入
	interval := f.infuse(1)
	setFunnel(key, f)
	return interval, nil
}
