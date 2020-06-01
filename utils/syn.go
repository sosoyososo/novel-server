package utils

import (
	"errors"
	"time"
)

/**
锁需要有超时机制：如果我们对锁没有做超时机制，则多锁必然会有死锁情况。
锁超时时间不能太长：超时时间不能过长，不然会在死锁发生时候，急速加剧死锁发生概率。
锁定时间不能太短：如果太短，接近任务执行时间，则会导致过多的任务丢弃。
100ms的选择：我们后台绝大部分任务的执行会在50ms之内完成，涉及到加锁的任务时间会更短，这种情况下，我理解100ms是一个比较合适的时间。
*/

var (
	synPoolFetchChan = make(chan int, 1)
	synPool          = map[string]chan int{}
	synPoolCtx       = map[string]map[string]interface{}{}
)

func init() {
	synPoolFetchChan <- 1
}

func SyncMultiKeysAction(keys []string, ac func()) error {
	t := time.NewTimer(time.Microsecond * 100)
	defer func() {
		t.Stop()
	}()

	for _, key := range keys {
		c := tradeSynPoolGet(key)
		select {
		case <-c:
			defer func() {
				c <- 1
			}()
		case <-t.C: //超时返回//不要继续等待
			return errors.New("锁超时")
		}
	}
	ac()
	return nil
}

func SyncSetContextValue(syncKey, objKey string, value interface{}) {
	SynKeyAction(syncKey, func() {
		synPoolCtx[syncKey][objKey] = value
	})
}

func SyncGetContextValue(syncKey, objKey string) interface{} {
	var ret interface{}
	SynKeyAction(syncKey, func() {
		if nil != synPoolCtx[syncKey] {
			ret = synPoolCtx[syncKey][objKey]
		}
	})
	return ret
}

func SynKeyAction(key string, ac func()) error {
	return SyncMultiKeysAction([]string{key}, ac)
}

// 同步操作
func tradeSynPoolGet(uid string) chan int {
	<-synPoolFetchChan
	defer func() {
		synPoolFetchChan <- 1
	}()
	c := synPool[uid]
	if nil != c {
		return c
	}
	c = make(chan int, 1)
	c <- 1
	synPool[uid] = c
	return c
}
