package datex

import "time"

var (
	TimeLocation, _ = time.LoadLocation("Asia/Shanghai")
	DateTimeFormat  = "2006-01-02 15:04:05"
	DateFormat      = "2006-01-02"
)

func ParseDateTimeLocation(dateStr string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormat, dateStr, TimeLocation)
}

func ParseDateLocation(dateStr string) (time.Time, error) {
	return time.ParseInLocation(DateFormat, dateStr, TimeLocation)
}

func FormatDateTime(t time.Time) string {
	return t.Format(DateTimeFormat)
}

func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// StartAndEndTimeDay 获取thatDay当日的开始时间和结束时间
func StartAndEndTimeDay(thatDay time.Time) (start, end time.Time) {
	start = time.Date(thatDay.Year(), thatDay.Month(), thatDay.Day(), 0, 0, 0, 0, thatDay.Location()) // 获取当日开始时间
	end = start.Add(24 * time.Hour).Add(-time.Millisecond * 1)                                        // 获取次日开始时间
	return start, end
}

func StartAndTimeOffsetDay(dateT time.Time, offsetDay int) time.Time {
	add := dateT.Add(-time.Hour * 24 * time.Duration(offsetDay))
	start := time.Date(add.Year(), add.Month(), add.Day(), 0, 0, 0, 0, add.Location()) // 获取当日开始时间r
	return start
}
