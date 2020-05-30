package timerTask

import (
	"time"

	"../utils"
)

const (
	slowDuration   = time.Second * 300
	normalDuration = time.Second * 90
	fastDuration   = time.Second * 10
	slow1Duration  = time.Hour
	slow2Duration  = time.Hour * 24
)

var (
	fastTimer   = time.NewTicker(fastDuration)
	normalTimer = time.NewTicker(normalDuration)
	slowTimer   = time.NewTicker(slowDuration)
	slow1Timer  = time.NewTicker(slow1Duration)
	slow2Timer  = time.NewTicker(slow2Duration)
)

func init() {
	utils.CallFuncInNewRecoveryRoutine(func() {
		for {
			<-fastTimer.C
			log.Logln("fast timer fired")
			handleRepeatTaskWithType(TaskSubTypeFast)
			handleOnTimeFireTask()
		}
	})

	utils.CallFuncInNewRecoveryRoutine(func() {
		for {
			<-normalTimer.C
			log.Logln("normal timer fired")
			handleRepeatTaskWithType(TaskSubTypeNormal)
		}
	})

	utils.CallFuncInNewRecoveryRoutine(func() {
		for {
			<-slowTimer.C
			log.Logln("slow timer fired")
			handleRepeatTaskWithType(TaskSubTypeSlow)
		}
	})
	utils.CallFuncInNewRecoveryRoutine(func() {
		for {
			<-slow1Timer.C
			log.Logln("slow1 timer fired")
			handleRepeatTaskWithType(TaskSubTypeSlow1)
		}
	})
	utils.CallFuncInNewRecoveryRoutine(func() {
		for {
			<-slow2Timer.C
			log.Logln("slow2 timer fired")
			handleRepeatTaskWithType(TaskSubTypeSlow2)
		}
	})
}

func handleRepeatTaskWithType(st TaskSubType) {
	for k, v := range *taskMap {
		if v.t == TaskTypeRepeat && v.st == st {
			log.Logf("开始执行任务 -- %v -- %v\n", v.t, v.name)
			ret := v.callBack(v.parameter)
			log.Logf("任务执行结束 %v\n", v.name)
			if ret == TaskHandleResultTypeRemove {
				defer delete(*taskMap, k)
			}
		}
	}
}

func handleOnTimeFireTask() {
	for k, v := range *taskMap {
		if v.t == TaskTypeFireDate && v.fireDate.After(time.Now()) {
			log.Logf("开始执行任务 -- %v -- %v\n", v.t, v.name)
			ret := v.callBack(v.parameter)
			log.Logf("任务执行结束 %v\n", v.name)
			if ret == TaskHandleResultTypeDelay {
				freshTime := v.fireDate.Add(time.Minute * 5)
				v.fireDate = &freshTime
				defer func() { (*taskMap)[k] = v }()
			} else {
				defer delete(*taskMap, k)
			}
		}
	}
}
