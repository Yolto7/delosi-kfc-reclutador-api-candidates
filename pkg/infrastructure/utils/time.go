package utils

import (
	"log"
	"time"
)

func nowInLocation(timeZone string) time.Time {
	location, err := time.LoadLocation(timeZone)
	if err != nil {
		log.Printf("invalid timezone '%s', falling back to UTC: %v", timeZone, err)
		location = time.UTC
	}
	return time.Now().In(location)
}

func NowDate(timeZone string) string {
	return nowInLocation(timeZone).Format("2006-01-02")
}

func NowTime(timeZone string) string {
	return nowInLocation(timeZone).Format("15:04:05")
}

func NowDateTime(timeZone string) string {
	return nowInLocation(timeZone).Format("2006-01-02 15:04:05")
}

func NowInTimezone(timeZone string) time.Time {
	return nowInLocation(timeZone)
}

func NowTimestamp(timeZone string) int64 {
	return nowInLocation(timeZone).Unix()
}

func Sleep(delay int) {
	t := time.NewTimer(time.Duration(delay) * time.Millisecond)
	defer t.Stop()
	<-t.C
}
