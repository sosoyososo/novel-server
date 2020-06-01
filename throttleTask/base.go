package throttleTask

import "time"

var (
	timeThrottle            = map[string]time.Time{}
	timeThrottle_handle_key = "throttle_task_map_handler"
)
