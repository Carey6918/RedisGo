package RedisGo

import "time"

type bucket struct {
	capacity int64     // 桶容量（最大令牌数）
	rate     float64   // 令牌产生速率，个/s
	quota    int64     // 剩余令牌数
	lastTime time.Time // 上次放入令牌的时间（由于惰性放入，也即上次拿走令牌的时间）
}

func (b *bucket) put() {
	putTokens := int64(time.Now().Sub(b.lastTime).Seconds() * b.rate)
	if putTokens < 1 {
		return
	}
	b.lastTime = time.Now()
	b.quota = b.quota + putTokens
	if b.quota > b.capacity {
		b.quota = b.capacity
	}
}

func (b *bucket) take(tokens int64) bool {
	b.put()
	if b.quota < tokens {
		return false
	}
	b.quota = b.quota - tokens
	return true
}

func isAllowedByTokenBucket(key string, rate float64, capacity int64) (bool, error) {
	// 其余操作与funnel类似，不再赘述
	b := &bucket{capacity, rate, capacity, time.Now()}
	return b.take(1), nil
}
