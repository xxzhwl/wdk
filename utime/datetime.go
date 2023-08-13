// Package utime 包描述
// Author: wanlizhan
// Date: 2023/6/10
package utime

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

func init() {
	os.Setenv("TZ", "Asia/Shanghai")
}

// Date 获取当前日期
func Date() string {
	return time.Now().Format(time.DateOnly)
}

// DateTime 获取当前日期-时间
func DateTime() string {
	return time.Now().Format(time.DateTime)
}

// TimeStamp 获取时间戳
func TimeStamp() int64 {
	return time.Now().Unix()
}

// IsWeekend 是不是周六日
func IsWeekend(t time.Time) bool {
	return time.Saturday == t.Weekday() || time.Sunday == t.Weekday()
}

// GetZeroHourTimestamp 获取0点时间戳
func GetZeroHourTimestamp() int64 {
	ts := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", ts)
	return t.UTC().Unix() - 8*3600
}

// GetNightTimestamp 获取0点时间戳 59
func GetNightTimestamp() int64 {
	return GetZeroHourTimestamp() + 86400 - 1
}

// DateTimeFormat 根据指定format格式化时间
func DateTimeFormat(format string) string {
	return time.Now().Format(format)
}

func StrToTime(value string) (time.Time, error) {
	if len(value) == 0 {
		return time.Time{}, nil
	}
	matchList := map[string]string{
		"2006-01-02 15:04:05":      `^\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}$`,
		"2006-01-02 15:04":         `^\d{4}-\d{2}-\d{2}\s\d{2}:\d{2}$`,
		"2006-01-02 15":            `^\d{4}-\d{2}-\d{2}\s\d{2}$`,
		"2006-01-02":               `^\d{4}-\d{2}-\d{2}$`,
		"2006-01":                  `^\d{4}-\d{2}$`,
		"2006-01-02T15:04:05.000Z": `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d{3}Z$`, //同RFC3339标准
	}
	for format, itm := range matchList {
		re := regexp.MustCompile(itm)
		if re.MatchString(value) {
			return _strToTime(value, format)
		}
	}
	return time.Time{}, fmt.Errorf("date '%s' can not convert to time.Time object", value)
}

// IsUTCTime 判断一个String是不是标准UTC时间
func IsUTCTime(value string) bool {
	itm := `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d{3}Z$`
	re := regexp.MustCompile(itm)
	return re.MatchString(value)
}

// 将string的时间格式类型(如2017-11-12)，根据format格式(如Y-m-d)，转成一个time.Time的可操作对像
func _strToTime(value string, format string) (time.Time, error) {
	if IsUTCTime(value) {
		return time.Parse(format, value)
	}
	//使用本地时间解析
	return time.ParseInLocation(format, value, time.Local)
}

// AddHour 时间增加时
func AddHour(t time.Time, num int) time.Time {
	return t.Add(time.Duration(num) * time.Hour)
}

// AddMin 时间增加分
func AddMin(t time.Time, num int) time.Time {
	return t.Add(time.Duration(num) * time.Minute)
}

// AddSec 时间增加秒
func AddSec(t time.Time, num int) time.Time {
	return t.Add(time.Duration(num) * time.Second)
}

// AddDay 时间增加日
func AddDay(t time.Time, num int) time.Time {
	return t.AddDate(0, 0, num)
}

// AddMouth 时间增加月
func AddMouth(t time.Time, num int) time.Time {
	return t.AddDate(0, num, 0)
}

// AddYear 时间增加年
func AddYear(t time.Time, num int) time.Time {
	return t.AddDate(num, 0, 0)
}
