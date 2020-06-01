package throttleTask

import (
	"time"

	"../utils"
)

func ThrottleDurationTask(key string, d time.Duration, t func() error) error {
	return utils.SynKeyAction(timeThrottle_handle_key, func() {
		nextExcThrottleTime := timeThrottle[key]
		if time.Now().After(nextExcThrottleTime) {
			err := t()
			if nil == err {
				timeThrottle[key] = time.Now()
			}
		}
	})
}
