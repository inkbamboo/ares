package datetime

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/jinzhu/now"
)

type LocalTime struct {
	now.Now `swaggerignore:"true"`
}

func Now() *LocalTime {
	return &LocalTime{
		*now.New(time.Now()),
	}
}

func NewLocalTime(stdTime ...time.Time) *LocalTime {
	c := new(LocalTime)
	if len(stdTime) > 0 {
		c.Now = *now.New(stdTime[0])
	} else {
		c.Now = *now.New(time.Time{})
	}
	return c
}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 || string(data) == "null" {
		*t = LocalTime{
			*now.New(time.Time{}),
		}
		return
	}
	var lTime time.Time
	lTime, err = now.Parse(strings.Trim(string(data), `"`))
	*t = LocalTime{
		*now.New(lTime),
	}
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	if !t.IsZero() {
		return []byte(fmt.Sprintf(`"%s"`, t.String())), nil
	}
	return []byte(`""`), nil
}

// Value 写入 mysql 时调用
func (t LocalTime) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return []byte(t.String()), nil
}

// Scan 检出 mysql 时调用
func (t *LocalTime) Scan(v interface{}) error {
	switch v.(type) {
	case time.Time:
		*t = *NewLocalTime(v.(time.Time))
	case *time.Time:
		*t = *NewLocalTime(*v.(*time.Time))
	case LocalTime:
		*t = v.(LocalTime)
	case *LocalTime:
		*t = *v.(*LocalTime)
	}
	return nil
}
func (t LocalTime) String() string {
	if t.IsZero() {
		return ""
	}
	return carbon.NewCarbon(t.Time).ToDateTimeString(time.Local.String())
}
func getAbsValue(value int64) int64 {
	return (value ^ (value >> 63)) - (value >> 63)
}
func (t *LocalTime) DiffInDays(lc ...*LocalTime) int64 {
	start := t
	var end *LocalTime
	if len(lc) > 0 {
		end = lc[0]
	} else {
		end = Now()
	}
	return int64(math.Floor(float64((end.Unix() - start.Unix()) / (carbon.HoursPerDay * carbon.SecondsPerHour))))
}

func (t *LocalTime) DiffAbsInDays(lc ...*LocalTime) int64 {
	return getAbsValue(t.DiffInDays(lc...))
}
func (t *LocalTime) DiffInHours(lc ...*LocalTime) int64 {
	start := t
	var end *LocalTime
	if len(lc) > 0 {
		end = lc[0]
	} else {
		end = Now()
	}
	return int64(math.Floor(float64((end.Unix() - start.Unix()) / carbon.SecondsPerHour)))
}

func (t *LocalTime) DiffAbsInHours(lc ...*LocalTime) int64 {
	return getAbsValue(t.DiffInHours(lc...))
}

func (t *LocalTime) DiffInMinutes(lc ...*LocalTime) int64 {
	start := t
	var end *LocalTime
	if len(lc) > 0 {
		end = lc[0]
	} else {
		end = Now()
	}
	return int64(math.Floor(float64((end.Unix() - start.Unix()) / carbon.SecondsPerMinute)))
}

func (t *LocalTime) DiffAbsInMinutes(lc ...*LocalTime) int64 {
	return getAbsValue(t.DiffInMinutes(lc...))
}

func (t *LocalTime) DiffInSeconds(lc ...*LocalTime) int64 {
	start := t
	var end *LocalTime
	if len(lc) > 0 {
		end = lc[0]
	} else {
		end = Now()
	}
	return end.Unix() - start.Unix()
}

func (t *LocalTime) DiffAbsInSeconds(lc ...*LocalTime) int64 {
	return getAbsValue(t.DiffInSeconds(lc...))
}
func (t LocalTime) TimeValue() time.Time {
	return t.Time
}
