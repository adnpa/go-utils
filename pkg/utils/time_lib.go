package utils

import "time"

func NowSecond() int64 {
	return time.Now().Unix()
}

// NowSecondWithLocation
// ex. NowSecondWithLocation("Asia/Shanghai")
func NowSecondWithLocation(name string) int64 {
	loc, _ := time.LoadLocation(name)
	t := time.Now()
	return t.In(loc).Unix()
}

func NowMilliSecond() int64 {
	return time.Now().UnixMilli()
}

func DateTime() string {
	return time.Now().Format(time.DateTime)
}

func Date() string {
	return time.Now().Format(time.DateOnly)
}

func DateTime2Second(datetime string) int64 {
	t, _ := time.Parse(time.DateTime, datetime)
	return t.Unix()
}

func Date2Second(date string) int64 {
	t, _ := time.Parse(time.DateOnly, date)
	return t.Unix()
}

func Second2DateTime(second int64) string {
	t := time.Unix(second, 0)
	return t.Format(time.DateTime)
}

func Second2Date(second int64) string {
	t := time.Unix(second, 0)
	return t.Format(time.DateOnly)
}
