package common

import "time"

func GetTimeNow() time.Time {
	return time.Now()
}

func GetTimeNowUTC() time.Time {
	return time.Now().UTC()
}

func GetTimeNowUnix(nano bool) int64 {
	if nano {
		return time.Now().UnixNano()
	}
	return time.Now().Unix()
}
