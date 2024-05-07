package utils

import "time"

const (
	F_DATETIME_DATE     = "2006-01-02"
	F_DATETIME_TIME     = "15:04:05"
	F_DATETIME_DATETIME = "2006-01-02 15:04:05"
)

type DateTime struct {
	t time.Time
	s string
}

func NewDateTime() DateTime {
	return DateTime{t: time.Now()}
}
func FromSecond(second int64) DateTime {
	p := DateTime{}
	p.t = time.Unix(second, 0)
	return p
}

func FromTime(t time.Time) DateTime {
	return DateTime{t: t}
}
func FromMillisecond(millisecond int64) DateTime {
	p := DateTime{}
	temp := millisecond / 1000
	p.t = time.Unix(temp, (millisecond%1000)*int64(time.Millisecond))
	return p
}
func FromLocal(year, month, day, hour, minute, second, millisecond int) DateTime {
	p := DateTime{}
	p.t = time.Date(year, time.Month(month), day, hour, minute, second, millisecond*int(time.Millisecond), time.Local)
	return p
}

func (p DateTime) Time() time.Time {
	return p.t
}
func (p DateTime) UnixSecond() int64 {
	return p.t.Unix()
}
func (p DateTime) UnixMillisecond() int64 {
	return p.t.UnixNano() / int64(time.Millisecond)
}

func (p DateTime) SetTime(t time.Time) DateTime {
	p.t = t
	return p
}
func (p DateTime) SetString(s string) DateTime {
	p.s = s
	return p
}

func ToDateTime(t time.Time) DateTime {
	return DateTime{t: t}
}
func (p DateTime) FormatDate() string {
	return p.t.Format(F_DATETIME_DATE)
}
func (p DateTime) FormatTime() string {
	return p.t.Format(F_DATETIME_TIME)
}
func (p DateTime) FormatDateTime() string {
	return p.t.Format(F_DATETIME_DATETIME)
}

func (p DateTime) ParseLocalDate() (time.Time, error) {
	return time.ParseInLocation(F_DATETIME_DATE, p.s, time.Local)
}
func (p DateTime) ParseLocalTime() (time.Time, error) {
	return time.ParseInLocation(F_DATETIME_TIME, p.s, time.Local)
}
func (p DateTime) ParseLocalDateTime() (time.Time, error) {
	return time.ParseInLocation(F_DATETIME_DATETIME, p.s, time.Local)
}
