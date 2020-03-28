package utils

import "time"

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func BeginTimestampOfDayForSec(timeStamp int64) int64 {
	t := time.Unix(timeStamp, 0)
	return timeStamp - int64(t.Hour()*3600+t.Minute()*60+t.Second())
}
