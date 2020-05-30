package timerTask

import (
	"time"

	"../utils"
)

var (
	log = utils.InfoLogger
)

type TaskType int

const (
	TaskTypeNone TaskType = iota
	TaskTypeRepeat
	TaskTypeFireDate
)

type TaskSubType int

const (
	TaskSubTypeNone   TaskSubType = iota
	TaskSubTypeSlow               //5Min
	TaskSubTypeNormal             //1.5Min
	TaskSubTypeFast               //10S
	TaskSubTypeSlow1              //1h
	TaskSubTypeSlow2              //1d
)

type TaskHandleResultType int

const (
	TaskHandleResultTypeNone TaskHandleResultType = iota
	TaskHandleResultTypeRemove
	TaskHandleResultTypeDelay //延迟5Min
)

type task struct {
	name      string
	callBack  func(interface{}) TaskHandleResultType
	t         TaskType
	st        TaskSubType
	parameter interface{}
	fireDate  *time.Time
}

var (
	taskMap     = &map[string]task{}
	suspendTask = map[string]*time.Time{}
)

/**
 * 多次使用会进行比较，使用最久的时间
 */
func SuspendRegisterdTaskBefore(name string, t time.Time) {
	if len(name) <= 0 {
		return
	}
	if time.Now().UnixNano() > t.UnixNano() {
		return
	}
	preT := suspendTask[name]
	if preT != nil && preT.UnixNano() > t.UnixNano() {
		//Ignore shortter suspend
	} else {
		suspendTask[name] = &t
	}
}

func ResumeRegisterdTask(name string) {
	delete(suspendTask, name)
}

func RegisterOnTimeFire(name string, fireDate *time.Time, callBack func(interface{}) TaskHandleResultType, parameter interface{}) {
	RegisterTask(name, TaskTypeFireDate, TaskSubTypeNone, fireDate, callBack, parameter)
}

func RegisterRepeat(name string, callBack func(interface{}) TaskHandleResultType, parameter interface{}) {
	RegisterTask(name, TaskTypeRepeat, TaskSubTypeNormal, nil, callBack, parameter)
}

func RegisterSlowRepeat(name string, callBack func(interface{}) TaskHandleResultType, parameter interface{}) {
	RegisterTask(name, TaskTypeRepeat, TaskSubTypeSlow, nil, callBack, parameter)
}

func RegisterFastRepeat(name string, callBack func(interface{}) TaskHandleResultType, parameter interface{}) {
	RegisterTask(name, TaskTypeRepeat, TaskSubTypeFast, nil, callBack, parameter)
}

func RegisterSlow1Repeat(name string, callBack func(interface{}) TaskHandleResultType, parameter interface{}) {
	RegisterTask(name, TaskTypeRepeat, TaskSubTypeSlow1, nil, callBack, parameter)
}

func RegisterSlow2Repeat(name string, callBack func(interface{}) TaskHandleResultType, parameter interface{}) {
	RegisterTask(name, TaskTypeRepeat, TaskSubTypeSlow2, nil, callBack, parameter)
}

func RegisterTask(name string, t TaskType, st TaskSubType, fireDate *time.Time, callBack func(interface{}) TaskHandleResultType, parameter interface{}) {
	log.Logf("注册任务 -- %v -- %v-- %v\n", t, st, name)
	if len((*taskMap)[name].name) > 0 {
		panic("重复注册任务")
	}
	(*taskMap)[name] = task{name, callBack, t, st, parameter, fireDate}
}

func HasTaskWithName(name string) bool {
	for k, _ := range *taskMap {
		if k == name {
			return true
		}
	}
	return false
}

func UpdateFireDate(name string, fireDate *time.Time) {
	for k, v := range *taskMap {
		if k == name {
			v.fireDate = fireDate
			(*taskMap)[k] = v
			break
		}
	}
}
