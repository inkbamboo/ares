// Package datetime 提供时间处理工具，封装 carbon 和 now 库，支持 JSON/SQL 序列化。
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

// LocalTime 封装 now.Now，提供本地时间的序列化和计算功能。
type LocalTime struct {
	now.Now `swaggerignore:"true"`
}

// Now 返回当前时间的 LocalTime 实例。
func Now() *LocalTime {
	return &LocalTime{
		*now.New(time.Now()),
	}
}

// NewLocalTime 根据标准 time.Time 创建 LocalTime 实例，不传参则创建零值实例。
func NewLocalTime(stdTime ...time.Time) *LocalTime {
	c := new(LocalTime)
	if len(stdTime) > 0 {
		c.Now = *now.New(stdTime[0])
	} else {
		c.Now = *now.New(time.Time{})
	}
	return c
}

// UnmarshalJSON 实现 JSON 反序列化，支持多种时间格式解析。
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

// MarshalJSON 实现 JSON 序列化，零值返回空字符串。
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
	if v == nil {
		*t = *NewLocalTime()
		return nil
	}
	switch v.(type) {
	case time.Time:
		*t = *NewLocalTime(v.(time.Time))
	case *time.Time:
		*t = *NewLocalTime(*v.(*time.Time))
	case LocalTime:
		*t = v.(LocalTime)
	case *LocalTime:
		*t = *v.(*LocalTime)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}
// String 返回格式化的时间字符串（yyyy-MM-dd HH:mm:ss），零值返回空字符串。
func (t LocalTime) String() string {
	if t.IsZero() {
		return ""
	}
	return carbon.NewCarbon(t.Time).ToDateTimeString(time.Local.String())
}

// getAbsValue 返回 int64 的绝对值（位运算实现）。
func getAbsValue(value int64) int64 {
	return (value ^ (value >> 63)) - (value >> 63)
}

// DiffInDays 计算与目标时间的天数差（有符号），不传参则与当前时间比较。
func (t *LocalTime) DiffInDays(lc ...*LocalTime) int64 {
	start := t
	var end *LocalTime
	if len(lc) > 0 {
		end = lc[0]
	} else {
		end = Now()
	}
	return int64(math.Floor(float64(end.Unix()-start.Unix()) / float64(carbon.HoursPerDay*carbon.SecondsPerHour)))
}

// DiffAbsInDays 计算与目标时间的天数差（绝对值）。
func (t *LocalTime) DiffAbsInDays(lc ...*LocalTime) int64 {
	return getAbsValue(t.DiffInDays(lc...))
}

// DiffInHours 计算与目标时间的小时差（有符号）。
func (t *LocalTime) DiffInHours(lc ...*LocalTime) int64 {
	start := t
	var end *LocalTime
	if len(lc) > 0 {
		end = lc[0]
	} else {
		end = Now()
	}
	return int64(math.Floor(float64(end.Unix()-start.Unix()) / float64(carbon.SecondsPerHour)))
}

// DiffAbsInHours 计算与目标时间的小时差（绝对值）。
func (t *LocalTime) DiffAbsInHours(lc ...*LocalTime) int64 {
	return getAbsValue(t.DiffInHours(lc...))
}

// DiffInMinutes 计算与目标时间的分钟差（有符号）。
func (t *LocalTime) DiffInMinutes(lc ...*LocalTime) int64 {
	start := t
	var end *LocalTime
	if len(lc) > 0 {
		end = lc[0]
	} else {
		end = Now()
	}
	return int64(math.Floor(float64(end.Unix()-start.Unix()) / float64(carbon.SecondsPerMinute)))
}

// DiffAbsInMinutes 计算与目标时间的分钟差（绝对值）。
func (t *LocalTime) DiffAbsInMinutes(lc ...*LocalTime) int64 {
	return getAbsValue(t.DiffInMinutes(lc...))
}

// DiffInSeconds 计算与目标时间的秒数差（有符号）。
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

// DiffAbsInSeconds 计算与目标时间的秒数差（绝对值）。
func (t *LocalTime) DiffAbsInSeconds(lc ...*LocalTime) int64 {
	return getAbsValue(t.DiffInSeconds(lc...))
}

// TimeValue 返回底层的 time.Time 值。
func (t LocalTime) TimeValue() time.Time {
	return t.Time
}
