package RedisGo

import (
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	ExpirationTime = 5 * time.Second
	RetryTimes     = 3
	RetryInterval  = 500 * time.Millisecond
)

type goID int
type key2Times map[string]int
type goLocks map[goID]key2Times

var locks goLocks

// 如果要实现可重入锁的话，需要获得当前线程的id
// 但是在go中，并没有提供获取goroutine id的操作，开发者解释如下：
// Please don't use goroutine local storage.
// It's highly discouraged.
// In fact, IIRC, we used to expose Goid,
// but it is hidden since we don't want people to do this
// 因此，我们在这里通过runtime.Stack的堆栈信息来获取GoID
func GoroutineID() goID {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		return -1
	}
	return goID(id)
}

func Lock(key string) bool {
	id := GoroutineID()
	if _, ok := locks[id]; !ok {
		locks[id] = make(map[string]int)
	}
	key2Times := locks[id]
	if _, ok := key2Times[key]; ok {
		key2Times[key]++
		time.AfterFunc(ExpirationTime, func() {
			unlock(key, id)
		})
		return true
	}
	ok := lock(key)
	if !ok {
		return false
	}
	key2Times[key] = 1
	time.AfterFunc(ExpirationTime, func() {
		unlock(key, id)
	})
	return true
}

func Unlock(key string) bool {
	return unlock(key, GoroutineID())
}

func lock(key string) bool {
	failTimes := 0
	for i := 0; i < RetryTimes; i++ {
		ok, err := GetClient().SetNX(key, true, ExpirationTime).Result()
		if err != nil {
			failTimes++
			continue
		}
		if !ok {
			time.Sleep(RetryInterval) // 如果没有从redis获取到锁，等待一段时间重试
			continue
		}
		return true
	}
	if failTimes > RetryTimes/2 { // 超过半数是redis连接错误，认为是redis的问题，可以获取到锁
		return true
	}
	return false // 否则，是真的没有从redis中获取到锁
}

func del(key string) bool {
	for i := 0; i < RetryTimes; i++ {
		err := GetClient().Del(key).Err()
		if err != nil {
			continue
		}
		return true
	}
	return false
}

// 封装这么一层是为了区分goroutine自身释放锁和超时释放内存锁、
// 因为time的超时是另起一个goroutine来跑的，无法释放调用方的锁
func unlock(key string, id goID) bool {
	if _, ok := locks[id]; !ok {
		return false
	}
	key2Times := locks[id]
	if _, ok := key2Times[key]; !ok {
		return false
	}
	key2Times[key]--
	if key2Times[key] <= 0 {
		del(key)
		delete(key2Times, key)
	}
	if len(key2Times) == 0 {
		delete(locks, key2Times)
	}
	return true
}
