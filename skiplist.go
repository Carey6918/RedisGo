package RedisGo

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type node struct {
	member interface{}
	score  int64
	next   []*node
}

type list struct {
	head   *node
	level  int
	mu     *sync.RWMutex
}

func (l *list) String() string {
	res := ""
	current := l.head
	for i := l.level - 1; i >= 0; i-- {
		res += fmt.Sprintf("level %v: ", i)
		for current != nil {
			res += fmt.Sprintf("%v ", current.score)
			current = current.next[i]
		}
		res += "\n"
		current = l.head
	}
	return res
}

const MaxLevel = 32
const Probability = 0.25

func (l *list) randLevel() (level int) {
	rand.Seed(time.Now().UnixNano())
	for level = 1; rand.Float32() < Probability && level < MaxLevel; level++ {
	}
	return
}

func newSkipList() *list {
	return &list{head: &node{
		member: nil,
		score:  0,
		next:   make([]*node, MaxLevel),
	}, level: 1,
		mu: &sync.RWMutex{}}
}

func (l *list) Insert(score int64, member interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.head
	update := make([]*node, MaxLevel)
	for i := l.level - 1; i >= 0; i-- {
		if current.next[i] == nil || current.next[i].score > score {
			update[i] = current
		} else {
			for current.next[i] != nil && current.next[i].score < score {
				current = current.next[i]
			}
			update[i] = current
		}
	}

	level := l.randLevel()
	if level > l.level {
		for i := l.level; i < level; i++ {
			update[i] = l.head
		}
		l.level = level
	}
	node := &node{
		member: member,
		score:  score,
		next:   make([]*node, level),
	}
	for i := 0; i < level; i++ {
		node.next[i] = update[i].next[i]
		update[i].next[i] = node
	}
}

func (l *list) Delete(score int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.head
	for i := l.level - 1; i >= 0; i-- {
		for current.next[i] != nil {
			if current.next[i].score == score {
				tmp := current.next[i]
				current.next[i] = tmp.next[i]
				tmp.next[i] = nil
			} else if current.next[i].score > score {
				break
			} else {
				current = current.next[i]
			}
		}
	}
}

func (l *list) Get(score int64) interface{} {
	return l.get(score).member
}

func (l *list) get(score int64) *node {
	l.mu.RLock()
	defer l.mu.RUnlock()
	current := l.head
	for i := l.level - 1; i >= 0; i-- {
		for current.next[i] != nil {
			if current.next[i].score == score {
				return current.next[i]
			} else if current.next[i].score > score {
				break
			} else {
				current = current.next[i]
			}
		}
	}
	return nil
}

func (l *list) Range(min, max int64) map[int64]interface{} {
	result := make(map[int64]interface{}, 0)
	current := l.head
	for i := l.level - 1; i >= 0; i-- {
		for current.next[i] != nil {
			if current.next[i].score > max {
				break
			} else if current.next[i].score < min {
				current = current.next[i]
			} else {
				result[current.next[i].score] = current.next[i].member
				current = current.next[i]
			}
		}
		current = l.head
	}
	return result
}
