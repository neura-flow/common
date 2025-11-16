package util

import (
	"fmt"
	"time"
)

const (
	EST = "America/New_York"
	UTC = "UTC"
	SH  = "Asia/Shanghai"

	YYYYMMDD            = "20060102"
	YYYY_MM_DD          = "2006-01-02"
	YYYY_MM_DD_HH_mm_ss = "2006-01-02 15:04:05" // 格式: yyyy-MM-dd HH:mm:ss
	ISO8601             = "2006-01-02T15:04:05.000Z07:00"
	ISO8601oOmitted     = "2006-01-02T15:04:05.999Z07:00"
)

func ParseTime(str, timezone string) (*time.Time, error) {
	return Parse(&str, String(YYYY_MM_DD_HH_mm_ss), &timezone)
}

func ParseDate(str, timezone string) (*time.Time, error) {
	return Parse(&str, String(YYYY_MM_DD), &timezone)
}

func Parse(str, layout, timezone *string) (*time.Time, error) {
	if AnyBlank(str, layout, timezone) {
		return nil, fmt.Errorf("invalid parameters")
	}
	tz, err := time.LoadLocation(*timezone)
	if err != nil {
		return nil, err
	}
	t, err := time.ParseInLocation(*layout, *str, tz)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func Trunc(t *time.Time, timezone *string) (*time.Time, error) {
	if t == nil || timezone == nil || *timezone == "" {
		return nil, fmt.Errorf("invalid parameters")
	}

	tz, err := time.LoadLocation(*timezone)
	if err != nil {
		return nil, err
	}

	d, err := time.ParseInLocation(YYYY_MM_DD, t.In(tz).Format(YYYY_MM_DD), tz)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func Hour(t *time.Time, timezone *string) (hour int, err error) {
	if t == nil || IsBlank(timezone) {
		return -1, fmt.Errorf("invalid parameters")
	}
	tz, err := time.LoadLocation(*timezone)
	if err != nil {
		return -1, err
	}
	hour = t.In(tz).Hour()
	return
}

func AddMin(t *time.Time, amount int, timezone *string) (*time.Time, error) {
	if t == nil || timezone == nil {
		return nil, fmt.Errorf("invalid parameters")
	}
	tz, err := time.LoadLocation(*timezone)
	if err != nil {
		return nil, err
	}
	t2 := t.In(tz).Add(time.Minute * time.Duration(amount))
	return &t2, nil
}

func AddSec(t *time.Time, amount int, timezone *string) (*time.Time, error) {
	if t == nil || timezone == nil {
		return nil, fmt.Errorf("invalid parameters")
	}
	tz, err := time.LoadLocation(*timezone)
	if err != nil {
		return nil, err
	}
	t2 := t.In(tz).Add(time.Second * time.Duration(amount))
	return &t2, nil
}

func AddDay(t *time.Time, amount int, timezone *string) (*time.Time, error) {
	if t == nil || timezone == nil {
		return nil, fmt.Errorf("invalid parameters")
	}
	tz, err := time.LoadLocation(*timezone)
	if err != nil {
		return nil, err
	}
	t2 := t.In(tz).AddDate(0, 0, amount)
	return &t2, nil
}

func Format(t *time.Time, layout, timezone *string) (*string, error) {
	if t == nil || layout == nil || timezone == nil {
		return nil, fmt.Errorf("invalid parameters")
	}
	if *layout == "" || *timezone == "" {
		return nil, fmt.Errorf("invalid parameters")
	}

	tz, err := time.LoadLocation(*timezone)
	if err != nil {
		return nil, err
	}

	s := t.In(tz).Format(*layout)
	return &s, nil
}

func FormatDateTime(t *time.Time, tz string) string {
	s, _ := Format(t, String(YYYY_MM_DD_HH_mm_ss), String(tz))
	if s == nil {
		return ""
	}
	return *s
}
func ParseDateTime(s, tz string) *time.Time {
	d, _ := Parse(String(s), String(YYYY_MM_DD_HH_mm_ss), String(tz))
	return d
}

func ParseUtcTimestamp(s string) *time.Time {
	i := ParseInt64(s)
	if i == nil {
		return nil
	}
	tz, _ := time.LoadLocation(UTC)
	tm := time.Unix(*i/1000, 0).In(tz)
	return &tm
}

func GetWeekDay(t *time.Time, timeZone string) (int, error) {
	tz, err := time.LoadLocation(timeZone)
	if err != nil {
		return -1, err
	}
	weekDay := t.In(tz).Weekday()
	return int(weekDay), nil
}
